package stores

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
