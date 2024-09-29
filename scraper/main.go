package main

import (
	"fmt"
	"scraper/stores"
)

func main() {
	// Define the store locator API endpoint (replace with the actual URL if needed)
	url := "https://www.7eleven.com.au/storelocator-retail/mulesoft/stores?lat=-33.8688197&long=151.2092955&dist=10"

	// Fetch the stores using the stores package
	storesResponse, err := stores.FetchStores(url)
	if err != nil {
		fmt.Println("Error fetching stores:", err)
		return
	}

	// Print the store details
	for _, store := range storesResponse.Stores {
		fmt.Printf("Store ID: %s, Name: %s, Distance: %f\n", store.StoreId, store.Name, store.Distance)
	}
}
