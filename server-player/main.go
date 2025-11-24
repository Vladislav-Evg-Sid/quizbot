package main

import (
	"encoding/json"
	"log"
	"net/http"
	"quiz-server-player/database"
	"quiz-server-player/models"

	"github.com/gorilla/mux"
)

func main() {
	// Инициализируем БД
	if err := database.Init(); err != nil {
		log.Fatal("Database init failed:", err)
	}

	r := mux.NewRouter()

	// Эндпоинты
	r.HandleFunc("/api/users/topics", handleAllTopics).Methods("GET") // Получение всех тем викторины

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	log.Println("🚀 Player service started on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func handleAllTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := database.GetAllTopics()
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, `{"success": false, "error": "Database error"}`, http.StatusInternalServerError)
		return
	}

	if topics == nil {
		http.Error(w, `{"success": false, "error": "User not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.AllTopicsResponse{
		Success: true,
		Topics:  topics,
	})
}
