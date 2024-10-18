package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"scraper/stores"
)

func test() {
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

func main() {
	// Set up context and Firestore client
	ctx := context.Background()

	// Create a Firestore client by providing your credentials file path (if needed)
	// Replace "path-to-your-service-account-key.json" with your actual service account key file
	// If you're running this on Google Cloud, it will pick up the credentials automatically
	client, err := firestore.NewClient(ctx, "fuel-data-2e457", option.WithCredentialsFile("path-to-your-service-account-key.json"))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	// Define the data you want to write
	docData := map[string]interface{}{
		"name":    "John Doe",
		"age":     30,
		"email":   "johndoe@example.com",
		"address": "123 Main St, Anytown, USA",
	}

	// Write data to a document in a collection
	_, err = client.Collection("users").Doc("user-1").Set(ctx, docData)
	if err != nil {
		log.Fatalf("Failed to write data to Firestore: %v", err)
	}

	fmt.Println("Data written successfully!")
}
