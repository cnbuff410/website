// Created by Kun Li(likunarmstrong@gmail.com) on 03/08/2014.
// All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"appengine"
	"appengine/datastore"
	"github.com/gorilla/mux"
)

const (
	fileKey           = "FILE"
	pathPrefix string = "static/posts"
	postKind          = "post_meta"
)

func init() {
	r := mux.NewRouter()
	s := r.PathPrefix("/blog").Subrouter()

	// Main page
	s.HandleFunc("/", blogMainHandler).Methods("GET")

	// API
	s.HandleFunc("/update", postUpdateHandler).Methods("GET", "POST")
	s.HandleFunc("/all", postFetchMetaHandler).Methods("GET")
	s.HandleFunc("/posts/{link}", postFetchContentHandler).Methods("GET")

	http.Handle("/", s)
}

func blogMainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello")
	//http.ServeFile(w, r, "web/blo")
}

func postUpdateHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	posts := getPosts(r)
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
	var postTitles []string
	var postLink []string
	postMetaMap := make(map[string][]string)
	for _, post := range posts {
		postTitles = append(postTitles, post.Title)
		postLink = append(postLink, post.Link)
	}
	postMetaMap["title"] = postTitles
	postMetaMap["link"] = postLink
	b, err := json.Marshal(postMetaMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "%v", string(b))
}

func postFetchContentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	num := vars["num"]
	c := appengine.NewContext(r)
	Posts := getPosts(r)
	if !strings.EqualFold(num, "all") && len(num) > 0 {
		n, e := strconv.Atoi(num)
		if e == nil {
			Posts = Posts[:n]
		}
	}
	b, err := json.Marshal(Posts)
	if err != nil {
		c.Errorf("marshal json error: %s", err.Error())
		fmt.Fprint(w, "")
		return
	}
	fmt.Fprint(w, string(b))
	//http.ServeFile(w, r, "web/post.html")
}
