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

	handler := handlers.NewJSONHandler(store)

	http.HandleFunc("/json/users", handler.HandleUsers)
	http.HandleFunc("/json/users/", handler.HandleUser)

	port := ":8081"
	log.Printf("Starting JSON HTTP server on %s", port)
	log.Printf("Try: curl http://localhost:8081/json/users")
	log.Fatal(http.ListenAndServe(port, nil))
}
