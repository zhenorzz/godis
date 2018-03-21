package main

import (
	"net/http"
	"log"
	"router"
)


func main() {
	cache := make(router.Cache)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cache.Router(w, r)
	})
	log.Fatal(http.ListenAndServe(":8070", nil))
}
