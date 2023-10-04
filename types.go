package main

type CreateAccountRequest struct {
	Prefix           string `json:"prefix"`
	Name             string `json:"name"`
	Affiliation      string `json:"affiliation"` // DB account table to be changed to match (institution -> affiliation)
	BH               string `json:"bh"`
	Available        string `json:"available"`
	Rating           string `json:"rating"`
	Role             string `json:"role"`
	Type             string `json:"type"`
	GeneralArea      string `json:"generalArea"`
	SpecialisedArea  string `json:"specialisedArea"`
	Trained          string `json:"trained"`
	PrimaryContact   string `json:"primaryContact"`
	SecondaryContact string `json:"secondaryContact"`
	Email            string `json:"email"`
	Published        string `json:"published"`
}
type Account struct {
	ID               int    `json:"id"`
	Prefix           string `json:"prefix"`
	Name             string `json:"name"`
	Affiliation      string `json:"affiliation"` // DB account table to be changed to match (institution -> affiliation)
	BH               string `json:"bh"`
	Available        string `json:"available"`
	Rating           string `json:"rating"`
	Role             string `json:"role"`
	Type             string `json:"type"`
	GeneralArea      string `json:"generalArea"`
	SpecialisedArea  string `json:"specialisedArea"`
	Trained          string `json:"trained"`
	PrimaryContact   string `json:"primaryContact"`
	SecondaryContact string `json:"secondaryContact"`
	Email            string `json:"email"`
	Published        string `json:"published"`
}

func NewAccount(prefix,
	name,
	affiliation,
	bh,
	available,
	rating,
	role,
	expType,
	generalArea,
	specialisedArea,
	trained,
	primaryContact,
	secondaryContact,
	email,
	published string) *Account {
	return &Account{
		Prefix:           prefix,
		Name:             name,
		Affiliation:      affiliation,
		BH:               bh,
		Available:        available,
		Rating:           rating,
		Role:             role,
		Type:             expType,
		GeneralArea:      generalArea,
		SpecialisedArea:  specialisedArea,
		Trained:          trained,
		PrimaryContact:   primaryContact,
		SecondaryContact: secondaryContact,
		Email:            email,
		Published:        published,
	}
}
