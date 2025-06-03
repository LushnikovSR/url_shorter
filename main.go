package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Response struct {
	Message string `json:"message"`
}

type SafeMap struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewSafeMap(initialSize int) *SafeMap {
	return &SafeMap{
		data: make(map[string]string, initialSize),
	}
}

func (sm *SafeMap) Set(key, value string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

func (sm *SafeMap) Get(key string) (string, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	val, ok := sm.data[key]
	return val, ok
}

func main() {

	safeMap := NewSafeMap(1000)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var resp Response
		resp.Message = "Hello"
		jsonResponse, err := json.Marshal(resp)
		if err != nil {
			//возвращает назад текст и меняет "Content-Type" на "text/plain"
			//http.Error(w, err.Error(), http.StatusInternalServerError)

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return

		}
		w.Write(jsonResponse)
	})

	http.HandleFunc("/name", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		name := query.Get("q")
		w.Header().Set("Content-Type", "application/json")
		var resp Response
		resp.Message = fmt.Sprintf("Hello %s!", name)
		jsonResponse, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Write(jsonResponse)
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		key := query.Get("k")
		value := query.Get("v")
		if safeMap.data[key] != "" {
			fail := fmt.Sprintf("This key:%s exists, try a different key or delete the existing key!", key)
			w.Write([]byte(fail))
			return
		}
		safeMap.Set(key, value)
		w.Write([]byte("Data successfully added"))
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		key := query.Get("k")
		value, exists := safeMap.Get(key)
		if !exists {
			text := fmt.Sprintf("Key: '%s' doesn`t exist", key)
			w.Write([]byte(text))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		var resp Response
		resp.Message = value
		jsonResponse, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Write(jsonResponse)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
