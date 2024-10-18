// firestore_client.go
package models

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"log"
)

// Declare a global Firestore client
var client *firestore.Client

// init function will run automatically when the package is imported
func init() {
	err := initializeFirestore()
	if err != nil {
		log.Fatalf("Error initializing Firestore: %v", err)
	}
}

// initializeFirestore sets up the Firestore client
func initializeFirestore() error {
	// Set up context
	ctx := context.Background()

	// Initialize Firestore client (Replace with the correct path to your credentials)
	var err error
	client, err = firestore.NewClient(ctx, "fuel-data-2e457",
		option.WithCredentialsFile("fuel-data-2e457-firebase-adminsdk-f32yg-a6bd68f23a.json"))
	if err != nil {
		return fmt.Errorf("failed to create Firestore client: %v", err)
	}

	fmt.Println("Firestore client initialized successfully!")
	return nil
}

// GetClient returns the Firestore client
func GetClient() *firestore.Client {
	return client
}
