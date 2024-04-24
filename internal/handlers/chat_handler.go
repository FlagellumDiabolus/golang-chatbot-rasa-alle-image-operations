package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const rasaServerURL = "http://localhost:5005/model/parse"

type request struct {
	Message string
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	//req := request{}
	//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//	http.Error(w, "Failed to parse request body", http.StatusBadRequest)
	//	return
	//}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	reqBody := string(reqBytes)

	if !isFileOperation(reqBody) {
		payload := map[string]string{
			"text": reqBody,
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, "Failed to marshal request payload", http.StatusInternalServerError)
			return
		}

		resp, err := http.Post(rasaServerURL, "application/json", bytes.NewReader(payloadBytes))
		if err != nil {
			//http.Error(w, fmt.Sprintf("Failed to send request to Rasa server: %v", err), http.StatusInternalServerError)
			fmt.Fprintf(w, "Connection to rasa not there, handling the non-rasa way\n")
			NonRasaWay(w, reqBody)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf("Rasa server returned non-OK status code: %d", resp.StatusCode), http.StatusInternalServerError)
			return
		}

		var rasaResp struct {
			Intent string `json:"intent"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&rasaResp); err != nil {
			http.Error(w, "Failed to decode Rasa server response", http.StatusInternalServerError)
			return
		}

		response := fmt.Sprintf("Intent: %s", rasaResp.Intent)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"response": response})
	} else {
		processFileOperation(w, reqBody)
	}
}

// no Rasa dependency
func NonRasaWay(w http.ResponseWriter, message string) {
	if isFileOperation(message) {
		res := processFileOperation(w, message)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"response": "%s"}`, res)
		_, err := w.Write([]byte(res))
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"response": "You said: %s"}`, message)
}
