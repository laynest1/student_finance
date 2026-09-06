package models

type Transaction struct {
	ID int `json:"id"`
	Amount float64 `json:"amount"`
	Category string `json:"category"`
	Date string `json:"date"`
	Description string `json:"description"`

}