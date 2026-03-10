package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"slices"
	"time"
)

// BTRData struct to hold the results of the MySQL query
type BTRData struct {
	Btr       string `json:"btr"`
	Beamline  string `json:"beamline"`
	PI        string `json:"pi"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// getBTR performs the MySQL query and returns the results
func getBTR(btrs []string, beamline, startTime, endTime, dateTime string) ([]BTRData, error) {
	var err error
	var rows *sql.Rows

	// Common to all queries, no matter the user parameters
	baseQuery := `
				SELECT
            br.schedule_entry_file_id AS btr,
            r.name AS beamline,
            pn.last_name,
            se.start_datetime,
            se.end_datetime
        FROM beampass.resource r
        JOIN beampass.schedule_entry se ON se.resource_id = r.id
        JOIN beampass.beamtime_request br ON se.beamtime_request_id = br.id
        JOIN beampass.project p ON br.project_id = p.id
        JOIN beampass.affiliation a ON p.lead_id = a.id
        JOIN beampass.person pn ON a.person_id = pn.id
    `
	var whereClauses []string
	var queryArgs []any

	whereClauses = append(whereClauses, "se.is_actual = true")

	if beamline != "" {
		whereClauses = append(whereClauses, "r.name = ?")
		queryArgs = append(queryArgs, beamline)
	}

	if btrs != nil {
		whereClauses = append(whereClauses, "br.schedule_entry_file_id IN (" + strings.Join(slices.Repeat([]string{"?"}, len(btrs)), ",") + ")")
		for _, btr := range(btrs) {
			queryArgs = append(queryArgs, btr)
		}
	}

	if dateTime != "" {
		whereClauses = append(whereClauses, "se.start_datetime < ? AND se.end_datetime > ?")
		queryArgs = append(queryArgs, dateTime, dateTime)
	} else {
		if startTime != "" {
			whereClauses = append(whereClauses, "se.start_datetime >= ?")
			queryArgs = append(queryArgs, startTime)
		}

		if endTime != "" {
			whereClauses = append(whereClauses, "se.end_datetime <= ?")
			queryArgs = append(queryArgs, endTime)
		}
	}

	query := baseQuery
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}
	query += " ORDER BY se.start_datetime;"

	if _verbose > 0 {
		log.Printf("QUERY: %s, queryArgs=%v", query, queryArgs)
	}

	rows, err = db.Query(query, queryArgs...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var results []BTRData
	for rows.Next() {
		var data BTRData
		if err := rows.Scan(&data.Btr, &data.Beamline, &data.PI, &data.StartTime, &data.EndTime); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		results = append(results, data)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return results, nil
}

func parseDate(s string) (string, error) {
	const layout = "2006-01-02 15:04:05"
	var err error
	if s == "" {
		return s, nil
	}
	var tstmp time.Time
	if len(s) == 8 {
		// Format: YYYYMMDD
		tstmp, err = time.ParseInLocation("20060102", s, time.Local)
	} else if len(s) == 10 {
		// Format: YYYY-MM-DD
		tstmp, err = time.ParseInLocation("2006-01-02", s, time.Local)
	} else {
		tstmp, err = time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	}
	if err != nil {
		return "", fmt.Errorf("[BeamPass.main.parseDate] time.ParseInLocation error: %w", err)
	}
	return tstmp.Format(layout), nil
}
