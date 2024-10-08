package api

import (
	"log"
	"net/http"

	"100.GO/internal/embed"
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
	r.GET("/cors", handler.EnableCORS(), handler.CorsMessage)

	go func() {
		log.Println("Server started on port 8081")
		t := embed.GetTemplate()
		http.HandleFunc("/cors", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(t)
		})
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8888")
}
