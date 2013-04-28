package main

import (
	"appengine"
	"html/template"
	"net/http"
)

const (
	FILE_KEY = "FILE"
)

func init() {
	http.HandleFunc("/", mainHandler)
	//http.HandleFunc("/upload", uploadFileHandler)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	r.ParseForm()
	if r.Method == "GET" {
		c.Infof("path", r.URL.Path)
		c.Infof("scheme", r.URL.Scheme)

		t := template.Must(template.ParseFiles("templates/main.html"))

		if err := t.Execute(w, nil); err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
	}
}

func check(err error, w http.ResponseWriter, c appengine.Context) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		c.Errorf("%v", err)
	}
}
