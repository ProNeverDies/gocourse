package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gocourse/internal/models"
	"gocourse/internal/repository/sqlconnect"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	teachers = make(map[int]models.Teacher) // Created a map for in memory storage ,as it searches and fetches data faster
	mutex    = &sync.Mutex{}                // This is our database
	netID    = 1
)

func init() {
	teachers[netID] = models.Teacher{
		ID:        netID,
		FirstName: "Akash",
		LastName:  "Kumar",
		Class:     "10th",
		Subject:   "Maths",
	}
	netID++
	teachers[netID] = models.Teacher{
		ID:        netID,
		FirstName: "Raj",
		LastName:  "Sharma",
		Class:     "9th",
		Subject:   "Science",
	}
	netID++
}

func TeacherHandler(w http.ResponseWriter, r *http.Request) {
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

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")
	// fmt.Println("IDStr", idStr)

	// This block handles GET /teachers and GET /teachers?first_name=...
	if idStr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")
		teacherList := make([]models.Teacher, 0, len(teachers))

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
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
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
	// teacher, exists := teachers[id]
	// if !exists {
	// 	http.Error(w, "Teacher not found", http.StatusNotFound)
	// 	return
	// }

	var teacher models.Teacher

	err = db.QueryRow("SELECT id,first_name,last_name,email,class,subject from teachers where id = ?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database Query Error", http.StatusNotFound)
		return
	}
	// FIX: Set Content-Type header before sending response.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)

}

func postTeacherHandler(w http.ResponseWriter, r *http.Request) {
	// mutex.Lock()
	// defer mutex.Unlock()

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var newTeachers []models.Teacher

	// FIX: Must pass a pointer to the slice for the decoder to populate it.
	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO teachers(first_name, last_name, email,class, subject) VALUES(?,?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Database statement preparation error", http.StatusInternalServerError)
		return
	}

	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))

	for i, newteacher := range newTeachers {
		res, err := stmt.Exec(newteacher.FirstName, newteacher.LastName, newteacher.Email, newteacher.Class, newteacher.Subject)
		if err != nil {
			http.Error(w, "Database insertion error", http.StatusInternalServerError)
			return
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Database retrieval error", http.StatusInternalServerError)
			return
		}
		newteacher.ID = int(lastID)
		addedTeachers[i] = newteacher
		// newteacher.ID = netID
		// teachers[netID] = newteacher
		// addedTeachers[i] = newteacher
		// netID++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(response)
}
