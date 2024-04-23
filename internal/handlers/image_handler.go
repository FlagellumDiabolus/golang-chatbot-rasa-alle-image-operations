package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang-chatbot-alle-image_operations/internal/database"
)

func SaveImageHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	err := database.SaveImage(req.Name, req.URL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to save image: %v", err), http.StatusInternalServerError)
		return
	}

	response := "Image saved successfully."
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

func RetrieveImageHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	imageURL, err := database.RetrieveImage(req.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve image: %v", err), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("Retrieved image link: %s", imageURL)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

func ListImagesHandler(w http.ResponseWriter, r *http.Request) {
	images, err := database.ListImages()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list images: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(images)
}
