package middlewares

import (
	"fmt"
	"net/http"
)

func SecurityHandlers(next http.Handler) http.Handler {
	fmt.Println("Security Middleware returned...")
	fmt.Println("Security Middleware returned...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("X-DNS-Prefetch-Control", "off")                                            //Prevent DNS prefetching in background
		w.Header().Set("X-Frame-Options", "DENY")                                                  //Defense Agaisnt Clickjacking (embedd iframes)
		w.Header().Set("X-XSS-Protection", "1; mode=block")                                        //Cross Site Scripting Protection
		w.Header().Set("X-Content-Type-Options", "nosniff")                                        //Mime Sniffing Protection (Multipurpose Internet Mail Exchange )
		w.Header().Set("Strict Transport-Security", "max-age=31536000; includeSubDomains;preload") //Interact Over HTTPS only
		w.Header().Set("Content-Security-Policy", "default-src 'self'")                            //Which resources are loaded on the page (same origin)
		w.Header().Set("Referrer-Policy", "no-referrer")                                           //Some or no referrer Info is passed
		w.Header().Set("X-Powered-By", "rubyOnRails")                                              //Custom Header created to avoid revealing tech stack
		w.Header().Set("Server", "")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Permissions-Policy", "geolocation=(self), microphone=()")

		next.ServeHTTP(w, r)
	})

}

// Basic Middleware template
// func securityHandlers(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// next.ServeHTTP(w,r)
// 	})

// }
