package main

import (
	"encoding/json"
	"net/http"
)

type JSONResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// Success helper
func writeJSONSuccess(w http.ResponseWriter, status int, message string, data any) error {
	return writeJSON(w, status, JSONResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error helper
func writeJSONError(w http.ResponseWriter, status int, message string, err error) error {
	errString := "unknown error"
	if err != nil {
		errString = err.Error()
	}

	return writeJSON(w, status, JSONResponse{
		Success: false,
		Message: message,
		Error:   errString,
	})
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(data)
	if err != nil {
		return err
	}
	return nil
}
