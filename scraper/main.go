package main

import (
	"context"
	"fmt"
	"scraper/fuel"
	"scraper/models"
	"scraper/stores"
	"sort"
	"time"
)

//func writeToFile(filename string, data interface{}) error {
//	filePath := filepath.Join(".", filename)
//	file, err := os.Create(filePath)
//	if err != nil {
//		return fmt.Errorf("failed to create file %s: %v", filePath, err)
//	}
//	defer file.Close()
//
//	encoder := json.NewEncoder(file)
//	encoder.SetIndent("", "  ")
//	if err := encoder.Encode(data); err != nil {
//		return fmt.Errorf("failed to encode JSON to file %s: %v", filePath, err)
//	}
//	return nil
//}
//
//	func downloadData() {
//		storesResponse, err := stores.FetchStores()
//		if err != nil {
//			fmt.Println("Error fetching stores:", err)
//			return
//		}
//
//		var allStores []stores.Store
//		var allFuelPrices []fuel.FuelPrice
//		var retryQueue []stores.Store
//
//		for _, store := range storesResponse.Stores {
//			fmt.Printf("Store ID: %s, Name: %s, Distance: %f\n", store.StoreId, store.Name, store.Distance)
//			if store.IsFuelStore {
//				allStores = append(allStores, store)
//				success := processStore(store, &allFuelPrices)
//				if !success {
//					retryQueue = append(retryQueue, store)
//				}
//			}
//		}
//
//		// Retry fetching fuel prices for stores in the retry queue
//		for len(retryQueue) > 0 {
//			retryQueue, allFuelPrices = retryFailedStores(retryQueue, allFuelPrices)
//		}
//
//		saveData(allStores, allFuelPrices)
//	}
//func readStoresFromFile(filename string) []stores.Store {
//	file, err := os.Open(filename)
//	if err != nil {
//		fmt.Printf("Failed to open stores file: %v\n", err)
//		return nil
//	}
//	defer file.Close()
//
//	var storeList []stores.Store
//	bytes, _ := io.ReadAll(file)
//	if err := json.Unmarshal(bytes, &storeList); err != nil {
//		fmt.Printf("Failed to parse stores JSON: %v\n", err)
//		return nil
//	}
//
//	return storeList
//}
//
//func readFuelPricesFromFile(filename string) []fuel.FuelPrice {
//	file, err := os.Open(filename)
//	if err != nil {
//		fmt.Printf("Failed to open fuel file: %v\n", err)
//		return nil
//	}
//	defer file.Close()
//
//	var fuelPriceList []fuel.FuelPrice
//	bytes, _ := io.ReadAll(file)
//	if err := json.Unmarshal(bytes, &fuelPriceList); err != nil {
//		fmt.Printf("Failed to parse fuel JSON: %v\n", err)
//		return nil
//	}
//
//	return fuelPriceList
//}
//
//func saveData(allStores []stores.Store, allFuelPrices []fuel.FuelPrice) {
//	timestamp := time.Now().Format("20060102-150405")
//	fuelFileName := fmt.Sprintf("fuel-%s.json", timestamp)
//	storeFileName := fmt.Sprintf("stores-%s.json", timestamp)
//
//	if err := writeToFile(fuelFileName, allFuelPrices); err != nil {
//		fmt.Printf("Failed to write fuel data to JSON file: %v\n", err)
//	} else {
//		fmt.Printf("Fuel data written to %s\n", fuelFileName)
//	}
//
//	if err := writeToFile(storeFileName, allStores); err != nil {
//		fmt.Printf("Failed to write store data to JSON file: %v\n", err)
//	} else {
//		fmt.Printf("Store data written to %s\n", storeFileName)
//	}
//}
//
//func parseCheapest() {
//	ctx := context.Background()
//	client := models.GetClient()
//	// storeCollection := client.Collection("store") // not needed?
//	fuelCollection := client.Collection("fuel")
//
//	storeFileName := "stores-20241022-211146.json"
//	fuelFileName := "fuel-20241022-211146.json"
//
//	storesList := readStoresFromFile(storeFileName)
//	fuelPrices := readFuelPricesFromFile(fuelFileName)
//
//	currentTime := time.Now()
//
//	storeMap := make(map[string]stores.Store)
//	for _, store := range storesList {
//		storeMap[store.StoreId] = store
//	}
//
//	type FuelStateData struct {
//		State     string
//		Price     float64
//		PriceDate time.Time
//		StoreID   string
//		StoreName string
//		Suburb    string
//		Address   string
//		Postcode  string
//	}
//
//	stateEANMap := make(map[string]map[string][]FuelStateData)
//
//	for _, price := range fuelPrices {
//		store, exists := storeMap[price.StoreNo]
//		if !exists {
//			continue
//		}
//		ean := price.EAN
//		state := store.Address.State
//
//		if _, exists := stateEANMap[state]; !exists {
//			stateEANMap[state] = make(map[string][]FuelStateData)
//		}
//		stateEANMap[state][ean] = append(stateEANMap[state][ean], FuelStateData{
//			State:     state,
//			Price:     float64(price.Price),
//			PriceDate: price.PriceDate,
//			StoreID:   store.StoreId,
//			StoreName: store.Name,
//			Suburb:    store.Address.Suburb,
//			Address:   store.Address.Address1, // address 2 usually empty
//			Postcode:  store.Address.Postcode,
//		})
//	}
//
//	for state, eansData := range stateEANMap {
//		for ean, prices := range eansData {
//			sort.Slice(prices, func(i, j int) bool {
//				return prices[i].Price < prices[j].Price
//			})
//
//			top3 := prices
//			if len(prices) > 3 {
//				top3 = prices[:3]
//			}
//
//			eanName := eans[ean]
//			fmt.Printf("State: %s, EAN: %s (%s)\n", state, ean, eanName)
//			for _, info := range top3 {
//				fmt.Printf("- Price: %.2f, Store: %s, Suburb: %s\n", info.Price, info.StoreName, info.Suburb)
//				_, _, err := fuelCollection.Add(ctx, map[string]interface{}{
//					"time":      currentTime,
//					"storeID":   info.StoreID,
//					"ean":       ean,
//					"price":     info.Price,
//					"priceDate": info.PriceDate,
//					"state":     info.State,
//					"suburb":    info.Suburb,
//					"address":   info.Address,
//					"postcode":  info.Postcode,
//				})
//
//				if err != nil {
//					fmt.Println("Failed to insert into firestore")
//				}
//			}
//			fmt.Println()
//		}
//	}
//	fmt.Println("Done!")
//}

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

