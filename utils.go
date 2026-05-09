package main

import "fmt"

// helper function to compose user affiliation from UserInfo
// First Last, Affliation: Department, University, City, State Zip, Country
func composeAffiliation(user UserInfo) string {
	dep := user.Department
	city := user.City
	state := user.State
	zip := user.Zip
	country := user.Country
	aff := fmt.Sprintf("%s %s, Affiliation: %s, %s, %s, %s, %s %s, %s",
		user.FirstName, user.LastName, *dep, *city, *state, *zip, *country)
	return aff
}
