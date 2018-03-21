package main

import (
	"net/http"
	"log"
	"router"
	"math/rand"
)

func main() {
	cache := make(router.Cache)
	ch := make(chan map[string]string)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cache.Router(w, r, ch)
	})
	go consumer(cache, ch)
	log.Fatal(http.ListenAndServe(":8070", nil))
}
func consumer(cache router.Cache, messages <-chan map[string]string) {
	for message := range messages {
		for _, v := range message {
			cache[randSeq(10)] = v
		}
	}

}
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}