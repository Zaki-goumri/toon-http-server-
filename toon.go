package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type TOONEncoder struct {
	indentSize int
}

func NewTOONEncoder() *TOONEncoder {
	return &TOONEncoder{
		indentSize: 2,
	}
}

func (e *TOONEncoder) Encode(v interface{}) (string, error) {
	return e.encodeValue(v, 0)
}

func (e *TOONEncoder) encodeValue(v interface{}, depth int) (string, error) {
	if v == nil {
		return "null", nil
	}

	val := reflect.ValueOf(v)
	
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return "null", nil
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		return e.encodeArray(val, depth)
	case reflect.Struct:
		return e.encodeObject(val, depth)
	case reflect.Map:
		return e.encodeMap(val, depth)
	case reflect.String:
		return e.encodeString(val.String()), nil
	case reflect.Bool:
		if val.Bool() {
			return "true", nil
		}
		return "false", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(val.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'f', -1, 64), nil
	default:
		return "", fmt.Errorf("unsupported type: %v", val.Kind())
	}
}

func (e *TOONEncoder) encodeArray(val reflect.Value, depth int) (string, error) {
	length := val.Len()
	if length == 0 {
		return "[]", nil
	}

	var sb strings.Builder
	
	firstElem := val.Index(0)
	for firstElem.Kind() == reflect.Ptr {
		if firstElem.IsNil() {
			return "", fmt.Errorf("nil element in array")
		}
		firstElem = firstElem.Elem()
	}

	if firstElem.Kind() == reflect.Struct {
		fields := e.getStructFields(firstElem.Type())
		
		sb.WriteString(fmt.Sprintf("[%d", length))
		for _, field := range fields {
			sb.WriteString(" ")
			sb.WriteString(field.toonName)
		}
		sb.WriteString("]\n")

		indent := strings.Repeat(" ", depth*e.indentSize)
		for i := 0; i < length; i++ {
			elem := val.Index(i)
			for elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			
			sb.WriteString(indent)
			for j, field := range fields {
				if j > 0 {
					sb.WriteString(",")
				}
				fieldVal := elem.FieldByIndex(field.index)
				encoded, err := e.encodeFieldValue(fieldVal)
				if err != nil {
					return "", err
				}
				sb.WriteString(encoded)
			}
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString(fmt.Sprintf("[%d]\n", length))
		indent := strings.Repeat(" ", depth*e.indentSize)
		for i := 0; i < length; i++ {
			elem := val.Index(i)
			sb.WriteString(indent)
			encoded, err := e.encodeValue(elem.Interface(), depth+1)
			if err != nil {
				return "", err
			}
			sb.WriteString(encoded)
			sb.WriteString("\n")
		}
	}

	return sb.String(), nil
}

func (e *TOONEncoder) encodeObject(val reflect.Value, depth int) (string, error) {
	var sb strings.Builder
	indent := strings.Repeat(" ", depth*e.indentSize)
	fields := e.getStructFields(val.Type())

	for _, field := range fields {
		fieldVal := val.FieldByIndex(field.index)
		sb.WriteString(indent)
		sb.WriteString(field.toonName)
		sb.WriteString(": ")
		
		encoded, err := e.encodeFieldValue(fieldVal)
		if err != nil {
			return "", err
		}
		sb.WriteString(encoded)
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func (e *TOONEncoder) encodeMap(val reflect.Value, depth int) (string, error) {
	var sb strings.Builder
	indent := strings.Repeat(" ", depth*e.indentSize)
	
	for _, key := range val.MapKeys() {
		mapVal := val.MapIndex(key)
		sb.WriteString(indent)
		sb.WriteString(e.encodeString(fmt.Sprint(key.Interface())))
		sb.WriteString(": ")
		
		encoded, err := e.encodeValue(mapVal.Interface(), depth+1)
		if err != nil {
			return "", err
		}
		sb.WriteString(encoded)
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func (e *TOONEncoder) encodeFieldValue(val reflect.Value) (string, error) {
	// Handle time.Time specially
	if val.Type() == reflect.TypeOf(time.Time{}) {
		t := val.Interface().(time.Time)
		return e.encodeString(t.Format(time.RFC3339)), nil
	}
	
	return e.encodeValue(val.Interface(), 0)
}

func (e *TOONEncoder) encodeString(s string) string {
	needsQuoting := s == "" || 
		s == "true" || s == "false" || s == "null" ||
		strings.ContainsAny(s, " \t\n\r,:[]{}\"") ||
		strings.HasPrefix(s, "-")

	if !needsQuoting {
		if _, err := strconv.ParseFloat(s, 64); err == nil {
			needsQuoting = true
		}
	}

	if !needsQuoting {
		return s
	}

	var sb strings.Builder
	sb.WriteString("\"")
	for _, ch := range s {
		switch ch {
		case '\\':
			sb.WriteString("\\\\")
		case '"':
			sb.WriteString("\\\"")
		case '\n':
			sb.WriteString("\\n")
		case '\r':
			sb.WriteString("\\r")
		case '\t':
			sb.WriteString("\\t")
		default:
			sb.WriteRune(ch)
		}
	}
	sb.WriteString("\"")
	return sb.String()
}

type structField struct {
	toonName string
	index    []int
}

func (e *TOONEncoder) getStructFields(t reflect.Type) []structField {
	var fields []structField
	
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		
		toonName := field.Tag.Get("toon")
		if toonName == "" {
			toonName = strings.ToLower(field.Name[:1]) + field.Name[1:]
		}
		
		fields = append(fields, structField{
			toonName: toonName,
			index:    field.Index,
		})
	}
	
	return fields
}

type TOONDecoder struct {
}

func NewTOONDecoder() *TOONDecoder {
	return &TOONDecoder{}
}

func (d *TOONDecoder) Decode(data string, v interface{}) error {
	lines := strings.Split(strings.TrimSpace(data), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty TOON data")
	}

	firstLine := strings.TrimSpace(lines[0])
	
	if strings.HasPrefix(firstLine, "[") {
		return d.decodeArray(lines, v)
	}
	
	return d.decodeObject(lines, v)
}

func (d *TOONDecoder) decodeObject(lines []string, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("v must be a non-nil pointer")
	}
	
	elem := rv.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("v must point to a struct")
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		field := d.findField(elem, key)
		if !field.IsValid() || !field.CanSet() {
			continue
		}

		if err := d.setFieldValue(field, value); err != nil {
			return err
		}
	}

	return nil
}

func (d *TOONDecoder) decodeArray(lines []string, v interface{}) error {
	if len(lines) == 0 {
		return fmt.Errorf("empty array")
	}

	header := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(header, "[") || !strings.HasSuffix(header, "]") {
		return fmt.Errorf("invalid array header")
	}

	headerParts := strings.Fields(header[1 : len(header)-1])
	if len(headerParts) == 0 {
		return fmt.Errorf("invalid array header")
	}

	length, err := strconv.Atoi(headerParts[0])
	if err != nil {
		return fmt.Errorf("invalid array length: %v", err)
	}

	fields := headerParts[1:]

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("v must be a non-nil pointer")
	}

	slice := rv.Elem()
	if slice.Kind() != reflect.Slice {
		return fmt.Errorf("v must point to a slice")
	}

	elemType := slice.Type().Elem()
	newSlice := reflect.MakeSlice(slice.Type(), 0, length)

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		values := strings.Split(line, ",")
		if len(values) != len(fields) {
			continue
		}

		elemPtr := reflect.New(elemType)
		elem := elemPtr.Elem()

		for j, fieldName := range fields {
			field := d.findField(elem, strings.TrimSpace(fieldName))
			if !field.IsValid() || !field.CanSet() {
				continue
			}

			value := strings.TrimSpace(values[j])
			if err := d.setFieldValue(field, value); err != nil {
				return err
			}
		}

		newSlice = reflect.Append(newSlice, elemPtr)
	}

	slice.Set(newSlice)
	return nil
}

func (d *TOONDecoder) findField(elem reflect.Value, name string) reflect.Value {
	t := elem.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		toonName := field.Tag.Get("toon")
		if toonName == "" {
			toonName = strings.ToLower(field.Name[:1]) + field.Name[1:]
		}
		if toonName == name {
			return elem.Field(i)
		}
	}
	return reflect.Value{}
}

func (d *TOONDecoder) setFieldValue(field reflect.Value, value string) error {
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
		value = value[1 : len(value)-1]
		value = strings.ReplaceAll(value, "\\n", "\n")
		value = strings.ReplaceAll(value, "\\r", "\r")
		value = strings.ReplaceAll(value, "\\t", "\t")
		value = strings.ReplaceAll(value, "\\\"", "\"")
		value = strings.ReplaceAll(value, "\\\\", "\\")
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(u)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(t))
		}
	}

	return nil
}
