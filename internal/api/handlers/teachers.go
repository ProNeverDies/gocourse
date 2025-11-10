package handlers

import (
	"encoding/json"
	"fmt"
	"gocourse/internal/models"
	"gocourse/internal/repository/sqlconnect"
	"log"
	"net/http"
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

// func TeacherHandler(w http.ResponseWriter, r *http.Request) {
// 	// fmt.Fprintf(w, "Hello Teachers Route")

// 	// Since there will be multiple methods and nested if condition it is better to use switch case
// 	switch r.Method {
// 	case http.MethodGet:
// 		getTeachersHandler(w, r)
// 		// w.Write([]byte("Hello GET method on teachers route"))
// 		// fmt.Println("Hello GET method on teachers route")
// 	case http.MethodPost:
// 		postTeacherHandler(w, r)
// 		// Parse the data imp for x-www-form-urlencoded
// 		// err := r.ParseForm()
// 		// if err != nil {
// 		// 	http.Error(w, "Error parsing form", http.StatusBadRequest)
// 		// 	return
// 		// }

// 		// fmt.Println("Form", r.Form)

// 		// response := make(map[string]interface{})

// 		// for key, value := range r.Form {
// 		// 	response[key] = value[0]
// 		// }

// 		// fmt.Println("Processed Response", response)

// 		// //Raw Body

// 		// body, err := io.ReadAll(r.Body)

// 		// if err != nil {
// 		// 	return
// 		// }
// 		// defer r.Body.Close()

// 		// fmt.Println("Raw Body", string(body))

// 		// // Unmarshall in case of json data
// 		// var userInstance user
// 		// err = json.Unmarshal(body, &userInstance)
// 		// if err != nil {
// 		// 	return
// 		// }

// 		// fmt.Println("User Instance", userInstance)
// 		// fmt.Println("Name", userInstance.Name)

// 		// w.Write([]byte("Hello Post method on teachers route"))
// 		// fmt.Println("Hello Post method on teachers route")
// 	case http.MethodPut:
// 		updateTeacherHandler(w, r)
// 		// w.Write([]byte("Hello Put method on teachers route"))
// 		// fmt.Println("Hello Put method on teachers route")
// 	case http.MethodPatch:
// 		patchTeacherHandler(w, r)
// 		// w.Write([]byte("Hello Patch method on teachers route"))
// 		// fmt.Println("Hello Patch method on teachers route")
// 	case http.MethodDelete:
// 		deleteTeacherHandler(w, r)
// 		// w.Write([]byte("Hello Delete method on teachers route"))
// 		// fmt.Println("Hello Delete method on teachers route")
// 	}

// 	// if r.Method == http.MethodGet {
// 	// 	w.Write([]byte("Hello GET method on teachers route"))
// 	// 	fmt.Println("Hello GET method on teachers route")
// 	// 	return
// 	// }

// 	// w.Write([]byte("Hello Teachers Route"))

// 	// fmt.Printf(r.Method) //http method which is sent to the route
// }

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {

	// path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	// idStr := strings.TrimSuffix(path, "/")
	// fmt.Println("IDStr", idStr)

	// This block handles GET /teachers and GET /teachers?first_name=...

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

	var teachers []models.Teacher

	teachers, err := sqlconnect.GetTeachersDbHandler(teachers, r)
	if err != nil {
		http.Error(w, "Error fetching teachers from database", http.StatusInternalServerError)
		return
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
		Count: len(teachers),
		Data:  teachers,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	// FIX: Added return to prevent code from falling through to the next block.

	// This block handles GET /teachers/{id}
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
	// 	return
	// }
	// teacher, exists := teachers[id]
	// if !exists {
	// 	http.Error(w, "Teacher not found", http.StatusNotFound)
	// 	return
	// }

}

func GetOneTeacherHandler(w http.ResponseWriter, r *http.Request) {

	idStr := r.PathValue("id")

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

	teacher, err := sqlconnect.GetTeacherByID(id)
	if err != nil {
		fmt.Println("Error fetching teacher by ID:", err)
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	// FIX: Set Content-Type header before sending response.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)

}

func PostTeacherHandler(w http.ResponseWriter, r *http.Request) {
	// mutex.Lock()
	// defer mutex.Unlock()

	var newTeachers []models.Teacher

	// FIX: Must pass a pointer to the slice for the decoder to populate it.
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	addedTeachers, err := sqlconnect.PostTeacherDBHandler(newTeachers)
	if err != nil {
		http.Error(w, "Error adding teachers to database", http.StatusInternalServerError)
		return
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
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
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

	updatedTeacherFromDB, err := sqlconnect.UpdateTeacher(id, updatedTeacher)
	if err != nil {
		http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacherFromDB)
}

//Patch /teachers/{id}  Not updates all the fields only the fields that we recieve

//Mtlb hum empty fields bhi pass kar sakte h during update

func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&updates)

	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = sqlconnect.PatchTeachers(updates)
	if err != nil {
		http.Error(w, "Error patching teachers", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func PatchOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	// idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := r.PathValue("id")
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

	updatedTeacher, err := sqlconnect.PatchOneTeacher(id, updates)
	if err != nil {
		http.Error(w, "Error patching teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)

}

func DeleteOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	// idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	err = sqlconnect.DeleteOneTeacher(id)
	if err != nil {
		http.Error(w, "Error deleting teacher", http.StatusInternalServerError)
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

func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var ids []int
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(ids) == 0 {
		http.Error(w, "No IDs provided", http.StatusBadRequest)
		return
	}

	deletedIds, err := sqlconnect.DeleteTeachers(ids)
	if err != nil {
		http.Error(w, "Error deleting teachers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status     string `json:"status"`
		DeletedIDs []int  `json:"deleted_ids"`
	}{
		Status:     "success",
		DeletedIDs: deletedIds,
	}
	json.NewEncoder(w).Encode(response)
}
