package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gocourse/internal/models"
	"gocourse/internal/repository/sqlconnect"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

// var (
// 	teachers = make(map[int]models.Teacher) // Created a map for in memory storage ,as it searches and fetches data faster
// 	mutex    = &sync.Mutex{}                // This is our database
// 	netID    = 1
// )

// func init() {
// 	teachers[netID] = models.Teacher{
// 		ID:        netID,
// 		FirstName: "Akash",
// 		LastName:  "Kumar",
// 		Class:     "10th",
// 		Subject:   "Maths",
// 	}
// 	netID++
// 	teachers[netID] = models.Teacher{
// 		ID:        netID,
// 		FirstName: "Raj",
// 		LastName:  "Sharma",
// 		Class:     "9th",
// 		Subject:   "Science",
// 	}
// 	netID++
// }

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
		updateTeacherHandler(w, r)
		// w.Write([]byte("Hello Put method on teachers route"))
		// fmt.Println("Hello Put method on teachers route")
	case http.MethodPatch:
		patchTeacherHandler(w, r)
		// w.Write([]byte("Hello Patch method on teachers route"))
		// fmt.Println("Hello Patch method on teachers route")
	case http.MethodDelete:
		deleteTeacherHandler(w, r)
		// w.Write([]byte("Hello Delete method on teachers route"))
		// fmt.Println("Hello Delete method on teachers route")
	}

	// if r.Method == http.MethodGet {
	// 	w.Write([]byte("Hello GET method on teachers route"))
	// 	fmt.Println("Hello GET method on teachers route")
	// 	return
	// }

	// w.Write([]byte("Hello Teachers Route"))

	// fmt.Printf(r.Method) //http method which is sent to the route
}

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{"id": true, "first_name": true, "last_name": true, "email": true, "class": true, "subject": true}
	return validFields[field]
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
		// firstName := r.URL.Query().Get("first_name")
		// lastName := r.URL.Query().Get("last_name")

		// if firstName != "" {
		// 	query += " AND first_name = ?"
		// 	args = append(args, firstName)
		// }
		// if lastName != "" {
		// 	query += " AND last_name = ?"
		// 	args = append(args, lastName)
		// }

		query := "SELECT id,first_name,last_name,email,class,subject from teachers WHERE 1=1"
		var args []interface{}

		query, args = addFilters(r, query, args)

		// r.URL.Query().Get("sortby") 	will only get the first value if multiple are provided
		query = addSorting(r, query)

		rows, err := db.Query(query, args...)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Database Query Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		teacherList := make([]models.Teacher, 0)

		for rows.Next() {
			var teacher models.Teacher
			err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "Database Scan Error", http.StatusInternalServerError)
				return
			}
			teacherList = append(teacherList, teacher)
		}

		// FIX: Corrected filtering logic.
		// for _, v := range teachers {
		// 	// If no filters are provided, include everyone.
		// 	if firstName == "" && lastName == "" {
		// 		 teacherList = append(teacherList, v)
		// 		 continue
		// 	}
		// 	// If filters are provided, match them.
		// 	if (firstName != "" && v.FirstName == firstName) || (lastName != "" && v.LastName == lastName) {
		// 		 teacherList = append(teacherList, v)
		// 	}
		// }
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

func addSorting(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortby"] // This gets all values for the key "sortby" in slice form

	// FIX: Reworked sorting logic to handle params with and without colons
	if len(sortParams) > 0 {
		var sortClauses []string
		for _, param := range sortParams {
			parts := strings.Split(param, ":")
			var field, order string

			if len(parts) == 1 {
				// Case: sortby=id
				field = parts[0]
				order = "asc" // Default to ascending order
			} else if len(parts) == 2 {
				// Case: sortby=id:desc
				field = parts[0]
				order = parts[1]
			} else {
				// Invalid format, skip
				continue
			}

			// FIX: Inverted validation logic. We should skip if EITHER is invalid.
			if !isValidSortField(field) || !isValidSortOrder(order) {
				continue
			}

			sortClauses = append(sortClauses, " "+field+" "+order)
		}

		// Only add the ORDER BY clause if we have at least one valid sort clause
		if len(sortClauses) > 0 {
			query += " ORDER BY"
			query += strings.Join(sortClauses, ", ")
		}
	}
	return query
}
func addFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += fmt.Sprintf(" AND %s = ?", dbField)
			args = append(args, value)
		}
	}
	return query, args
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

// PUT /teachers/{id}
func updateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)

	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingTeacher models.Teacher

	err = db.QueryRow("Select id,first_name,last_name,email,class,subject from teachers where id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return
	}

	updatedTeacher.ID = existingTeacher.ID

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, id)
	if err != nil {
		http.Error(w, "Error Updating Teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)
}

//Patch /teachers/{id}  Not updates all the fields only the fields that we recieve

//Mtlb hum empty fields bhi pass kar sakte h during update

func patchTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}

	err = json.NewDecoder(r.Body).Decode(&updates)

	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingTeacher models.Teacher

	err = db.QueryRow("Select id,first_name,last_name,email,class,subject from teachers where id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return
	}

	//APPLY UPDATES

	// for key, value := range updates {
	// 	switch key {
	// 	case "first_name":
	// 		existingTeacher.FirstName = value.(string)
	// 	case "last_name":
	// 		existingTeacher.LastName = value.(string)
	// 	case "email":
	// 		existingTeacher.Email = value.(string)
	// 	case "class":
	// 		existingTeacher.Class = value.(string)
	// 	case "subject":
	// 		existingTeacher.Subject = value.(string)
	// 	}
	// }

	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	fmt.Println("TeacherVal field 0", teacherVal.Type().Field(0))
	fmt.Println("TeacherVal field 1", teacherVal.Type().Field(1))

	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			fmt.Println("K from the refelct mechanism", k)
			field := teacherVal.Type().Field(i)
			fmt.Println(field.Tag.Get("json"))

			if field.Tag.Get("json") == k+",omitempty" {
				if teacherVal.Field(i).CanSet() {
					teacherVal.Field(i).Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}

		}
	}

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, id)
	if err != nil {
		http.Error(w, "Error Updating Teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingTeacher)

}

func deleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)

	if err != nil {
		http.Error(w, "Error Deleting Teacher", http.StatusInternalServerError)
		return
	}
	// fmt.Println(result.RowsAffected())
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error retrieving Delete result", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	// w.WriteHeader(http.StatusNoContent)

	w.Header().Set("Content-Type", "application/json")

	response := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "success",
		Message: fmt.Sprintf("Teacher with ID %d deleted successfully", id),
	}

	json.NewEncoder(w).Encode(response)

}
