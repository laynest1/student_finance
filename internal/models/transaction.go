package models

type Transaction struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	Amount float64 `json:"amount"`
	Category string `json:"category"`
	Date string `json:"date"`
	Description string `json:"description"`

}