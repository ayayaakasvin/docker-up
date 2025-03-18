package middleware

import "github.com/gin-gonic/gin"

func LoadRouterWithMiddleware(router *gin.Engine, middlewares... gin.HandlerFunc) {
	router.Use(middlewares...)
}