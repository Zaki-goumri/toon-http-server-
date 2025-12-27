package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/zaki/toon-go-server/pkg/models"
	"github.com/zaki/toon-go-server/pkg/toon"
)

type ToonHandler struct {
	store   *models.UserStore
	encoder *toon.TOONEncoder
	decoder *toon.TOONDecoder
}

func NewToonHandler(store *models.UserStore) *ToonHandler {
	return &ToonHandler{
		store:   store,
		encoder: toon.NewTOONEncoder(),
		decoder: toon.NewTOONDecoder(),
	}
}

func (h *ToonHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetUsers(w, r)
	case http.MethodPost:
		h.handleCreateUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ToonHandler) HandleUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetUser(w, r, id)
	case http.MethodPut:
		h.handleUpdateUser(w, r, id)
	case http.MethodDelete:
		h.handleDeleteUser(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ToonHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users := h.store.GetAll()

	toonData, err := h.encoder.Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/toon")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(toonData))
}

func (h *ToonHandler) handleGetUser(w http.ResponseWriter, r *http.Request, id int) {
	user, exists := h.store.Get(id)
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	toonData, err := h.encoder.Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/toon")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(toonData))
}

func (h *ToonHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := h.decoder.Decode(string(body), &user); err != nil {
		http.Error(w, fmt.Sprintf("Invalid TOON data: %v", err), http.StatusBadRequest)
		return
	}

	created := h.store.Create(&user)

	toonData, err := h.encoder.Encode(created)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/toon")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(toonData))
}

func (h *ToonHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request, id int) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := h.decoder.Decode(string(body), &user); err != nil {
		http.Error(w, fmt.Sprintf("Invalid TOON data: %v", err), http.StatusBadRequest)
		return
	}

	updated, exists := h.store.Update(id, &user)
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	toonData, err := h.encoder.Encode(updated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/toon")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(toonData))
}

func (h *ToonHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request, id int) {
	if !h.store.Delete(id) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// JSONHandler handles JSON requests
type JSONHandler struct {
	store *models.UserStore
}

func NewJSONHandler(store *models.UserStore) *JSONHandler {
	return &JSONHandler{
		store: store,
	}
}

func (h *JSONHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetUsers(w, r)
	case http.MethodPost:
		h.handleCreateUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *JSONHandler) HandleUser(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/json/users/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetUser(w, r, id)
	case http.MethodPut:
		h.handleUpdateUser(w, r, id)
	case http.MethodDelete:
		h.handleDeleteUser(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *JSONHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users := h.store.GetAll()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func (h *JSONHandler) handleGetUser(w http.ResponseWriter, r *http.Request, id int) {
	user, exists := h.store.Get(id)
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *JSONHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON data: %v", err), http.StatusBadRequest)
		return
	}

	created := h.store.Create(&user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *JSONHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request, id int) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON data: %v", err), http.StatusBadRequest)
		return
	}

	updated, exists := h.store.Update(id, &user)
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

func (h *JSONHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request, id int) {
	if !h.store.Delete(id) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
