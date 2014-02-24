package main

import (
	"html/template"
	"net/http"

	"appengine"
)

const (
	fileKey = "FILE"
)

func init() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/blog", blogHandler)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Enable cache
	w.Header().Set("cache-control", "public")
	w.Header().Set("max-age", "7200")

	r.ParseForm()
	if r.Method == "GET" {
		c.Infof("path: %s", r.URL.Path)
		c.Infof("scheme: %s", r.URL.Scheme)

		t := template.Must(template.ParseFiles("frontend/index.html"))

		if err := t.Execute(w, nil); err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
	}
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("frontend/blog_main.html"))

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
