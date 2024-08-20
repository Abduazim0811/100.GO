package api

import (
	"100.GO/internal/pkg/token"
	"100.GO/internal/storage"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() {
	handler := storage.Handler()

	r := gin.Default()

	r.POST("/register", handler.CreateUser)
	r.POST("/verify", handler.VerifyCode)
	r.POST("/login", handler.Login)

	r.POST("/origins", token.Protected(), handler.CreateOrigin)
	r.GET("/origins", token.Protected(), handler.GetOrigin)
	r.GET("/origins/:id", token.Protected(), handler.GetbyIdOrigin)
	r.PUT("/origins/:id", token.Protected(), handler.UpdateOrigin)
	r.DELETE("/origins/:id", token.Protected(), handler.DeleteOrigin)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8888")
}
