package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Used for serving localhost:4200
var CorsDefault gin.HandlerFunc = cors.New(cors.Config{
	AllowOrigins:     []string{"http://localhost:4200"},
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowHeaders:     []string{"Content-Type", "Authorization"},
	AllowCredentials: true,
});