package main

import "errors"

func successResponse(message string, data any) map[string]any {
	return map[string]any{
		"success": true,
		"message": message,
		"data":    data,
	}
}

func errorResponse(message string, err error) map[string]any {
	if err == nil {
		err = errors.New("unknown error")
	}
	return map[string]any{
		"success": false,
		"message": message,
		"error":   err.Error(),
	}
}
