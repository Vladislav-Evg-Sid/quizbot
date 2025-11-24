package main

import (
	"encoding/json"
	"log"
	"net/http"
	"quiz-server-admin/database"
	"quiz-server-admin/models"

	"github.com/gorilla/mux"
)

func main() {
	// Инициализируем БД
	if err := database.Init(); err != nil {
		log.Fatal("Database init failed:", err)
	}

	r := mux.NewRouter()

	// Эндпоинты
	r.HandleFunc("/api/users/start", handleStart).Methods("POST") // /start

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	log.Println("🚀 Admin service started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	var req models.StartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success": false, "error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	user := &models.User{
		TelegramID: req.TelegramID,
		Name:       req.Name,
		Username:   req.Username,
		IsAdmin:    false,
	}

	if err := database.CreateOrUpdateUser(user); err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, `{"success": false, "error": "Database error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.StartResponse{
		Success: true,
		User:    *user,
	})
}
