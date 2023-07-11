package config

import "go.mongodb.org/mongo-driver/bson/primitive"

// Propane constants
const (
	PropaneGradeID   = 6
	PropaneStationID = "56cf1815982d82b0f3000012"
)

func PropaneDispensers() []primitive.ObjectID {
	dispenser1, _ := primitive.ObjectIDFromHex("6475f6e88072bce37f4f57b6")
	dispenser2, _ := primitive.ObjectIDFromHex("6475f7028072bce37f4f57b7")
	return []primitive.ObjectID{dispenser1, dispenser2}
}

// dispenser1, _ := primitive.ObjectIDFromHex("6475f6e88072bce37f4f57b6")
// var PropaneDispensers = []primitive.ObjectID{primitive.ObjectIDFromHex("6475f6e88072bce37f4f57b6"), primitive.ObjectIDFromHex("6475f7028072bce37f4f57b7")}

// PropaneTankLookup function
func PropaneTankLookup(dispenserID string) int {
	// m := map[string]int{"56e7593f982d82eeff262cd5": 475, "56e7593f982d82eeff262cd6": 476} // old set
	m := map[string]int{"6475f6e88072bce37f4f57b6": 475, "6475f7028072bce37f4f57b7": 476}
	return m[dispenserID]
}
