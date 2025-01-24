package services

import (
	"testing"

	"github.com/Mekambee/Swift-Codes-Api/internal/database"
	"github.com/Mekambee/Swift-Codes-Api/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSwiftService_Basic(t *testing.T) {
	err := database.ConnectAndMigrate("localhost", "myuser", "mysecretpassword", "swiftdb", "5432")
	assert.NoError(t, err)
	defer database.DB.Close()

	database.DB.Exec("TRUNCATE swift_codes RESTART IDENTITY")

	records := []models.SwiftCodeData{
		{
			SwiftCode:     "TESTPLW1XXX",
			BankName:      "TEST BANK",
			Address:       "Some Address",
			CountryISO2:   "PL",
			CountryName:   "POLAND",
			IsHeadquarter: true,
		},
		{
			SwiftCode:     "TESTPLW1ABC",
			BankName:      "TEST BANK",
			Address:       "Branch Address",
			CountryISO2:   "PL",
			CountryName:   "POLAND",
			IsHeadquarter: false,
		},
	}
	err = SaveSwiftCodes(records)
	assert.NoError(t, err, "Should save without error")

	sc, err := GetSwiftCode("TESTPLW1XXX")
	assert.NoError(t, err)
	assert.Equal(t, "TESTPLW1XXX", sc.SwiftCode)
	assert.True(t, sc.IsHeadquarter)

	branches, err := GetBranchesByHQ("TESTPLW1XXX")
	assert.NoError(t, err)
	assert.Len(t, branches, 1, "Should find one branch matching first 8 chars")

	err = DeleteSwiftCode("TESTPLW1XXX", "TEST BANK", "PL")
	assert.NoError(t, err, "Should delete HQ record")

	_, err = GetSwiftCode("TESTPLW1XXX")
	assert.Error(t, err, "Should not find the HQ after deletion")
}
