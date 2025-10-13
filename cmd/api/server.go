package main

import (
	"crypto/tls"
	"fmt"
	mw "gocourse/internal/api/middlewares"
	"log"
	"net/http"
)

type user struct {
	Name string `json:"name"`
	City string `json:"city"`
	Age  int    `json:"age"`
}

// Handler Function
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello Root Route")
	w.Write([]byte("Hello Root Route"))
}

func teacherHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello Teachers Route")

	// Since there will be multiple methods and nested if condition it is better to use switch case
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET method on teachers route"))
		fmt.Println("Hello GET method on teachers route")
	case http.MethodPost:

		// Parse the data imp for x-www-form-urlencoded
		// err := r.ParseForm()
		// if err != nil {
		// 	http.Error(w, "Error parsing form", http.StatusBadRequest)
		// 	return
		// }

		// fmt.Println("Form", r.Form)

		// response := make(map[string]interface{})

		// for key, value := range r.Form {
		// 	response[key] = value[0]
		// }

		// fmt.Println("Processed Response", response)

		// //Raw Body

		// body, err := io.ReadAll(r.Body)

		// if err != nil {
		// 	return
		// }
		// defer r.Body.Close()

		// fmt.Println("Raw Body", string(body))

		// // Unmarshall in case of json data
		// var userInstance user
		// err = json.Unmarshal(body, &userInstance)
		// if err != nil {
		// 	return
		// }

		// fmt.Println("User Instance", userInstance)
		// fmt.Println("Name", userInstance.Name)

		w.Write([]byte("Hello Post method on teachers route"))
		fmt.Println("Hello Post method on teachers route")
	case http.MethodPut:
		w.Write([]byte("Hello Put method on teachers route"))
		fmt.Println("Hello Put method on teachers route")
	case http.MethodPatch:
		w.Write([]byte("Hello Patch method on teachers route"))
		fmt.Println("Hello Patch method on teachers route")
	case http.MethodDelete:
		w.Write([]byte("Hello Delete method on teachers route"))
		fmt.Println("Hello Delete method on teachers route")
	}

	// if r.Method == http.MethodGet {
	// 	w.Write([]byte("Hello GET method on teachers route"))
	// 	fmt.Println("Hello GET method on teachers route")
	// 	return
	// }

	// w.Write([]byte("Hello Teachers Route"))

	// fmt.Printf(r.Method) //http method which is sent to the route
}

func studentHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello students Route")
	switch r.Method {
	case http.MethodGet:
		// fmt.Println(r.URL.Path)
		// path := strings.TrimPrefix(r.URL.Path, "/students/")
		// userID := strings.TrimSuffix(path, "/")

		// fmt.Println("User ID", userID)

		// fmt.Println("User Query", r.URL.Query())
		// queryParams := r.URL.Query()

		// sortby := queryParams.Get("sortby")
		// key := queryParams.Get("key")
		// sortorder := queryParams.Get("sortorder")

		// if sortorder == "" {
		// 	sortorder = "desc"
		// }
		// fmt.Printf("Sortby: %s, Key: %s, Sortorder: %s\n", sortby, key, sortorder)
		w.Write([]byte("Hello GET method on students route"))
		fmt.Println("Hello GET method on students route")
	case http.MethodPost:
		w.Write([]byte("Hello Post method on students route"))
		fmt.Println("Hello Post method on students route")
	case http.MethodPut:
		w.Write([]byte("Hello Put method on students route"))
		fmt.Println("Hello Put method on students route")
	case http.MethodPatch:
		w.Write([]byte("Hello Patch method on students route"))
		fmt.Println("Hello Patch method on students route")
	case http.MethodDelete:
		w.Write([]byte("Hello Delete method on students route"))
		fmt.Println("Hello Delete method on students route")
	}
	// w.Write([]byte("Hello students Route"))
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello execs Route")
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello GET method on execs route"))
		fmt.Println("Hello GET method on execs route")
	case http.MethodPost:
		w.Write([]byte("Hello Post method on execs route"))
		fmt.Println("Hello Post method on execs route")
	case http.MethodPut:
		w.Write([]byte("Hello Put method on execs route"))
		fmt.Println("Hello Put method on execs route")
	case http.MethodPatch:
		w.Write([]byte("Hello Patch method on execs route"))
		fmt.Println("Hello Patch method on execs route")
	case http.MethodDelete:
		w.Write([]byte("Hello Delete method on execs route"))
		fmt.Println("Hello Delete method on execs route")
	}
	// w.Write([]byte("Hello execs Route"))
}

func main() {
	port := ":3000"

	cert := "cert.pem"

	key := "key.pem"

	mux := http.NewServeMux() //helps us to use multiple api routes where each route can have their handler func

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/teachers", teacherHandler)

	mux.HandleFunc("/students", studentHandler)

	mux.HandleFunc("/execs", execsHandler)

	// http.HandleFunc("/", rootHandler)

	// http.HandleFunc("/teachers", teacherHandler)

	// http.HandleFunc("/students", studentHandler)

	// http.HandleFunc("/execs", execsHandler)

	tlsconfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := http.Server{
		Addr:    port,
		Handler: mw.Compression(mw.ResponseTimeMiddleware(mw.Cors(mw.SecurityHandlers(mux)))),
		// Handler: mw.Cors(mux),
		// Handler: middlewares.SecurityHandlers(mux),1
		// Handler:mux,
		TLSConfig: tlsconfig,
	}

	// server := http.Server{
	// 	Addr:      port,
	// 	Handler:   nil,
	// 	TLSConfig: tlsconfig,
	// }

	fmt.Printf("Server is running on port %v\n", port)
	// err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatal("Error Starting the server", err)
	}
}
