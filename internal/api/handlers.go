package api

import (
	"net/http"
	"strings"

	"github.com/Mekambee/Swift-Codes-Api/internal/models"
	"github.com/Mekambee/Swift-Codes-Api/internal/services"
	"github.com/gin-gonic/gin"
)

// Endpoint 1: GET /v1/swift-codes/{swiftCode}
func GetSwiftCodeHandler(c *gin.Context) {
	swiftCode := c.Param("swiftCode")

	sc, err := services.GetSwiftCode(swiftCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Swift code not found"})
		return
	}

	if sc.IsHeadquarter {
		branches, err := services.GetBranchesByHQ(swiftCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve branches"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"address":       sc.Address,
			"bankName":      sc.BankName,
			"countryISO2":   sc.CountryISO2,
			"countryName":   sc.CountryName,
			"isHeadquarter": sc.IsHeadquarter,
			"swiftCode":     sc.SwiftCode,
			"branches":      branches,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address":       sc.Address,
		"bankName":      sc.BankName,
		"countryISO2":   sc.CountryISO2,
		"countryName":   sc.CountryName,
		"isHeadquarter": sc.IsHeadquarter,
		"swiftCode":     sc.SwiftCode,
	})
}

// Endpoint 2: GET /v1/swift-codes/country/{countryISO2}
func GetByCountryHandler(c *gin.Context) {
	iso2 := c.Param("countryISO2")
	iso2 = strings.ToUpper(iso2)

	data, err := services.GetSwiftByCountryISO2(iso2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve data"})
		return
	}
	if len(data) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No records found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"countryISO2": iso2,
		"countryName": data[0].CountryName,
		"swiftCodes":  data,
	})
}

// Endpoint 3: POST /v1/swift-codes/
/*
{
"address": string,
"bankName": string,
"countryISO2": string,
"countryName": string,
“isHeadquarter”: bool,
"swiftCode": string,
}
*/
type CreateSwiftCodeRequest struct {
	Address       string `json:"address" binding:"required"`
	BankName      string `json:"bankName" binding:"required"`
	CountryISO2   string `json:"countryISO2" binding:"required"`
	CountryName   string `json:"countryName" binding:"required"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode" binding:"required"`
}

func CreateSwiftCodeHandler(c *gin.Context) {
	var req CreateSwiftCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.CountryISO2 = strings.ToUpper(req.CountryISO2)
	req.CountryName = strings.ToUpper(req.CountryName)

	sc := models.SwiftCodeData{
		SwiftCode:     req.SwiftCode,
		BankName:      req.BankName,
		Address:       req.Address,
		CountryISO2:   req.CountryISO2,
		CountryName:   req.CountryName,
		IsHeadquarter: req.IsHeadquarter,
	}

	err := services.SaveSwiftCodes([]models.SwiftCodeData{sc})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Swift code created"})
}

// Endpoint 4: DELETE /v1/swift-codes/{swiftCode}
/*
{
    "bankName": "BANK NAME",
    "countryISO2": "AA"
}
*/
type DeleteSwiftCodeRequest struct {
	BankName    string `json:"bankName"`
	CountryISO2 string `json:"countryISO2"`
}

func DeleteSwiftCodeHandler(c *gin.Context) {
	swiftCode := c.Param("swiftCode")

	var req DeleteSwiftCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.CountryISO2 = strings.ToUpper(req.CountryISO2)

	err := services.DeleteSwiftCode(swiftCode, req.BankName, req.CountryISO2)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Swift code deleted"})
}
