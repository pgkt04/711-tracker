package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"scraper/fuel"
	"scraper/stores"
	"sort"
	"time"
)

//func old() {
//	ctx := context.Background()
//	client := models.GetClient()
//	fuelCollection := client.Collection("fuel")
//	storeCollection := client.Collection("stores")
//
//	// Fetch the stores using the stores package
//	storesResponse, err := stores.FetchStores()
//	if err != nil {
//		fmt.Println("Error fetching stores:", err)
//		return
//	}
//
//	eans := map[string]string{
//		"52": "Special Unleaded 91",
//		"53": "Special Diesel",
//		"57": "Special E10",
//		"56": "Supreme+ 98",
//		"55": "Extra 95",
//	}
//
//	// Struct to store fuel price information
//	type FuelInfo struct {
//		Timestamp         time.Time
//		StoreID           string
//		StoreName         string
//		Price             float64
//		PriceDate         time.Time
//		IsRecentlyUpdated bool
//	}
//
//	// Map to store top 5 cheapest prices for each EAN
//	cheapestPrices := make(map[string][]FuelInfo)
//
//	// Map to keep track of stores to be added to Firestore
//	storeData := make(map[string]map[string]interface{})
//
//	for _, store := range storesResponse.Stores {
//		fmt.Printf("Store ID: %s, Name: %s, Distance: %f\n", store.StoreId, store.Name, store.Distance)
//
//		// skip non-fuel stores.
//		if !store.IsFuelStore {
//			continue
//		}
//
//		fuelPrices, err := fuel.GetFuelPrices(store.StoreId)
//		if err != nil {
//			fmt.Printf("Failed to get fuel price for store ID %s: %v\n", store.StoreId, err)
//			continue
//		}
//
//		for _, price := range fuelPrices.Data {
//			ean := price.EAN
//			fuelInfo := FuelInfo{
//				StoreID:           store.StoreId,
//				StoreName:         store.Name,
//				Price:             float64(price.Price),
//				PriceDate:         price.PriceDate,
//				IsRecentlyUpdated: price.IsRecentlyUpdated,
//			}
//
//			// Update the top 5 cheapest prices list
//			cheapestPrices[ean] = append(cheapestPrices[ean], fuelInfo)
//
//			// Sort by price and maintain only the top 5
//			sort.SliceStable(cheapestPrices[ean], func(i, j int) bool {
//				return cheapestPrices[ean][i].Price < cheapestPrices[ean][j].Price
//			})
//
//			// Keep only the top 5 cheapest prices
//			if len(cheapestPrices[ean]) > 5 {
//				cheapestPrices[ean] = cheapestPrices[ean][:5]
//			}
//		}
//
//		// Store the store information for later Firestore addition
//		storeData[store.StoreId] = map[string]interface{}{
//			"storeId":     store.StoreId,
//			"name":        store.Name,
//			"distance":    store.Distance,
//			"isFuelStore": store.IsFuelStore,
//		}
//
//		// Tone it down.
//		time.Sleep(500 * time.Millisecond) // Adjust this duration based on the API's rate limit
//	}
//
//	// Store the top 5 cheapest prices in Firestore
//	for ean, infos := range cheapestPrices {
//		for _, info := range infos {
//			_, _, err := fuelCollection.Add(ctx, map[string]interface{}{
//				"timestamp":         time.Now(),
//				"ean":               ean,
//				"price":             info.Price,
//				"priceDate":         info.PriceDate,
//				"isRecentlyUpdated": info.IsRecentlyUpdated,
//				"storeNo":           info.StoreID,
//			})
//			if err != nil {
//				fmt.Printf("Failed to write fuel price to Firestore: %v \n", err)
//				continue
//			}
//			fmt.Println("Added entry!")
//		}
//	}
//
//	// Now, add store data to Firestore
//	for _, storeInfo := range storeData {
//		_, _, err := storeCollection.Add(ctx, storeInfo)
//		if err != nil {
//			fmt.Printf("Failed to write store info to Firestore: %v \n", err)
//		}
//	}
//
//	// Print the top 5 cheapest prices for each EAN
//	for ean, infos := range cheapestPrices {
//		fmt.Printf("EAN: %s (%s)\n", ean, eans[ean])
//		for _, info := range infos {
//			fmt.Printf("- Price: %.2f at Store ID: %s, Store Name: %s, Date: %s\n",
//				info.Price, info.StoreID, info.StoreName, info.PriceDate)
//		}
//	}
//}

func downloadData() {
	storesResponse, err := stores.FetchStores()
	if err != nil {
		fmt.Println("Error fetching stores:", err)
		return
	}

	var allStores []stores.Store
	var allFuelPrices []fuel.FuelPrice
	var retryQueue []stores.Store

	for _, store := range storesResponse.Stores {
		fmt.Printf("Store ID: %s, Name: %s, Distance: %f\n", store.StoreId, store.Name, store.Distance)
		if store.IsFuelStore {
			allStores = append(allStores, store)
			success := processStore(store, &allFuelPrices)
			if !success {
				retryQueue = append(retryQueue, store)
			}
		}
	}

	// Retry fetching fuel prices for stores in the retry queue
	for len(retryQueue) > 0 {
		retryQueue, allFuelPrices = retryFailedStores(retryQueue, allFuelPrices)
	}

	saveData(allStores, allFuelPrices)
}

