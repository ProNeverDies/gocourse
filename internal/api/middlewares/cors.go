package middlewares

import (
	"fmt"
	"net/http"
)

// Alowed Origins
var allowedOrigins = []string{
	"https://localhost:3000", // For HTTPS
	"http://localhost:3000",  // For HTTP
	"https://example.com",
}

func Cors(next http.Handler) http.Handler {
	fmt.Println("Cors Middleware... (initializing)")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")

		if isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			// Only block if origin is present and not in the list
			if origin != "" {
				http.Error(w, "Not allowed by CORS", http.StatusForbidden)
				return
			}
		}

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS") // Add OPTIONS
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// FIX 2: Handle OPTIONS preflight requests correctly
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
		// fmt.Println("CORS Middleware ends") // This is more useful here
	})
}

func isOriginAllowed(origin string) bool {
	// FIX 1: Allow requests with no origin (like Postman)
	if origin == "" {
		return true
	}

	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}
