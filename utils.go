package main

import (
	"fmt"
	"strings"
)

// helper to safely dereference string pointers
func strPtr(s *string) string {
	if s == nil {
		return ""
	}
	return strings.TrimSpace(*s)
}

// helper function to compose user affiliation from UserInfo
// First Last, Affiliation: University, City, State Zip, Country
func composeAffiliation(user UserInfo) string {
	parts := []string{}

	/*
		if v := strPtr(user.Department); v != "" {
			parts = append(parts, v)
		}
	*/

	if v := strPtr(user.Organization); v != "" {
		parts = append(parts, v)
	}

	location := strings.TrimSpace(
		fmt.Sprintf("%s, %s %s",
			strPtr(user.City),
			strPtr(user.State),
			strPtr(user.Zip),
		),
	)

	location = strings.Trim(location, ", ")

	if location != "" {
		parts = append(parts, location)
	}

	if v := strPtr(user.Country); v != "" {
		parts = append(parts, v)
	}

	affiliation := strings.Join(parts, ", ")

	return fmt.Sprintf(
		"%s %s, Affiliation: %s",
		user.FirstName,
		user.LastName,
		affiliation,
	)
}
