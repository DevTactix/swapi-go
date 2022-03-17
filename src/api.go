package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const PORT = ":8080"
const API_BASE_URL = "https://swapi.dev/api/"

type RouteHandler func(http.ResponseWriter, *http.Request)
type Route struct {
	Slug    string
	Handler RouteHandler
}

type Character struct {
	Name   string
	Height string
	Mass   string
}

var routes = []Route{
	{Slug: "/", Handler: rootRoute},
	{Slug: "/darth", Handler: curryApiRoute("people/4")},
	{Slug: "/luke", Handler: curryApiRoute("people/1")},
}

func rootRoute(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(rw, r)
		return
	}
	io.WriteString(rw, "Hello world!")
}

func curryApiRoute(slug string) RouteHandler {
	log.Println("curryApiRoute:", slug)
	return func(rw http.ResponseWriter, r *http.Request) {
		apiCall := fmt.Sprintf("%s%s", API_BASE_URL, slug)
		log.Println(apiCall)
		fmt.Fprintln(rw, apiCall)

		if r, err := http.Get(apiCall); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			log.Println("Bad API request:", err)
		} else {
			fmt.Println(r)
			if data, err := ioutil.ReadAll(r.Body); err != nil {
				log.Println("API call failed:", err)
			} else {
				var jsonString Character
				json.Unmarshal(data, &jsonString)
				fmt.Fprintf(rw, "Name: %s \nHeight: %s\nMass: %s\n", jsonString.Name, jsonString.Height, jsonString.Mass)
			}
		}
	}

}

func initRoutes(routes []Route) {
	for _, route := range routes {
		http.HandleFunc(route.Slug, route.Handler)
	}
}

func main() {
	initRoutes(routes)
	fmt.Println("Server running on localhost", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
