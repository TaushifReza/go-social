package main

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// Look for the "json" tag on the struct field
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		// If the tag is "-", ignore it
		if name == "-" {
			return ""
		}

		return name
	})
}

type JSONResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
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
func writeJSONError(w http.ResponseWriter, status int, message string, err any) error {
	// If err is nil, provide a fallback, otherwise use the err as-is
	var errData any = "unknown error"
	if err != nil {
		errData = err
	}

	return writeJSON(w, status, JSONResponse{
		Success: false,
		Message: message,
		Error:   errData,
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
