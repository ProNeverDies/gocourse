package handlers

import (
	"fmt"
	"net/http"
)

func StudentHandler(w http.ResponseWriter, r *http.Request) {
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
