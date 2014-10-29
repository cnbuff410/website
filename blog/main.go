// Created by Kun Li(likunarmstrong@gmail.com) on 03/08/2014.
// All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"appengine"
	"appengine/datastore"
)

const (
	pathPrefix string = "static/posts"
	postKind          = "post_meta"
)

func init() {
	r := mux.NewRouter()
	s := r.PathPrefix("/blog").Subrouter()

	// API
	s.HandleFunc("/update", postUpdateHandler).Methods("GET", "POST")
	s.HandleFunc("/all", postFetchMetaHandler).Methods("GET")
	s.HandleFunc("/{link}", postFetchContentHandler).Methods("GET").Queries("mode", "raw")

	// Serve page
	s.HandleFunc("/{link}", postServeHandler).Methods("GET")
	s.HandleFunc("/", blogServeHandler).Methods("GET")

	http.Handle("/", s)
}

func blogServeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/blog.html")
}

func postServeHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Infof("In post serve")
	http.ServeFile(w, r, "web/post.html")
}

func postUpdateHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	posts := getPostsMeta(r)
	if len(posts) == 0 {
		fmt.Fprintf(w, "There is no article to update")
		return
	}
	var keys []*datastore.Key
	for _, post := range posts {
		keys = append(keys, datastore.NewKey(c, postKind, post.Link, 0, nil))
	}
	if _, err := datastore.PutMulti(c, keys, posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Update complete: %v posts", len(posts))
}

func postFetchMetaHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var posts []Post
	_, err := datastore.NewQuery(postKind).Order("-Date").GetAll(c, &posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "%v", string(b))
}

func postFetchContentHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Infof("In post fetch")
	vars := mux.Vars(r)
	filename := vars["link"]

	content, err := getPostContent(filename)
	if err != nil {
		c.Errorf("Read content from %v error: %v", filename, err.Error())
		return
	}
	c.Infof("%v", content)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, content)
}
