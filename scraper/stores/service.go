package stores

import (
	"encoding/json"
	"errors"
	"net/http"
)

// FetchStores fetches the store data from the provided URL and returns a parsed response
func FetchStores(url string) (StoresResponse, error) {
	// Perform the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return StoresResponse{}, err
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return StoresResponse{}, errors.New("error: status code " + resp.Status)
	}

	// Decode the response into the StoresResponse struct
	var storesResponse StoresResponse
	if err := json.NewDecoder(resp.Body).Decode(&storesResponse); err != nil {
		return StoresResponse{}, err
	}

	return storesResponse, nil
}
