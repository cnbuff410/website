// Created by Kun Li(likunarmstrong@gmail.com) on 03/08/2014.
// All rights reserved.

package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	fileKey = "FILE"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainHandler).Methods("GET")
	r.HandleFunc("/blog", redirectBlogHandler).Methods("GET")
	http.Handle("/", r)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("cache-control", "public")
	//w.Header().Set("max-age", "7200")
	http.ServeFile(w, r, "./web/main.html")
}

func redirectBlogHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.URL.Path+"/", http.StatusSeeOther)
}
