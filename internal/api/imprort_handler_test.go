package api_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Mekambee/Swift-Codes-Api/internal/api"
	"github.com/Mekambee/Swift-Codes-Api/internal/database"
)

func TestImportSwiftCodesHandler(t *testing.T) {
	err := database.ConnectAndMigrate("localhost", "myuser", "mysecretpassword", "swiftdb", "5432")
	assert.NoError(t, err)
	defer database.DB.Close()

	database.DB.Exec("TRUNCATE swift_codes RESTART IDENTITY")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/v1/swift-codes/import", api.ImportSwiftCodesHandler)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	testFile := filepath.Join("testdata", "import_endpoint_test.xlsx")
	file, err := os.Open(testFile)
	assert.NoError(t, err, "Should open file")
	defer file.Close()

	part, err := writer.CreateFormFile("file", "import_test.xlsx")
	assert.NoError(t, err)

	_, err = io.Copy(part, file)
	assert.NoError(t, err)

	writer.Close()

	req, _ := http.NewRequest("POST", "/v1/swift-codes/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code, "Expected 200 OK for import")

	var count int
	row := database.DB.QueryRow("SELECT COUNT(*) FROM swift_codes")
	err = row.Scan(&count)
	assert.NoError(t, err)
	assert.True(t, count > 0, "Should import at least one record")
}
