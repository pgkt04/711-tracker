package main

import (
	"context"
	"fmt"
	"scraper/fuel"
	"scraper/models"
	"scraper/stores"
)

func main() {
	ctx := context.Background()
	client := models.GetClient()
	collection := client.Collection("fuel")

	// Fetch the stores using the stores package
	storesResponse, err := stores.FetchStores()
	if err != nil {
		fmt.Println("Error fetching stores:", err)
		return
	}

	for _, store := range storesResponse.Stores {
		fmt.Printf("Store ID: %s, Name: %s, Distance: %f\n", store.StoreId, store.Name, store.Distance)
		if store.IsFuelStore {
			fuelPrices, err := fuel.GetFuelPrices(store.StoreId)
			if err != nil {
				fmt.Printf("Failed to get fuel price for store ID %s: %v\n", store.StoreId, err)
				return
			}
			for _, price := range fuelPrices.Data {
				_, _, err := collection.Add(ctx, map[string]interface{}{
					"ean":               price.EAN,
					"price":             price.Price,
					"priceDate":         price.PriceDate,
					"isRecentlyUpdated": price.IsRecentlyUpdated,
					"storeNo":           price.StoreNo,
				})
				if err != nil {
					fmt.Printf("Failed to write fuel price to Firestore: %v \n", err)
				}
			}
		}
	}
}
