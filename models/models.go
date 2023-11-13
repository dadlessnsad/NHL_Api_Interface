package models

type AllTeamsResponse struct {
	Teams []Team `json:"teams"`
}

type Team struct {
	Name   string         `json:"name"`
	Abbrev string         `json:"abbreviation"`
	Roster RosterResponse `json:"roster"`
}

type RosterResponse struct {
	Forwards   []Player `json:"forwards"`
	Defensemen []Player `json:"defensemen"`
	Goalies    []Player `json:"goalies"`
}

type Player struct {
	ID                  int    `json:"id"`
	Headshot            string `json:"headshot"`
	FirstName           Name   `json:"firstName"`
	FullName            string `json:"fullName"`
	LastName            Name   `json:"lastName"`
	SweaterNumber       int    `json:"sweaterNumber"`
	PositionCode        string `json:"positionCode"`
	ShootsCatches       string `json:"shootsCatches"`
	HeightInInches      int    `json:"heightInInches"`
	WeightInPounds      int    `json:"weightInPounds"`
	HeightInCentimeters int    `json:"heightInCentimeters"`
	WeightInKilograms   int    `json:"weightInKilograms"`
	BirthDate           string `json:"birthDate"`
	BirthCity           Name   `json:"birthCity"`
	BirthCountry        string `json:"birthCountry"`
}

type Name struct {
	Default string `json:"default"`
}
