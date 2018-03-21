package router

import (
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"math/rand"
)

type Cache map[string]string
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

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
	case "PUT":
		result := cache.put(r)
		fmt.Fprint(w, result)
		fmt.Println(cache)
		break
	case "DELETE":
		result := cache.delete(r)
		fmt.Fprint(w, result)
		fmt.Println(cache)
		break
	default:
		fmt.Fprint(w, "does not support this method")
	}

}

//get resource
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
		return  key + " is not exist"
	}
}

//add resource
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
	for _, v := range response {
		//if _, ok := cache[k]; ok {
		//	return k + "%s already exist"
		//}
		cache[randSeq(10)] = v
	}
	return "success"
}

//update resource
func (cache Cache) put(r *http.Request) string {
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
			cache[k] = v
		} else {
			return k + " is not exist"
		}
	}
	return "success"
}

//delete resource
func (cache Cache) delete(r *http.Request) string {
	uri := r.RequestURI
	pathInfo := strings.Split(uri, "/")
	if len(pathInfo) == 0 {
		return "no path info"
	}
	key := pathInfo[1]
	if _, ok := cache[key]; ok {
		delete(cache,key)
	} else {
		return  key + " is not exist"
	}
	return "success"
}
