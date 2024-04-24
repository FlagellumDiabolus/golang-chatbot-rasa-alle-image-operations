package handlers

import (
	"chatbot-ai/internal/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func isFileOperation(message string) bool {
	return strings.Contains(strings.ToLower(message), "save") ||
		strings.Contains(strings.ToLower(message), "retrieve") ||
		strings.Contains(strings.ToLower(message), "get")
}

func processFileOperation(w http.ResponseWriter, message string) string {
	switch {
	case strings.Contains(strings.ToLower(message), "save"):
		saveImage(w, message)
		return "image saved"
	case strings.Contains(strings.ToLower(message), "retrieve") || strings.Contains(strings.ToLower(message), "get"):
		retrieveImage(w, message)
		return "image retrieved"
	default:
		http.Error(w, "File operation not supported", http.StatusBadRequest)
	}
	return ""
}

func saveImage(w http.ResponseWriter, message string) {
	imageName := extractImageName(message)
	if imageName == "" {
		http.Error(w, "Image name not provided", http.StatusBadRequest)
		return
	}

	imageURL := extractImageURL(message)
	if imageURL == "" {
		http.Error(w, "Image URL not provided", http.StatusBadRequest)
		return
	}

	err := database.SaveImage(imageName, imageURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to save image: %v", err), http.StatusInternalServerError)
		return
	}

	response := "Image saved successfully."
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

func retrieveImage(w http.ResponseWriter, message string) {
	imageName := extractImageName(message)
	if imageName == "" {
		http.Error(w, "Image name not provided", http.StatusBadRequest)
		return
	}

	imageURL, err := database.RetrieveImage(imageName)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("Image '%s' not found", imageName), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to retrieve image: %v", err), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("Retrieved image link: %s", imageURL)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

func extractImageName(message string) string {
	parts := strings.Fields(message)
	if len(parts) < 3 {
		return ""
	}
	return parts[2]
}

func extractImageURL(message string) string {
	parts := strings.Fields(message)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}
