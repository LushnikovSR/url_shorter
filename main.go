package main

import (
	"fmt"
	"net/http"
)

func main() {

	storage := make(map[string]string, 1000)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
		//io.WriteString(w, "Hello from web-server!")
	})

	http.HandleFunc("/name", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		name := query.Get("q")
		fmt.Fprintf(w, "Hello %s!", name)
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		key := query.Get("k")
		value := query.Get("v")
		fmt.Fprintf(w, "key: %s - value: %s", key, value)
		if storage[key] == "" {
			storage[key] = value
		} else {
			fmt.Fprintf(w, "This key:%s exists, try a different key or delete the existing key!", key)
		}
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		key := query.Get("k")
		fmt.Fprintf(w, "key: %s \n", key)
		if value, exists := storage[key]; exists {
			fmt.Fprintf(w, "%s", value)
		} else {
			fmt.Fprintf(w, "This key:%s doesn`t exist", key)
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
