package main

import (
	"log"
	"net/http"

	"github.com/zaki/toon-go-server/pkg/handlers"
	"github.com/zaki/toon-go-server/pkg/models"
)

func main() {
	store := models.NewUserStore()

	// Add some initial data
	store.Create(&models.User{Name: "Alice", Email: "alice@example.com"})
	store.Create(&models.User{Name: "Bob", Email: "bob@example.com"})

	handler := handlers.NewToonHandler(store)

	http.HandleFunc("/users", handler.HandleUsers)
	http.HandleFunc("/users/", handler.HandleUser)

	port := ":8080"
	log.Printf("Starting TOON HTTP server on %s", port)
	log.Printf("Try: curl http://localhost:8080/users")
	log.Fatal(http.ListenAndServe(port, nil))
}
