package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		fmt.Printf("📥 %s %s | IP: %s\n",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)

		next(w, r)

		fmt.Printf("✅ %s %s | %v\n",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	}
}
