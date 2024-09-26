package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Define the Go structs that match the JSON structure
type Address struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Suburb   string `json:"suburb"`
	State    string `json:"state"`
	Postcode string `json:"postcode"`
	Extra    string `json:"extra"`
}

type OpeningHours struct {
	DayOfWeek int    `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type Region struct {
	CountryId string `json:"countryId"`
	RegionId  int    `json:"regionId"`
	Region    string `json:"region"`
}

type Store struct {
	StoreId                string         `json:"storeId"`
	Distance               float64        `json:"distance"`
	Name                   string         `json:"name"`
	Location               []float64      `json:"location"`
	Centre                 string         `json:"centre"`
	Address                Address        `json:"address"`
	Phone                  string         `json:"phone"`
	Fax                    string         `json:"fax"`
	AllHours               bool           `json:"allHours"`
	IsActive               bool           `json:"isActive"`
	IsDigitalDisplay       bool           `json:"isDigitalDisplay"`
	IsFuelStore            bool           `json:"isFuelStore"`
	HasKiosk               bool           `json:"hasKiosk"`
	Features               []string       `json:"features"`
	ParcelMate             []string       `json:"ParcelMate"`
	FuelOptions            []string       `json:"fuelOptions"`
	Atm                    bool           `json:"atm"`
	IsBrandNewStore        bool           `json:"isBrandNewStore"`
	IsFranchiseOpp         bool           `json:"isFranchiseOpp"`
	FranchiseSuburb        string         `json:"franchiseSuburb"`
	FranchiseEstimatedCost string         `json:"franchiseEstimatedCost"`
	AllowStoreDelivery     bool           `json:"allowStoreDelivery"`
	OpeningHours           []OpeningHours `json:"openingHours"`
	SpecialOpeningHours    []string       `json:"specialOpeningHours"`
	TimeSlot               string         `json:"timeSlot"`
	Region                 Region         `json:"region"`
	HideFromStorelocator   bool           `json:"hideFromStorelocator"`
}

type StoresResponse struct {
	Stores []Store `json:"stores"`
}

func main() {
	url := "https://www.7eleven.com.au/storelocator-retail/mulesoft/stores?lat=-33.8688197&long=151.2092955&dist=10" // Replace with the actual URL

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		return
	}

	// Decode the response into the Go struct
	var storesResponse StoresResponse
	if err := json.NewDecoder(resp.Body).Decode(&storesResponse); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	// Print the decoded data (for debugging purposes)
	for _, store := range storesResponse.Stores {
		fmt.Printf("Store ID: %s, Name: %s, Distance: %f\n", store.StoreId, store.Name, store.Distance)
	}
}
