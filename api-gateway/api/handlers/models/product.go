package models

type Product struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Prays  int64  `json:"prays"`
	Amount int64  `json:"amount"`
}
