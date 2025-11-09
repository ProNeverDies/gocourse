package routes

import (
	"gocourse/internal/api/handlers"
	"net/http"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux() //helps us to use multiple api routes where each route can have their handler func

	mux.HandleFunc("/", handlers.RootHandler)

	mux.HandleFunc("/teachers/", handlers.TeacherHandler) // Using "/teachers/" to catch all sub-paths

	mux.HandleFunc("GET /teachers/", handlers.TeacherHandler)
	mux.HandleFunc("GET /teachers/{id}", handlers.TeacherHandler)
	mux.HandleFunc("POST /teachers/", handlers.TeacherHandler)
	mux.HandleFunc("PATCH /teachers/", handlers.TeacherHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.TeacherHandler)
	mux.HandleFunc("PUT /teachers/", handlers.TeacherHandler)
	mux.HandleFunc("DELETE /teachers/", handlers.TeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.TeacherHandler)

	mux.HandleFunc("/students", handlers.StudentHandler)

	mux.HandleFunc("/execs", handlers.ExecsHandler)

	return mux
}
