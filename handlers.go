package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// BtrHandler handles the /btr endpoint
func BtrHandler(w http.ResponseWriter, r *http.Request) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	const layout = "2006-01-02 15:04:05"
	const shortLayout = "20060102"

	resourceName := r.URL.Query().Get("beamline")
	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")

	var startTime, endTime time.Time
	var err error

	now := time.Now()
	// Case 1: Neither start nor end provided
	if startTimeStr == "" && endTimeStr == "" {
		endTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		startTime = endTime.Add(-24 * time.Hour)
	} else if startTimeStr != "" && endTimeStr == "" {
		// Case 2: start provided, end not
		startTime, err = parseDate(startTimeStr)
		if err != nil {
			http.Error(w, "Invalid start_time format", http.StatusBadRequest)
			return
		}
		endTime = startTime.Add(24 * time.Hour)
	} else {
		// Case 3: both start and end provided
		startTime, err = parseDate(startTimeStr)
		if err != nil {
			http.Error(w, "Invalid start_time format", http.StatusBadRequest)
			return
		}
		endTime, err = parseDate(endTimeStr)
		if err != nil {
			http.Error(w, "Invalid end_time format", http.StatusBadRequest)
			return
		}
	}
	startTimeStr = startTime.Format(layout)
	endTimeStr = endTime.Format(layout)

	// Validate parameters
	if resourceName == "" || startTimeStr == "" || endTimeStr == "" {
		http.Error(w, `{"error": "Missing required parameters: name, start_time, and end_time"}`, http.StatusBadRequest)
		return
	}

	// Execute the getBTR function
	data, err := getBTR(resourceName, startTimeStr, endTimeStr)
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
