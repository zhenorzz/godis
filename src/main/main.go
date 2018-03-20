package main

import (
	"fmt"
	"net/http"
	"log"
	"strings"
	"io/ioutil"
)

var cache map[string]string

func main() {
	cache = make(map[string]string)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		pathInfo := strings.Split(uri, "/")
		switch r.Method {
		case "POST":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			value := string(body)
			cache[pathInfo[2]] = value
			break
		case "GET":
			fmt.Fprintf(w, "%s : %s", pathInfo[2], cache[pathInfo[2]])
			break
		}

		fmt.Println(r.Method, pathInfo)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