func processStore(store stores.Store, allFuelPrices *[]fuel.FuelPrice) bool {
	fuelPrices, err := fuel.GetFuelPrices(store.StoreId)
	if err != nil {
		fmt.Printf("Failed to get fuel price for store ID %s: %v\n", store.StoreId, err)
		return false
	}

	*allFuelPrices = append(*allFuelPrices, fuelPrices.Data...)
	time.Sleep(500 * time.Millisecond) // Rate limiting
	return true
}

func retryFailedStores(retryQueue []stores.Store, allFuelPrices []fuel.FuelPrice) ([]stores.Store, []fuel.FuelPrice) {
	var nextRetryQueue []stores.Store

	for _, store := range retryQueue {
		fmt.Printf("Retrying store ID: %s\n", store.StoreId)
		success := processStore(store, &allFuelPrices)
		if !success {
			nextRetryQueue = append(nextRetryQueue, store)
		}
	}

	return nextRetryQueue, allFuelPrices
}

func saveData(allStores []stores.Store, allFuelPrices []fuel.FuelPrice) {
	timestamp := time.Now().Format("20060102-150405")
	fuelFileName := fmt.Sprintf("fuel-%s.json", timestamp)
	storeFileName := fmt.Sprintf("stores-%s.json", timestamp)

	if err := writeToFile(fuelFileName, allFuelPrices); err != nil {
		fmt.Printf("Failed to write fuel data to JSON file: %v\n", err)
	} else {
		fmt.Printf("Fuel data written to %s\n", fuelFileName)
	}

	if err := writeToFile(storeFileName, allStores); err != nil {
		fmt.Printf("Failed to write store data to JSON file: %v\n", err)
	} else {
		fmt.Printf("Store data written to %s\n", storeFileName)
	}
}

func writeToFile(filename string, data interface{}) error {
	filePath := filepath.Join(".", filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON to file %s: %v", filePath, err)
	}
	return nil
}

var eans = map[string]string{
	"52": "Special Unleaded 91",
	"53": "Special Diesel",
	"57": "Special E10",
	"56": "Supreme+ 98",
	"55": "Extra 95",
	//"94":
}

type StoreFuelData struct {
	Store stores.Store
	Price fuel.FuelPrice
}

func readStoresFromFile(filename string) []stores.Store {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Failed to open stores file: %v\n", err)
		return nil
	}
	defer file.Close()

	var storeList []stores.Store
	bytes, _ := io.ReadAll(file)
	if err := json.Unmarshal(bytes, &storeList); err != nil {
		fmt.Printf("Failed to parse stores JSON: %v\n", err)
		return nil
	}

	return storeList
}

func readFuelPricesFromFile(filename string) []fuel.FuelPrice {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Failed to open fuel file: %v\n", err)
		return nil
	}
	defer file.Close()

	var fuelPriceList []fuel.FuelPrice
	bytes, _ := io.ReadAll(file)
	if err := json.Unmarshal(bytes, &fuelPriceList); err != nil {
		fmt.Printf("Failed to parse fuel JSON: %v\n", err)
		return nil
	}

	return fuelPriceList
}

func parseCheapest() {
	storeFileName := "stores-20241018-155736.json" // Replace with actual path
	fuelFileName := "fuel-20241018-155736.json"    // Replace with actual path

	storesList := readStoresFromFile(storeFileName)
	fuelPrices := readFuelPricesFromFile(fuelFileName)

	storeMap := make(map[string]stores.Store)
	for _, store := range storesList {
		storeMap[store.StoreId] = store
	}

	type FuelStateData struct {
		State     string
		Price     float64
		StoreID   string
		StoreName string
		Suburb    string
	}

	stateEANMap := make(map[string]map[string][]FuelStateData)

	for _, price := range fuelPrices {
		store, exists := storeMap[price.StoreNo]
		if !exists {
			continue
		}
		ean := price.EAN
		state := store.Address.State

		if _, exists := stateEANMap[state]; !exists {
			stateEANMap[state] = make(map[string][]FuelStateData)
		}
		stateEANMap[state][ean] = append(stateEANMap[state][ean], FuelStateData{
			State:     state,
			Price:     float64(price.Price),
			StoreID:   store.StoreId,
			StoreName: store.Name,
			Suburb:    store.Address.Suburb,
		})
	}

	for state, eansData := range stateEANMap {
		for ean, prices := range eansData {
			sort.Slice(prices, func(i, j int) bool {
				return prices[i].Price < prices[j].Price
			})

			top3 := prices
			if len(prices) > 3 {
				top3 = prices[:3]
			}

			eanName := eans[ean]
			fmt.Printf("State: %s, EAN: %s (%s)\n", state, ean, eanName)
			for _, info := range top3 {
				fmt.Printf("- Price: %.2f, Store: %s, Suburb: %s\n", info.Price, info.StoreName, info.Suburb)
			}
			fmt.Println()
		}
	}
}

func main() {
	//downloadData()
	parseCheapest()
}
