package sqlconnect

import (
	"database/sql"
	"fmt"
	"gocourse/internal/models"
	"log"
	"reflect"
	"strconv"

	"net/http"
	"strings"
)

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{"id": true, "first_name": true, "last_name": true, "email": true, "class": true, "subject": true}
	return validFields[field]
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

func GetTeachersDbHandler(teachers []models.Teacher, r *http.Request) ([]models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	query := "SELECT id,first_name,last_name,email,class,subject from teachers WHERE 1=1"
	var args []interface{}

	query, args = addFilters(r, query, args)

	// r.URL.Query().Get("sortby") 	will only get the first value if multiple are provided
	query = addSorting(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		// http.Error(w, "Database Query Error", http.StatusInternalServerError)
		return nil, err
	}
	defer rows.Close()
	// teacherList := make([]models.Teacher, 0)

	for rows.Next() {
		var teacher models.Teacher
		err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
		if err != nil {
			fmt.Println(err)
			// http.Error(w, "Database Scan Error", http.StatusInternalServerError)
			return nil, err
		}
		teachers = append(teachers, teacher)
	}
	return teachers, nil
}

func GetTeacherByID(id int) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	defer db.Close()

	var teacher models.Teacher

	err = db.QueryRow("SELECT id,first_name,last_name,email,class,subject from teachers where id = ?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		// http.Error(w, "Teacher not found", http.StatusNotFound)
		return models.Teacher{}, err
	} else if err != nil {
		// http.Error(w, "Database Query Error", http.StatusNotFound)
		return models.Teacher{}, err
	}
	return teacher, nil
}

func PostTeacherDBHandler(newTeachers []models.Teacher) ([]models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers(first_name, last_name, email,class, subject) VALUES(?,?, ?, ?, ?)")
	if err != nil {
		// http.Error(w, "Database statement preparation error", http.StatusInternalServerError)
		return nil, err
	}

	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))

	for i, newteacher := range newTeachers {
		res, err := stmt.Exec(newteacher.FirstName, newteacher.LastName, newteacher.Email, newteacher.Class, newteacher.Subject)
		if err != nil {
			// http.Error(w, "Database insertion error", http.StatusInternalServerError)
			return nil, err
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			// http.Error(w, "Database retrieval error", http.StatusInternalServerError)
			return nil, err
		}
		newteacher.ID = int(lastID)
		addedTeachers[i] = newteacher
		// newteacher.ID = netID
		// teachers[netID] = newteacher
		// addedTeachers[i] = newteacher
		// netID++
	}
	return addedTeachers, nil
}

func UpdateTeacher(id int, updatedTeacher models.Teacher) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	defer db.Close()

	var existingTeacher models.Teacher

	err = db.QueryRow("Select id,first_name,last_name,email,class,subject from teachers where id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		// http.Error(w, "Teacher not found", http.StatusNotFound)
		return models.Teacher{}, err
	}
	if err != nil {
		// http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return models.Teacher{}, err
	}

	updatedTeacher.ID = existingTeacher.ID

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, id)
	if err != nil {
		// http.Error(w, "Error Updating Teacher", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	return updatedTeacher, nil
}

func PatchTeachers(updates []map[string]interface{}) error {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return err
	}
	defer db.Close()

	tx, err := db.Begin() //Starts a new transacation
	if err != nil {
		// http.Error(w, "Database transaction error", http.StatusInternalServerError)
		return err
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			// http.Error(w, "Invalid or missing teacher ID", http.StatusBadRequest)
			tx.Rollback()
			return err
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println(err)
			// http.Error(w, "Error Converting ID to Int", http.StatusBadRequest)
			tx.Rollback()
			return err
		}

		var teacherFromDb models.Teacher

		err = db.QueryRow("SELECT id,first_name,last_name,email,class,subject from teachers where id = ?", id).Scan(&teacherFromDb.ID, &teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class, &teacherFromDb.Subject)

		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				// http.Error(w, fmt.Sprintf("Teacher with ID %d not found", id), http.StatusNotFound)
				return err
			}
			// http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
			return err
		}

		//Apply updates using reflection

		teacherVal := reflect.ValueOf(&teacherFromDb).Elem()
		teacherType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				continue // skip updating the id field
			}

			for i := 0; i < teacherVal.NumField(); i++ {
				field := teacherType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := teacherVal.Field(i)

					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)

						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							log.Printf("cannot convert %v to %v", val.Type(), fieldVal.Type())
							return err
						}
					}
					break
				}
			}
		}

		_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", teacherFromDb.FirstName, teacherFromDb.LastName, teacherFromDb.Email, teacherFromDb.Class, teacherFromDb.Subject, id)

		if err != nil {
			tx.Rollback()
			// http.Error(w, "Error Updating Teacher", http.StatusInternalServerError)
			return err
		}

	}

	err = tx.Commit()
	if err != nil {
		// http.Error(w, "Database commit error", http.StatusInternalServerError)
		return err
	}
	return nil
}

func PatchOneTeacher(id int, updates map[string]interface{}) (models.Teacher, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	defer db.Close()

	var existingTeacher models.Teacher

	err = db.QueryRow("Select id,first_name,last_name,email,class,subject from teachers where id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err == sql.ErrNoRows {
		// http.Error(w, "Teacher not found", http.StatusNotFound)
		return models.Teacher{}, err
	}
	if err != nil {
		// http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return models.Teacher{

			//APPLY UPDATES
		}, err
	}

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
		// http.Error(w, "Error Updating Teacher", http.StatusInternalServerError)
		return models.Teacher{}, err
	}
	return existingTeacher, nil
}

func DeleteOneTeacher(id int) error {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return err
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)

	if err != nil {
		// http.Error(w, "Error Deleting Teacher", http.StatusInternalServerError)
		return err
	}
	// fmt.Println(result.RowsAffected())
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// http.Error(w, "Error retrieving Delete result", http.StatusInternalServerError)
		return err
	}

	if rowsAffected == 0 {
		// http.Error(w, "Teacher not found", http.StatusNotFound)
		return err
	}
	return nil
}

func DeleteTeachers(ids []int) ([]int, error) {
	db, err := ConnectDb()
	if err != nil {
		// http.Error(w, "Database connection error", http.StatusInternalServerError)
		return nil, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		// http.Error(w, "Database transaction error", http.StatusInternalServerError)
		return nil, err
	}

	stmt, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		// http.Error(w, "Database statement preparation error", http.StatusInternalServerError)
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	deletedIds := []int{}

	for _, id := range ids {
		res, err := stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			// http.Error(w, "Error deleting teacher", http.StatusInternalServerError)
			return nil, err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			tx.Rollback()
			// http.Error(w, "Error retrieving delete result", http.StatusInternalServerError)
			return nil, err
		}

		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
		if rowsAffected < 1 {
			tx.Rollback()
			// http.Error(w, fmt.Sprintf("Teacher with ID %d not found", id), http.StatusNotFound)
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		// http.Error(w, "Database commit error", http.StatusInternalServerError)
		return nil, err
	}

	if len(deletedIds) == 0 {
		// http.Error(w, "IDs do not exist", http.StatusNotFound)
		return nil, err
	}
	return deletedIds, nil
}
