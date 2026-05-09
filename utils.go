package main

import "fmt"

// helper function to compose user affiliation from UserInfo
// First Last, Affliation: Department, University, City, State Zip, Country
func composeAffiliation(user UserInfo) string {
	aff := fmt.Sprintf("%s %s, Affiliation: %s, %s, %s, %s, %s %s, %s",
		user.FirstName, user.LastName, user.Department,
		user.City, user.State, user.Zip, user.Country)

	return aff
}
