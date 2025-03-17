package middleware

import (
	"github.com/ayayaakasvin/auth/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// var CorsDefaultConfig cors.Config = cors.Config{
// 	AllowOrigins:     []string{"http://localhost:4200"},
// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 	AllowHeaders:     []string{"Content-Type", "Authorization"},
// 	AllowCredentials: true,
// }


func CorsWithConfig (addresses config.ServiceAddresses) gin.HandlerFunc {
	var CorsDefaultConfig cors.Config = cors.Config{
		AllowOrigins: 		addresses.Addresses,
		AllowMethods: 		[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:		[]string{"Content-Type", "Authorization"},
		AllowCredentials: 	true,
	}

	return cors.New(CorsDefaultConfig)
}