package config

// Propane constants
const (
	PropaneGradeID   = 6
	PropaneStationID = "56cf1815982d82b0f3000012"
)

// PropaneTankLookup function
func PropaneTankLookup(dispenserID string) int {
	m := map[string]int{"56e7593f982d82eeff262cd5": 475, "56e7593f982d82eeff262cd6": 476}
	return m[dispenserID]
}
