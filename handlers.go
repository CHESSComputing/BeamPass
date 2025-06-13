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
	beamline := r.URL.Query().Get("beamline")
	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")

	var err error

	if startTimeStr, err = parseDate(startTimeStr); err != nil {
		http.Error(w, "Invalid end_time format", http.StatusBadRequest)
		return
	}
	if endTimeStr, err = parseDate(endTimeStr); err != nil {
		http.Error(w, "Invalid end_time format", http.StatusBadRequest)
		return
	}

	// Validate parameters
	if beamline == "" {
		http.Error(w, `{"error": "Missing required parameters: beamline"}`, http.StatusBadRequest)
		return
	}

	// Execute the getBTR function
	data, err := getBTR(beamline, startTimeStr, endTimeStr)
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
