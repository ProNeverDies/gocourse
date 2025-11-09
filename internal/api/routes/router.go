package routes

import (
	"gocourse/internal/api/handlers"
	"net/http"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux() //helps us to use multiple api routes where each route can have their handler func

	mux.HandleFunc("/", handlers.RootHandler)

	// mux.HandleFunc("/teachers/", handlers.TeacherHandler) // Using "/teachers/" to catch all sub-paths

	mux.HandleFunc("GET /teachers/", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers/", handlers.PostTeacherHandler)
	mux.HandleFunc("PATCH /teachers/", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teachers/", handlers.DeleteTeacherHandler)

	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)
	mux.HandleFunc("GET /teachers/{id}", handlers.GetOneTeacherHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchOneTeacherHandler)

	mux.HandleFunc("/students", handlers.StudentHandler)

	mux.HandleFunc("/execs", handlers.ExecsHandler)

	return mux
}
