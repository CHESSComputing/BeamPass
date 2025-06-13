package main

import (
	"fmt"
	"log"
	"time"
)

// BTRData struct to hold the results of the MySQL query
type BTRData struct {
	Btr       string `json:"btr"`
	Beamline  string `json:"beamline"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// getBTR performs the MySQL query and returns the results
func getBTR(beamline string, startTime, endTime string) ([]BTRData, error) {
	query := `
		SELECT br.schedule_entry_file_id as btr, r.name as beamline, se.start_datetime, se.end_datetime
		FROM beampass.resource r
		JOIN beampass.schedule_entry se ON se.resource_id = r.id
		JOIN beampass.beamtime_request br ON se.beamtime_request_id = br.id
		WHERE r.name = ? AND se.is_actual = true AND se.start_datetime >= ? AND se.end_datetime <= ?
		ORDER BY se.start_datetime;
	`
	if _verbose > 0 {
		log.Printf("QUERY: %s, beamline=%s startTime=%s endTime=%s", query, beamline, startTime, endTime)
	}

	rows, err := db.Query(query, beamline, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var results []BTRData
	for rows.Next() {
		var data BTRData
		if err := rows.Scan(&data.Btr, &data.Beamline, &data.StartTime, &data.EndTime); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		results = append(results, data)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return results, nil
}

func parseDate(s string) (time.Time, error) {
	if len(s) == 8 {
		// Format: YYYYMMDD
		return time.ParseInLocation("20060102", s, time.Local)
	}
	return time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
}
