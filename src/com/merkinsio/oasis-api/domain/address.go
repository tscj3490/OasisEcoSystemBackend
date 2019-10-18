package domain

// Address Holds the user location
type Address struct {
	Street   string `json:"street" bson:"street"`
	ZipCode  string `json:"zipCode" bson:"zipCode"`
	City     string `json:"city" bson:"city"`
	Province string `json:"province" bson:"province"`
}
