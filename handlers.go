package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// BtrHandler handles the /btr endpoint
func BtrHandler(w http.ResponseWriter, r *http.Request) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	resourceName := r.URL.Query().Get("name")
	startTime := r.URL.Query().Get("start_time")
	endTime := r.URL.Query().Get("end_time")

	// Validate parameters
	if resourceName == "" || startTime == "" || endTime == "" {
		http.Error(w, `{"error": "Missing required parameters: name, start_time, and end_time"}`, http.StatusBadRequest)
		return
	}

	// Execute the getBTR function
	data, err := getBTR(resourceName, startTime, endTime)
	if err != nil {
		log.Printf("Error in getBTR: %v", err)
		http.Error(w, fmt.Sprintf(`{"error": "Internal server error: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// Encode and send the response
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
	}
}
