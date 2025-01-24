package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Mekambee/Swift-Codes-Api/internal/api"
	"github.com/Mekambee/Swift-Codes-Api/internal/database"
	"github.com/Mekambee/Swift-Codes-Api/internal/services"
)

// Tests GET /v1/swift-codes/{swiftCode} after loading XLSX data into the DB.
func TestIntegration_GetSwiftCode(t *testing.T) {
	err := database.ConnectAndMigrate("localhost", "myuser", "mysecretpassword", "swiftdb", "5432")
	assert.NoError(t, err, "DB connection should succeed")
	defer database.DB.Close()

	_, _ = database.DB.Exec("TRUNCATE swift_codes RESTART IDENTITY")

	testFile := filepath.Join("testdata", "integration_test.xlsx")
	data, err := services.ParseSwiftXLSX(testFile)
	assert.NoError(t, err, "Parsing XLSX should succeed")

	err = services.SaveSwiftCodes(data)
	assert.NoError(t, err, "Saving to DB should succeed")

	router := api.SetupRouter()

	swiftCode := "AIPOPLP1XXX"
	req, _ := http.NewRequest("GET", "/v1/swift-codes/"+swiftCode, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK")

	bodyBytes, _ := io.ReadAll(w.Body)
	var resp map[string]interface{}
	err = json.Unmarshal(bodyBytes, &resp)
	assert.NoError(t, err, "Response should be valid JSON")

	assert.Equal(t, swiftCode, resp["swiftCode"], "Swift code should match")
	assert.Equal(t, true, resp["isHeadquarter"], "Should be HQ (ends with XXX)")
	assert.Equal(t, "PL", resp["countryISO2"], "countryISO2 should match the XLSX data")

	branchesVal, ok := resp["branches"]
	assert.True(t, ok, "HQ response should have a 'branches' field")
	branchesArr, ok := branchesVal.([]interface{})
	assert.True(t, ok, "'branches' should be an array")
	_ = branchesArr

	t.Logf("HQ response: %v", resp)
}

// Tests GET /v1/swift-codes/country/{countryISO2}.
func TestIntegration_GetCountry(t *testing.T) {
	err := database.ConnectAndMigrate("localhost", "myuser", "mysecretpassword", "swiftdb", "5432")
	assert.NoError(t, err)
	defer database.DB.Close()

	_, _ = database.DB.Exec("TRUNCATE swift_codes RESTART IDENTITY")

	testFile := filepath.Join("testdata", "integration_test.xlsx")
	data, err := services.ParseSwiftXLSX(testFile)
	assert.NoError(t, err, "Parsing XLSX should succeed")
	err = services.SaveSwiftCodes(data)
	assert.NoError(t, err, "Saving to DB should succeed")

	router := api.SetupRouter()

	req, _ := http.NewRequest("GET", "/v1/swift-codes/country/CL", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK")

	var resp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err, "Response should be valid JSON")

	assert.Equal(t, "CL", resp["countryISO2"], "Should return data for CL")

	swiftCodesVal, ok := resp["swiftCodes"].([]interface{})
	assert.True(t, ok, "'swiftCodes' should be an array")
	assert.True(t, len(swiftCodesVal) > 0, "Expected at least one record for CL")
}

// Tests POST /v1/swift-codes/ for creating a new SWIFT code.
func TestIntegration_PostNewSwiftCode(t *testing.T) {
	err := database.ConnectAndMigrate("localhost", "myuser", "mysecretpassword", "swiftdb", "5432")
	assert.NoError(t, err)
	defer database.DB.Close()

	_, _ = database.DB.Exec("TRUNCATE swift_codes RESTART IDENTITY")

	router := api.SetupRouter()

	body := `{
		"address": "Some address",
		"bankName": "Test Bank",
		"countryISO2": "us",
		"countryName": "united states",
		"isHeadquarter": true,
		"swiftCode": "TESTUSNYXXX"
	}`

	req, _ := http.NewRequest("POST", "/v1/swift-codes/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK")

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Swift code created", resp["message"], "Should confirm creation")

	row := database.DB.QueryRow(`
		SELECT swift_code, country_iso2, is_headquarter
		FROM swift_codes
		WHERE swift_code = 'TESTUSNYXXX'
	`)
	var sc, iso2 string
	var hq bool
	err = row.Scan(&sc, &iso2, &hq)
	assert.NoError(t, err)
	assert.Equal(t, "TESTUSNYXXX", sc)
	assert.Equal(t, "US", iso2)
	assert.True(t, hq)
}

// Tests DELETE /v1/swift-codes/{swiftCode} with bankName and countryISO2 in the JSON body.
func TestIntegration_DeleteSwiftCode(t *testing.T) {
	err := database.ConnectAndMigrate("localhost", "myuser", "mysecretpassword", "swiftdb", "5432")
	assert.NoError(t, err)
	defer database.DB.Close()

	_, _ = database.DB.Exec("TRUNCATE swift_codes RESTART IDENTITY")

	_, _ = database.DB.Exec(`
		INSERT INTO swift_codes (swift_code, bank_name, address, country_iso2, country_name, is_headquarter)
		VALUES ('TESTUSNYXXX', 'TEST BANK', 'Some address', 'US', 'UNITED STATES', true)
	`)

	router := api.SetupRouter()

	body := `{"bankName":"TEST BANK","countryISO2":"US"}`
	req, _ := http.NewRequest("DELETE", "/v1/swift-codes/TESTUSNYXXX", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK")

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Swift code deleted", resp["message"])

	row := database.DB.QueryRow(`
		SELECT COUNT(*) FROM swift_codes
		WHERE swift_code = 'TESTUSNYXXX'
	`)
	var count int
	err = row.Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count, "Record should be deleted")
}
