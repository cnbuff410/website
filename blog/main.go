// Created by Kun Li(likunarmstrong@gmail.com) on 03/08/2014.
// All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"appengine"
	"github.com/gorilla/mux"
)

const (
	fileKey = "FILE"
)

func init() {
	r := mux.NewRouter()
	s := r.PathPrefix("/blog").Subrouter()
	// Main page
	s.HandleFunc("/blog", blogMainHandler).Methods("GET")
	s.HandleFunc("/blog/{post}", blogPostHandler).Methods("GET")

	// Data fetching
	s.HandleFunc("/posts/{num}", fetchPostHandler).Methods("GET")
	http.Handle("/", s)
}

func blogMainHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/blog_index.html")
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

func fetchPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	num := vars["num"]
	c := appengine.NewContext(r)
	Posts, err := getPosts(r)
	if !strings.EqualFold(num, "all") && len(num) > 0 {
		n, e := strconv.Atoi(num)
		if e == nil {
			Posts = Posts[:n]
		}
	}
	if err != nil {
		c.Errorf("fetch post error: %s", err.Error())
		return
	}
	b, err := json.Marshal(Posts)
	if err != nil {
		c.Errorf("marshal json error: %s", err.Error())
		fmt.Fprint(w, "")
		return
	}
	fmt.Fprint(w, string(b))
}
