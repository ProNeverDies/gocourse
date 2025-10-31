package main

import (
	"crypto/tls"
	"fmt"
	mw "gocourse/internal/api/middlewares"
	"gocourse/internal/api/routes"
	"gocourse/internal/repository/sqlconnect"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type user struct {
	Name string `json:"name"`
	City string `json:"city"`
	Age  int    `json:"age"`
}

func main() {
	err := godotenv.Load("../../cmd/.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	_, err = sqlconnect.ConnectDb()

	if err != nil {
		log.Fatal("Database connection failed:", err)
		return
	}

	port := os.Getenv("API_PORT")

	cert := "cert.pem"

	key := "key.pem"

	// http.HandleFunc("/", rootHandler)

	// http.HandleFunc("/teachers", teacherHandler)

	// http.HandleFunc("/students", studentHandler)

	// http.HandleFunc("/execs", execsHandler)

	tlsconfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	// rl := mw.NewRateLimiter(5, 1*time.Minute)

	// hppOptions := mw.HPPOptions{
	// 	CheckQuery: 					true,
	// 	CheckBody: 						true,
	// 	CheckBodyOnlyForContentType: 	"application/x-www-form-urlencoded",
	// 	Whitelist: 						[]string{"sortBy", "sortOrder", "name", "age", "class"},
	// }
	// secureMux := mw.Cors(rl.Middleware(mw.ResponseTimeMiddleware(mw.SecurityHandlers(mw.Compression(mw.Hpp(hppOptions)(mux))))))
	// secureMux := applyMiddleware(mux, mw.Hpp(hppOptions), mw.Compression, mw.SecurityHandlers, mw.ResponseTimeMiddleware, rl.Middleware, mw.Cors)

	// secureMux := mw.SecurityHandlers(mux)
	// FIX: Chained all necessary middlewares. The request flows from outside in: Cors -> Compression -> SecurityHandlers -> Mux.

	router := routes.Router()
	secureMux := mw.Cors(mw.Compression(mw.SecurityHandlers(router)))
	// Efficency and logical ordereing of handling request is important while chaining middlewares
	server := http.Server{
		Addr:    port,
		Handler: secureMux,
		// Handler: mw.Cors(rl.Middleware(mw.ResponseTimeMiddleware(mw.SecurityHandlers(mw.Compression(mw.Hpp(hppOptions)(mux)))))),
		// Handler: rl.Middleware(mw.Compression(mw.ResponseTimeMiddleware(mw.Cors(mw.SecurityHandlers(mux))))),
		// Handler: mw.Cors(mux),
		// Handler: middlewares.SecurityHandlers(mux),1
		// Handler:mux,
		TLSConfig: tlsconfig,
	}

	// server := http.Server{
	// 	Addr: 		port,
	// 	Handler: 	nil,
	// 	TLSConfig: 	tlsconfig,
	// }

	fmt.Printf("Server is running on port %v\n", port)
	// err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatal("Error Starting the server", err)
	}
}
