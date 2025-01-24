package services

import (
	"strings"

	"github.com/Mekambee/Swift-Codes-Api/internal/models"
	"github.com/xuri/excelize/v2"
)

func ParseSwiftXLSX(filePath string) ([]models.SwiftCodeData, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	var result []models.SwiftCodeData

	for i := 1; i < len(rows); i++ {
		row := rows[i]
		if len(row) < 8 {
			continue
		}

		iso2 := strings.ToUpper(row[0])
		swift := row[1]
		bankName := row[3]
		address := row[4]
		town := row[5]
		countryName := strings.ToUpper(row[6])
		fullAddress := address + ", " + town

		isHQ := false
		if strings.HasSuffix(strings.ToUpper(swift), "XXX") {
			isHQ = true
		}

		data := models.SwiftCodeData{
			SwiftCode:     swift,
			BankName:      bankName,
			Address:       fullAddress,
			CountryISO2:   iso2,
			CountryName:   countryName,
			IsHeadquarter: isHQ,
		}

		result = append(result, data)
	}

	return result, nil
}
