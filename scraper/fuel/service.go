package fuel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Function to perform the GET request
func GetFuelPrices(storeNo string) (*FuelResponse, error) {
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
