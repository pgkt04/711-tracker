package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

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

func fetchStores() (*StoresResponse, error) {
	url := "https://www.7eleven.com.au/storelocator-retail/mulesoft/stores" // Replace with the actual URL

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: status code %d", resp.StatusCode)
	}

	// Decode the response into the Go struct
	var storesResponse StoresResponse
	if err := json.NewDecoder(resp.Body).Decode(&storesResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &storesResponse, nil
}

func test_stores() {
	storesResponse, err := fetchStores()
	if err != nil {
		fmt.Printf("Error fetching stores: %v\n", err)
		return
	}

	// Use the fetched store data
	for _, store := range storesResponse.Stores {
		if !store.IsFuelStore {
			continue
		}
		fmt.Printf("Store ID: %s, Name: %s, Distance: %f\n", store.StoreId, store.Name, store.Distance)
	}
}

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

// Function to perform the GET request
func getFuelPrices(storeNo string) (*FuelResponse, error) {
	// Create the URL with the provided store number
	url := fmt.Sprintf("https://www.7eleven.com.au/storelocator-retail/mulesoft/fuelPrices?storeNo=%s", storeNo)

	// Perform the GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check for a successful response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Read the response body using io.ReadAll (replaces ioutil.ReadAll)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON response
	var result FuelResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Return the parsed response
	return &result, nil
}

func test_fuel() {
	// Example store number
	storeNo := "2362"

	// Get fuel prices for the store
	fuelPrices, err := getFuelPrices(storeNo)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Handle the case where no data is returned
	if len(fuelPrices.Data) == 0 {
		fmt.Println("No fuel prices available.")
		return
	}

	// Print the fuel prices
	for _, price := range fuelPrices.Data {
		fmt.Printf("EAN: %s, Price: %d, Price Date: %s, Recently Updated: %t\n",
			price.EAN, price.Price, price.PriceDate.Format(time.RFC1123), price.IsRecentlyUpdated)
	}
}

func main() {
	test_stores()
	// test_fuel()
}