var eans = map[string]string{
	"52": "Special Unleaded 91",
	"53": "Special Diesel",
	"57": "Special E10",
	"56": "Supreme+ 98",
	"55": "Extra 95",
	"54": "LPG",
}

// combines downloading data and parsing it without writing to disk
func downloadAndParseCheapest() {
	// Step 1: Fetch Stores
	storesResponse, err := stores.FetchStores()
	if err != nil {
		fmt.Println("Error fetching stores:", err)
		return
	}

	var allStores []stores.Store
	var allFuelPrices []fuel.FuelPrice
	var retryQueue []stores.Store

	// Step 2: Process Each Store to Fetch Fuel Prices
	for _, store := range storesResponse.Stores {
		fmt.Printf("Store ID: %s, Name: %s \n", store.StoreId, store.Name)
		if store.IsFuelStore {
			allStores = append(allStores, store)
			success := processStore(store, &allFuelPrices)
			if !success {
				retryQueue = append(retryQueue, store)
			}
		}
	}

	// Step 3: Retry Fetching Fuel Prices for Failed Stores
	for len(retryQueue) > 0 {
		retryQueue, allFuelPrices = retryFailedStores(retryQueue, allFuelPrices)
	}

	// Step 4: Process Fuel Prices to Find Top 3 Cheapest per EAN per State
	ctx := context.Background()
	client := models.GetClient()
	fuelCollection := client.Collection("fuel")
	currentTime := time.Now()

	// Create a map for quick store lookup
	storeMap := make(map[string]stores.Store)
	for _, store := range allStores {
		storeMap[store.StoreId] = store
	}

	type FuelStateData struct {
		State     string
		Price     float64
		PriceDate time.Time
		StoreID   string
		StoreName string
		Suburb    string
		Address   string
		Postcode  string
	}

	stateEANMap := make(map[string]map[string][]FuelStateData)

	// Organize fuel prices by state and EAN
	for _, price := range allFuelPrices {
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
			PriceDate: price.PriceDate,
			StoreID:   store.StoreId,
			StoreName: store.Name,
			Suburb:    store.Address.Suburb,
			Address:   store.Address.Address1, // address 2 usually empty
			Postcode:  store.Address.Postcode,
		})
	}

	// Iterate through each state and EAN to find top 3 cheapest prices
	for state, eansData := range stateEANMap {
		for ean, prices := range eansData {
			// Sort prices in ascending order
			sort.Slice(prices, func(i, j int) bool {
				return prices[i].Price < prices[j].Price
			})

			// Select top 3 cheapest prices
			top3 := prices
			if len(prices) > 3 {
				top3 = prices[:3]
			}

			eanName := eans[ean]
			fmt.Printf("State: %s, EAN: %s (%s)\n", state, ean, eanName)
			for _, info := range top3 {
				fmt.Printf("- Price: %.2f, Store: %s, Suburb: %s\n", info.Price, info.StoreName, info.Suburb)
				_, _, err := fuelCollection.Add(ctx, map[string]interface{}{
					"time":      currentTime,
					"storeID":   info.StoreID,
					"ean":       ean,
					"price":     info.Price,
					"priceDate": info.PriceDate,
					"state":     info.State,
					"suburb":    info.Suburb,
					"address":   info.Address,
					"postcode":  info.Postcode,
				})

				if err != nil {
					fmt.Println("Failed to insert into Firestore:", err)
				}
			}
			fmt.Println()
		}
	}
	fmt.Println("Download and Parsing Completed Successfully!")
}

func main() {
	//downloadData()
	//parseCheapest()
	downloadAndParseCheapest()
}
