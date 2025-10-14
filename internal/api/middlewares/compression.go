package middlewares

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

func Compression(next http.Handler) http.Handler {
	fmt.Println("Compression Middleware...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The response write method by deafult will not compress and send the header that why we have to use gzip Response Writer
		//Check if the client accepts the gunzip encoding
		fmt.Println("Compression Middleware being returned...")
		if !strings.Contains(r.Header.Get("Accept-encoding"), "gzip") {
			next.ServeHTTP(w, r)
		}

		// Set the response header

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		w = &gzipResponseWriter{ResponseWriter: w, Writer: gz}

		next.ServeHTTP(w, r)
		fmt.Println("Sent Response from Compression Middleware")
	})
}

// gzipResponseWriter wraps http.ResponseWriter to write gzipped responses
type gzipResponseWriter struct { // We try to create an interface that implements response writer
	// When we hover over it we can see the functions that it can implement
	http.ResponseWriter
	Writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}
