package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	index, err := ioutil.ReadFile("assets/index.html")
	if err != nil {
		log.Fatal(err)
	}

	r := httprouter.New()
	r.GET("/assets/*file", serveFile)
	r.GET("/index", func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		w.Header().Add("Content-Type", "text/html")
		w.Write(index)
	})

	r.GET("/api/entries", listEntries)

	err = http.ListenAndServe(":8081", r)
	if err != nil {
		log.Fatal(err)
	}
}

func serveFile(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	file := p.ByName("file")

	b, err := ioutil.ReadFile("assets/" + file)
	if err != nil {
		log.Printf("Could not read file %s\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var contentType string
	if strings.HasSuffix(file, ".html") {
		contentType = "text/html"
	} else if strings.HasSuffix(file, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(file, ".js") {
		contentType = "application/javascript"
	} else {
		log.Printf("Could not determine content type %s\n", file)
	}

	w.Header().Add("Content-Type", contentType)
	w.Write(b)
}

type entry struct {
	Broker     string    `json:"broker"`
	TripNumber int       `json:"trip_number"`
	Date       time.Time `json:"date"`
}

func listEntries(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")

	e := []entry{
		entry{Broker: "Joe", TripNumber: 1, Date: <-time.After(time.Second)},
		entry{Broker: "Toe", TripNumber: 2, Date: <-time.After(time.Second)},
		entry{Broker: "Moe", TripNumber: 3, Date: <-time.After(time.Second)},
	}

	if err := json.NewEncoder(w).Encode(&e); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error encoding: %s\n", err)
	}
}
