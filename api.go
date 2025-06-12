package main

import (
	"fmt"
)

// BTRData struct to hold the results of the MySQL query
type BTRData struct {
	ScheduleEntryFileID string `json:"schedule_entry_file_id"`
	ResourceName        string `json:"resource_name"`
	StartDatetime       string `json:"start_datetime"`
	EndDatetime         string `json:"end_datetime"`
}

// getBTR performs the MySQL query and returns the results
func getBTR(resourceName string, startTime, endTime string) ([]BTRData, error) {
	query := `
		SELECT br.schedule_entry_file_id, r.name, se.start_datetime, se.end_datetime
		FROM beampass.resource r
		JOIN beampass.schedule_entry se ON se.resource_id = r.id
		JOIN beampass.beamtime_request br ON se.beamtime_request_id = br.id
		WHERE r.name = ? AND se.is_actual = true AND se.start_datetime >= ? AND se.end_datetime <= ?
		ORDER BY se.start_datetime;
	`

	rows, err := db.Query(query, resourceName, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var results []BTRData
	for rows.Next() {
		var data BTRData
		var fileID, stime, etime string
		if err := rows.Scan(&data.ScheduleEntryFileID, &data.ResourceName, &data.StartDatetime, &data.EndDatetime); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		results = append(results, data)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return results, nil
}
