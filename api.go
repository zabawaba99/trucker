package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/CloudCom/firego"
	"github.com/julienschmidt/httprouter"
)

type api struct {
	fb *firego.Firebase
}

func NewAPI(fbURL, fbSecret string) (*api, error) {
	if fbURL == "" || fbSecret == "" {
		return nil, errors.New("fbURL or fbSecret are empty")
	}
	fb := firego.New(fbURL)
	fb.Auth(fbSecret)
	return &api{fb: fb}, nil
}

func (a *api) saveEntry(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var e entry
	if err := json.NewDecoder(req.Body).Decode(&e); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("%s", err)
		return
	}

	if _, err := a.fb.Push(e); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s", err)
		return
	}
}

func (a *api) listEntries(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")

	var v map[string]entry
	if err := a.fb.Value(&v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var entries []entry
	for _, e := range v {
		entries = append(entries, e)
	}

	if err := json.NewEncoder(w).Encode(&entries); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error encoding: %s\n", err)
	}
}
