package main

import (
	"fmt"
	"net/http"
	"log"
	"strings"
	"encoding/json"
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
			var response map[string]string
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(body, &response)
			if err != nil {
				panic(err)
			}
			if len(response) == 0 {
				fmt.Fprintf(w, "%v", false)
				return
			}
			for k, v := range response {
				if _, ok := cache[k]; ok {
					fmt.Fprintf(w, "%s already exist", k)
					return
				}
				cache[k] = v
			}
			fmt.Println(cache)
			break
		case "GET":
			fmt.Fprintf(w, "%s : %s", pathInfo[2], cache[pathInfo[2]])
			break
		}

		fmt.Println(r.Method, pathInfo)
	})
	log.Fatal(http.ListenAndServe(":8070", nil))
}
