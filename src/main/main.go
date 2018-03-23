package main

import (
	"net/http"
	"log"
	"router"
	"sync"
)

func main() {
	cache := router.Cache{Data:make(map[string]string),Mutex: &sync.Mutex{}}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cache.Router(w, r)
	})
	log.Fatal(http.ListenAndServe(":8070", nil))
}
