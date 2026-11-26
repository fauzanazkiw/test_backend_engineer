package response

import (
	"encoding/json"
	"log"
	"net/http"
)

// JSON menulis payload apa adanya sebagai response JSON.
func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("gagal encode response: %v", err)
	}
}

// Data mengembalikan response berbentuk {"data": ...}.
func Data(w http.ResponseWriter, status int, data any) {
	JSON(w, status, map[string]any{"data": data})
}

// Message mengembalikan response berbentuk {"message": "..."}.
func Message(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]any{"message": message})
}

// MessageData mengembalikan response berbentuk {"message": "...", "data": ...}.
func MessageData(w http.ResponseWriter, status int, message string, data any) {
	JSON(w, status, map[string]any{
		"message": message,
		"data":    data,
	})
}

// Error mengembalikan response error berbentuk {"message": "..."}.
func Error(w http.ResponseWriter, status int, message string) {
	Message(w, status, message)
}

// ValidationError mengembalikan response validasi sesuai standar PRD.
func ValidationError(w http.ResponseWriter, fieldErrors map[string][]string) {
	JSON(w, http.StatusBadRequest, map[string]any{
		"message": "validation failed",
		"errors":  fieldErrors,
	})
}
