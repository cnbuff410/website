// Created by Kun Li(likunarmstrong@gmail.com) on 03/08/2014.
// All rights reserved.

package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"appengine"

	"github.com/gorilla/mux"
)

const (
	fileKey = "FILE"
)

func init() {
	r := mux.NewRouter()

	// Main page
	r.HandleFunc("/", mainHandler).Methods("GET")
	r.HandleFunc("/blog", blogMainHandler).Methods("GET")
	r.HandleFunc("/blog/{post}", blogPostHandler).Methods("GET")

	http.Handle("/", r)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	// Enable cache
	w.Header().Set("cache-control", "public")
	w.Header().Set("max-age", "7200")

	t := template.Must(template.ParseFiles(
		"web/html/main.html",
		"web/html/chrome/head.html",
		"web/html/chrome/foot.html",
		"web/html/chrome/nav.html"))

	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func blogMainHandler(w http.ResponseWriter, r *http.Request) {
	// Enable cache
	w.Header().Set("cache-control", "public")
	w.Header().Set("max-age", "7200")

	t := template.Must(template.ParseFiles(
		"web/html/blog_main.html",
		"web/html/chrome/head.html",
		"web/html/chrome/foot.html",
		"web/html/chrome/nav-blog.html"))

	c := appengine.NewContext(r)
	Posts, err := getPosts(c, r)
	Posts = Posts[:15]
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, Posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func blogPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["post"]
	content, err := ioutil.ReadFile(filepath.Join(pathPrefix, fileName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t := template.Must(template.ParseFiles(
		"web/html/blog_post.html",
		"web/html/chrome/head.html",
		"web/html/chrome/foot.html",
		"web/html/chrome/nav-blog.html"))

	if err := t.Execute(w, &PostContent{Content: byte2html(content)}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
