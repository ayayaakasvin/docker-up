package middleware

import "github.com/gin-gonic/gin"

func MiddlewareAdd(router *gin.Engine, middlewares... gin.HandlerFunc) {
	router.Use(middlewares...)
}