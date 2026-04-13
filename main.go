package main

import (
	"gin_demo/internal/handler"
	"gin_demo/middleware"
	"log"
	"net/http"
	"os"

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

	// ✅ FIX HERE
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("🚀 Server started on port:", port)

	err := http.ListenAndServe(":"+port, enableCORS(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // allow all (dev)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
