package router

import (
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type Cache map[string]string

func (cache Cache) Router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		result := cache.post(r)
		fmt.Fprint(w, result)
		fmt.Println(cache)
		break
	case "GET":
		result := cache.get(r)
		fmt.Fprint(w, result)
		fmt.Println(cache)
		break
	}

}

func (cache Cache) post(r *http.Request) string {
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
		return "none json in body"
	}
	for k, v := range response {
		if _, ok := cache[k]; ok {
			return fmt.Sprintf("%s already exist", k)
		}
		cache[k] = v
	}
	return "success"
}

func (cache Cache) get(r *http.Request) string{
	uri := r.RequestURI
	pathInfo := strings.Split(uri, "/")
	if len(pathInfo) == 0 {
		return "no path info"
	}
	key := pathInfo[1]
	if _, ok := cache[key]; ok {
		return cache[key]
	} else {
		return fmt.Sprintf("%s is not exist", key)
	}
}
