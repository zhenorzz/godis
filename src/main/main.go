package main

import (
	"net/http"
	"log"
	"router"
	"storage"
	"sync"
	"fmt"
)

func main() {
	postChan := make(chan map[string]string)
	putChan := make(chan []string)
	delChan := make(chan string)
	cache := router.Cache{Data: make(map[string]string), Mutex: &sync.RWMutex{}}
	err := storage.Read(&cache)
	if err != nil {
		panic(err)
	}
	fmt.Println(cache)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			result := cache.Post(r, postChan)
			fmt.Fprint(w, result)
			break
		case "GET":
			result := cache.Get(r)
			fmt.Fprint(w, result)
			break
		case "PUT":
			result := cache.Put(r, putChan)
			fmt.Fprint(w, result)
			break
		case "DELETE":
			result := cache.Delete(r, delChan)
			fmt.Fprint(w, result)
			break
		default:
			fmt.Fprint(w, "does not support this method")
		}
	})

	//get message from channel
	go func() {
		for {
			select {
			case addData := <-postChan:
				storage.Write(addData)
			case updateData := <-putChan:
				storage.Update(updateData)
			case key := <-delChan:
				storage.Delete(key)
			}
		}
	}()

	log.Fatal(http.ListenAndServe(":8070", nil))
}
