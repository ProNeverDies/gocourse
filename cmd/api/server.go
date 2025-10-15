package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	mw "gocourse/internal/api/middlewares"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type user struct {
	Name string `json:"name"`
	City string `json:"city"`
	Age  int    `json:"age"`
}
type Teacher struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Class     string `json:"class,omitempty"`
	Subject   string `json:"subject,omitempty"`
}

var (
	teachers = make(map[int]Teacher) // Created a map for in memory storage ,as it searches and fetches data faster
	mutex    = &sync.Mutex{}         // This is our database
	netID    = 1
)

func init() {
	teachers[netID] = Teacher{
		ID:        netID,
		FirstName: "Akash",
		LastName:  "Kumar",
		Class:     "10th",
		Subject:   "Maths",
	}
	netID++
	teachers[netID] = Teacher{
		ID:        netID,
		FirstName: "Raj",
		LastName:  "Sharma",
		Class:     "9th",
		Subject:   "Science",
	}
	netID++
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")
	// fmt.Println("IDStr", idStr)

	// This block handles GET /teachers and GET /teachers?first_name=...
	if idStr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")
		teacherList := make([]Teacher, 0, len(teachers))

		// FIX: Corrected filtering logic.
		for _, v := range teachers {
			// If no filters are provided, include everyone.
			if firstName == "" && lastName == "" {
				teacherList = append(teacherList, v)
				continue
			}
			// If filters are provided, match them.
			if (firstName != "" && v.FirstName == firstName) || (lastName != "" && v.LastName == lastName) {
				teacherList = append(teacherList, v)
			}
		}
		response := struct {
			Status string    `json:"status"`
			Count  int       `json:"count"`
			Data   []Teacher `json:"data"`
		}{
			Status: "sucess",
			// FIX: Count should reflect the number of items in the filtered list.
			Count: len(teacherList),
			Data:  teacherList,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		// FIX: Added return to prevent code from falling through to the next block.
		return
	}

	// This block handles GET /teachers/{id}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}
	teacher, exists := teachers[id]
	if !exists {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	// FIX: Set Content-Type header before sending response.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)

}

func postTeacherHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []Teacher
	// FIX: Must pass a pointer to the slice for the decoder to populate it.
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	addedTeachers := make([]Teacher, len(newTeachers))

	for i, newteacher := range newTeachers {
		newteacher.ID = netID
		teachers[netID] = newteacher
		addedTeachers[i] = newteacher
		netID++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string    `json:"status"`
		Count  int       `json:"count"`
		Data   []Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(response)
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
		getTeachersHandler(w, r)
		// w.Write([]byte("Hello GET method on teachers route"))
		// fmt.Println("Hello GET method on teachers route")
	case http.MethodPost:
		postTeacherHandler(w, r)
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

		// w.Write([]byte("Hello Post method on teachers route"))
		// fmt.Println("Hello Post method on teachers route")
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
		// fmt.Println("Query Params:", r.URL.Query())
		// fmt.Println("Query Params:", r.URL.Query().Get("name"))

		// err := r.ParseForm()
		// if err != nil {
		// 	return
		// }
		// fmt.Println("Form Data:", r.Form)

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

	mux.HandleFunc("/teachers/", teacherHandler) // Using "/teachers/" to catch all sub-paths

	mux.HandleFunc("/students", studentHandler)

	mux.HandleFunc("/execs", execsHandler)

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
	secureMux := mw.Cors(mw.Compression(mw.SecurityHandlers(mux)))
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
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatal("Error Starting the server", err)
	}
}

type Middleware func(http.Handler) http.Handler

// func applyMiddleware(handler http.Handler, middlewares ...Middleware) http.Handler {
// 	for _, middleware := range middlewares {
// 		handler = middleware(handler)
// 	}
// 	return handler
// }
