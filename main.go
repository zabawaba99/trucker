package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	var (
		port           string
		firebaseURL    string
		firebaseSecret string
	)

	flag.StringVar(&port, "port", stringOr(os.Getenv("PORT"), "8080"), "The port that the application will run on")
	flag.StringVar(&firebaseURL, "firebase-url", os.Getenv("FIREBASE_URL"), "The URL used to communicate with Firebase")
	flag.StringVar(&firebaseSecret, "firebase-secret", os.Getenv("FIREBASE_SECRET"), "The secret used to authenticate with Firebase")
	flag.Parse()

	a, err := NewAPI(firebaseURL, firebaseSecret)
	if err != nil {
		log.Fatal(err)
	}

	index, err := ioutil.ReadFile("assets/index.html")
	if err != nil {
		log.Fatal(err)
	}

	r := httprouter.New()
	r.GET("/assets/*file", serveFile)
	r.GET("/", func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		w.Header().Add("Content-Type", "text/html")
		w.Write(index)
	})

	r.GET("/api/entries", a.listEntries)
	r.POST("/api/entries", a.saveEntry)

	err = http.ListenAndServe(":"+port, r)
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
