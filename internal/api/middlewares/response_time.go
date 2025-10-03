package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Recieved Request in Response Time")
		start := time.Now()

		// Create a custom response writer to capture the status code

		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		duration := time.Since(start)

		wrappedWriter.Header().Set("X-Response-Time", duration.String())

		next.ServeHTTP(wrappedWriter, r)

		duration = time.Since(start)

		fmt.Printf("Method: %s, URL: %s, Status: %d, Duration: %v\n", r.Method, r.URL, wrappedWriter.status, duration.String())
		fmt.Println("Sent Response from Response Time Middleware")

	})
}

//Resopse Writer

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}
