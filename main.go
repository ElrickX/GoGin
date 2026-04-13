package main

import (
	"log"
	"net/http"

	"gin_demo/internal/handler"
	"gin_demo/middleware"

	"github.com/golang-jwt/jwt/v5"

	_ "gin_demo/docs" // swag 生成的文件

	httpSwagger "github.com/swaggo/http-swagger"
)

// ===== Swagger 基本信息 =====

// @title Demo API
// @version 1.0
// @description Simple SaaS API with JWT Auth
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// ===== MAIN =====

func main() {

	http.HandleFunc("/login", middleware.LoggingMiddleware(handler.LoginHandler))
	http.HandleFunc("/report", middleware.Chain(handler.ReportHandler,
		middleware.LoggingMiddleware,
		middleware.AuthMiddleware,
	))

	// Swagger UI
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	port := ":8080"

	log.Println("🚀 Server started on http://localhost:8080")
	log.Println("To view API Docs, refer to http://localhost:8080/swagger")

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
