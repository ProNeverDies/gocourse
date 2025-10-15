package routes

import (
	"gocourse/internal/api/handlers"
	"net/http"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux() //helps us to use multiple api routes where each route can have their handler func

	mux.HandleFunc("/", handlers.RootHandler)

	mux.HandleFunc("/teachers/", handlers.TeacherHandler) // Using "/teachers/" to catch all sub-paths

	mux.HandleFunc("/students", handlers.StudentHandler)

	mux.HandleFunc("/execs", handlers.ExecsHandler)

	return mux
}
