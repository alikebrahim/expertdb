package main

import (
	"time"
)

type CreateAccountRequest struct {
	Name           string `json:"name"`
	Affiliation    string `json:"affiliation"`
	PrimaryContact string `json:"primaryContact"`
	This is a test
}
type Account struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Affiliation    string    `json:"affiliation"`
	PrimaryContact string    `json:"primaryContact"`
	CreatedAt      time.Time `json:"created_at"`
}

func NewAccount(name, affiliation, primaryContact string) *Account {
	return &Account{
		Name:           name,
		Affiliation:    affiliation,
		PrimaryContact: primaryContact,
		CreatedAt:      time.Now().UTC(),
	}
}
