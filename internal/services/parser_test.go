package services

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests if the parser correctly reads a row without "XXX" (branch).
func TestParseSwiftXLSX_Basic(t *testing.T) {
	testFile := filepath.Join("testdata", "test_basic.xlsx")

	records, err := ParseSwiftXLSX(testFile)
	assert.NoError(t, err, "Parsing should not return an error")
	assert.NotEmpty(t, records, "Expected at least one record")

	first := records[0]
	assert.Equal(t, "ABCDPLW1", first.SwiftCode, "SWIFT code should match")
	assert.False(t, first.IsHeadquarter, "Should be branch if no 'XXX' suffix")
	assert.Equal(t, "PL", first.CountryISO2, "ISO2 should be 'PL'")
	assert.Equal(t, "POLAND", first.CountryName, "Country name should be 'POLAND'")
	assert.Contains(t, first.Address, "ADDRESS1", "Address should contain 'ADDRESS1'")
	assert.Contains(t, first.Address, "CITY1", "Address should contain 'CITY1'")
}

// Tests if the parser recognizes "XXX" at the end as headquarter.
func TestParseSwiftXLSX_HQ(t *testing.T) {
	testFile := filepath.Join("testdata", "test_hq.xlsx")

	records, err := ParseSwiftXLSX(testFile)
	assert.NoError(t, err, "Parsing should succeed")
	assert.NotEmpty(t, records, "Expected at least one record")

	first := records[0]
	assert.Equal(t, "ABIEBGS1XXX", first.SwiftCode, "SWIFT code should match")
	assert.True(t, first.IsHeadquarter, "SWIFT ending with 'XXX' should be HQ")
	assert.Equal(t, "BG", first.CountryISO2, "ISO2 should be 'BG'")
	assert.Equal(t, "BULGARIA", first.CountryName, "Country name should be 'BULGARIA'")
}

// Tests if an entirely empty file returns no records.
func TestParseSwiftXLSX_CompletelyEmptyFile(t *testing.T) {
	testFile := filepath.Join("testdata", "test_completely_empty.xlsx")

	records, err := ParseSwiftXLSX(testFile)
	assert.NoError(t, err, "Parsing an empty file should not crash")
	assert.Empty(t, records, "No records expected for an empty file")
}

// Tests if a file with only headers returns no records.
func TestParseSwiftXLSX_OnlyHeaders(t *testing.T) {
	testFile := filepath.Join("testdata", "test_only_headers.xlsx")

	records, err := ParseSwiftXLSX(testFile)
	assert.NoError(t, err, "Parsing should succeed with only headers")
	assert.Empty(t, records, "No data rows expected when only headers exist")
}

// Tests if rows with insufficient columns are skipped.
func TestParseSwiftXLSX_MissingColumns(t *testing.T) {
	testFile := filepath.Join("testdata", "test_missing_cols.xlsx")

	records, err := ParseSwiftXLSX(testFile)
	assert.NoError(t, err, "Parsing should not fail on missing columns")
	assert.Len(t, records, 2, "Expected 2 valid records after skipping incomplete rows")
}

// Tests if ISO2 and countryName are converted to uppercase.
func TestParseSwiftXLSX_Uppercase(t *testing.T) {
	testFile := filepath.Join("testdata", "test_uppercase.xlsx")

	records, err := ParseSwiftXLSX(testFile)
	assert.NoError(t, err, "Parsing should succeed")
	assert.NotEmpty(t, records, "Expected at least one record")

	first := records[0]
	assert.Equal(t, "PL", first.CountryISO2, "ISO2 should be uppercase 'PL'")
	assert.Equal(t, "POLAND", first.CountryName, "Country name should be uppercase 'POLAND'")
}

// Tests if invalid files produce an error.
func TestParseSwiftXLSX_InvalidFile(t *testing.T) {
	testFile := filepath.Join("testdata", "not_really_excel.xlsx")

	_, err := ParseSwiftXLSX(testFile)
	assert.Error(t, err, "Should return an error for a non-Excel file")
}

// Tests if special characters (e.g., diacritics) are preserved in the address.
func TestParseSwiftXLSX_SpecialCharacters(t *testing.T) {
	testFile := filepath.Join("testdata", "test_special_chars.xlsx")

	records, err := ParseSwiftXLSX(testFile)
	assert.NoError(t, err, "Parsing should succeed")
	assert.NotEmpty(t, records, "Expected at least one record")

	first := records[0]
	assert.Contains(t, first.Address, "ZAŻÓŁĆ GĘŚLĄ", "Special characters should remain intact")
}

// Tests if mixed-case ISO2/countryName are turned uppercase.
func TestParseSwiftXLSX_MixedCase(t *testing.T) {
	testFile := filepath.Join("testdata", "test_mixed_case.xlsx")

	records, err := ParseSwiftXLSX(testFile)
	assert.NoError(t, err, "Parsing should succeed")
	assert.NotEmpty(t, records, "Expected at least one record")

	first := records[0]
	assert.Equal(t, "PL", first.CountryISO2, "ISO2 should be forced to uppercase")
	assert.Equal(t, "POLAND", first.CountryName, "Country name should be forced to uppercase")
}
