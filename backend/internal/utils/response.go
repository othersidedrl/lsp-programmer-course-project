// Package utils berisi pembantu umum: format respons JSON dan JWT.
package utils

import (
	"encoding/json"
	"net/http"
)

// Response adalah format standar seluruh respons API.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// JSON menulis respons JSON dengan status HTTP tertentu.
func JSON(w http.ResponseWriter, status int, payload Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// OK adalah pembantu untuk respons sukses (HTTP 200).
func OK(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusOK, Response{Success: true, Message: message, Data: data})
}

// Error adalah pembantu untuk respons gagal.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, Response{Success: false, Message: message})
}
