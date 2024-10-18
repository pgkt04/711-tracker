package fuel

import "time"

// Define the structure to match the JSON response
type FuelPrice struct {
	EAN               string    `json:"ean"`
	Price             int       `json:"price"`
	PriceDate         time.Time `json:"priceDate"`
	IsRecentlyUpdated bool      `json:"isRecentlyUpdated"`
	StoreNo           string    `json:"storeNo"`
}

type FuelResponse struct {
	Data []FuelPrice `json:"data"`
}
