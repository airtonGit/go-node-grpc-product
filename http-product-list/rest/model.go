package rest

import "time"

type productDiscount struct {
	Percent      float32 `json:"pct"`
	ValueInCents int32   `json:"value_in_cents"`
}

type product struct {
	ID           string          `json:"id"`
	PriceInCents int32           `json:"price_in_cents"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Discount     productDiscount `json:"discount"`
}

type user struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	DateOfBirth time.Time `json:"date_of_birth"`
}
