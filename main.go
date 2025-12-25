package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	store := NewUserStore()
	
	// Add some data
	store.Create(&User{Name: "Alice", Email: "alice@example.com"})
	store.Create(&User{Name: "Bob", Email: "bob@example.com"})
	
	encoder := NewTOONEncoder()
	decoder := NewTOONDecoder()
	
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetUsers(w, r, store, encoder)
		case http.MethodPost:
			handleCreateUser(w, r, store, encoder, decoder)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	
	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/users/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		
		switch r.Method {
		case http.MethodGet:
			handleGetUser(w, r, store, encoder, id)
		case http.MethodPut:
			handleUpdateUser(w, r, store, encoder, decoder, id)
		case http.MethodDelete:
			handleDeleteUser(w, r, store, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	
	port := ":8080"
	log.Printf("Starting TOON HTTP server on %s", port)
	log.Printf("Try: curl http://localhost:8080/users")
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleGetUsers(w http.ResponseWriter, r *http.Request, store *UserStore, encoder *TOONEncoder) {
	users := store.GetAll()
	
	toon, err := encoder.Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/toon")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(toon))
}

func handleGetUser(w http.ResponseWriter, r *http.Request, store *UserStore, encoder *TOONEncoder, id int) {
	user, exists := store.Get(id)
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	toon, err := encoder.Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/toon")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(toon))
}

func handleCreateUser(w http.ResponseWriter, r *http.Request, store *UserStore, encoder *TOONEncoder, decoder *TOONDecoder) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	
	var user User
	if err := decoder.Decode(string(body), &user); err != nil {
		http.Error(w, fmt.Sprintf("Invalid TOON data: %v", err), http.StatusBadRequest)
		return
	}
	
	created := store.Create(&user)
	
	toon, err := encoder.Encode(created)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/toon")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(toon))
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request, store *UserStore, encoder *TOONEncoder, decoder *TOONDecoder, id int) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	
	var user User
	if err := decoder.Decode(string(body), &user); err != nil {
		http.Error(w, fmt.Sprintf("Invalid TOON data: %v", err), http.StatusBadRequest)
		return
	}
	
	updated, exists := store.Update(id, &user)
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	toon, err := encoder.Encode(updated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/toon")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(toon))
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request, store *UserStore, id int) {
	if !store.Delete(id) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

