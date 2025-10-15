package handlers

import (
	"fmt"
	"net/http"
)

func ExecsHandler(w http.ResponseWriter, r *http.Request) {
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
