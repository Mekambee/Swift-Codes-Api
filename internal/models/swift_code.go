package models

type SwiftCodeData struct {
	ID            int64  `json:"-"`
	SwiftCode     string `json:"swiftCode"`
	BankName      string `json:"bankName"`
	Address       string `json:"address"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
}
