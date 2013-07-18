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
	http.HandleFunc("/blog", blogHandler)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Infof("Visitor from country:", r.Header.Get("X-AppEngine-CountryCountry"))
	c.Infof("Visitor from region:", r.Header.Get("X-AppEngine-RegionName"))
	c.Infof("Visitor from city:", r.Header.Get("X-AppEngine-CityName"))
	c.Infof("Visitor from location:", r.Header.Get("X-AppEngine-CityLatLong"))

	// Enable cache
	w.Header().Set("cache-control", "public")
	w.Header().Set("max-age", "7200")

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

func blogHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t := template.Must(template.ParseFiles("templates/blog_main.html"))

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
