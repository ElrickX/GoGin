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

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
