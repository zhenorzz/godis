package router

import (
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type Cache map[string]string
func (cache Cache)Router(w http.ResponseWriter, r *http.Request) {
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
				fmt.Fprintf(w, "%s al ready exist", k)
				return
			}
			cache[k] = v
		}
		fmt.Println(cache)
		break
	case "GET":
		fmt.Fprintf(w, "%s : %s", pathInfo[1], cache[pathInfo[1]])
		break
	}

	fmt.Println(r.Method, pathInfo)
}