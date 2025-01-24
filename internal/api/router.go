package api

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		swiftCodes := v1.Group("/swift-codes")
		{
			swiftCodes.GET("/:swiftCode", GetSwiftCodeHandler)
			swiftCodes.GET("/country/:countryISO2", GetByCountryHandler)
			swiftCodes.POST("/", CreateSwiftCodeHandler)
			swiftCodes.DELETE("/:swiftCode", DeleteSwiftCodeHandler)
			swiftCodes.POST("/import", ImportSwiftCodesHandler)
		}
	}

	return router
}
