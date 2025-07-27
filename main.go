package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// getHandler handles GET requests and URL parameters
func getHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	param1 := params.Get("param1")
	param2 := params.Get("param2")

	response := fmt.Sprintf("GET request received with param1: %s, param2: %s", param1, param2)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

// postHandler handles POST requests
func postHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	response := fmt.Sprintf("POST request received with body: %s", string(body))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

// jsonHandler handles JSON POST requests
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	responseData := map[string]interface{}{
		"message": "JSON request received",
		"data":    data,
	}
	responseJSON, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// uploadHandler handles file uploads
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the file to a temporary location
	tempFile, err := ioutil.TempFile("", header.Filename)
	if err != nil {
		http.Error(w, "Error creating temporary file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("File uploaded successfully: %s", header.Filename)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

// downloadHandler handles file downloads
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "sample.txt" // Replace with an actual file path for testing
	fileName := filepath.Base(filePath)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}

func main() {
	// Register handlers for different routes
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/json", jsonHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)

	// Start the server
	fmt.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
