package api

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := SetupRouter()
	assert.NotNil(t, r, "SetupRouter should return a router")

	routesInfo := r.Routes()
	var foundGetHQ bool
	for _, route := range routesInfo {
		if route.Method == "GET" && route.Path == "/v1/swift-codes/:swiftCode" {
			foundGetHQ = true
			break
		}
	}
	assert.True(t, foundGetHQ, "Should have GET /v1/swift-codes/:swiftCode route")
}
