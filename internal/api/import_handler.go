package api

import (
	"net/http"
	"path/filepath"

	"github.com/Mekambee/Swift-Codes-Api/internal/services"
	"github.com/gin-gonic/gin"
)

func ImportSwiftCodesHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file found"})
		return
	}

	tempPath := filepath.Join("/tmp", file.Filename)
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot save file"})
		return
	}

	swiftData, err := services.ParseSwiftXLSX(tempPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing XLSX"})
		return
	}

	err = services.SaveSwiftCodes(swiftData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Import successful"})
}
