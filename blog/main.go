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
	"appengine/memcache"
)

const (
	pathPrefix         = "static/posts"
	postKind           = "post_meta"
	allPostKey         = "allPost"
	postContentKeyBase = "post_"
)

func init() {
	r := mux.NewRouter()
	s := r.PathPrefix("/blog").Subrouter()

	// API
	s.HandleFunc("/update", updatePostCacheHandler).Methods("POST")
	s.HandleFunc("/all", fetchPostMetaHandler).Methods("GET")
	s.HandleFunc("/{link}", fetchPostContentHandler).Methods("GET").Queries("mode", "raw")

	// Serve page
	s.HandleFunc("/{link}", postHandler).Methods("GET")
	s.HandleFunc("/", blogHandler).Methods("GET")

	http.Handle("/", s)
}

func blogHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/blog.html")
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/post.html")
}

func updatePostCacheHandler(w http.ResponseWriter, r *http.Request) {
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

	// Update memcache
	value, err := json.Marshal(posts)
	if err != nil {
		c.Errorf("error marshal posts: %v", err)
		return
	}
	item := &memcache.Item{
		Key:   allPostKey,
		Value: value,
	}
	if err := memcache.Add(c, item); err == memcache.ErrNotStored {
		if err := memcache.Set(c, item); err != nil {
			c.Errorf("error setting item: %v", err)
		}
	} else if err != nil {
		c.Errorf("error adding item: %v", err)
	}

	var content, postContentKey string
	for _, post := range posts {
		postContentKey = postContentKeyBase + post.FileName
		content, err = getPostContent(post.FileName)
		if err != nil {
			c.Errorf("Read content from %v error: %v", post.FileName, err)
			continue
		}
		item := &memcache.Item{
			Key:   postContentKey,
			Value: []byte(content),
		}
		if err := memcache.Add(c, item); err == memcache.ErrNotStored {
			if err := memcache.Set(c, item); err != nil {
				c.Errorf("error setting item: %v", err)
			}
		} else if err != nil {
			c.Errorf("error adding item: %v", err)
		}
	}

	fmt.Fprintf(w, "Update datastore complete: %v posts", len(posts))
}

func fetchPostMetaHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	var value []byte
	// Get the item from the memcache
	if item, err := memcache.Get(c, allPostKey); err == memcache.ErrCacheMiss {
		c.Infof("info of all posts not in the cache")
		var posts []Post
		_, err := datastore.NewQuery(postKind).Order("-Date").GetAll(c, &posts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		value, err = json.Marshal(posts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Add the item to the memcache, if the key does not already exist
		item := &memcache.Item{
			Key:   allPostKey,
			Value: value,
		}
		if err := memcache.Add(c, item); err == memcache.ErrNotStored {
			c.Infof("item with key %q already exists", item.Key)
		} else if err != nil {
			c.Errorf("error adding item: %v", err)
		}
	} else if err != nil {
		c.Errorf("error getting item: %v", err)
	} else {
		c.Infof("Hit the cache")
		value = item.Value
	}

	w.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "%v", string(value))
}

func fetchPostContentHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	vars := mux.Vars(r)
	filename := vars["link"]

	var content string
	memcacheKey := postContentKeyBase + filename
	// Get the item from the memcache
	if item, err := memcache.Get(c, memcacheKey); err == memcache.ErrCacheMiss {
		c.Infof("post content not in the cache")
		content, err = getPostContent(filename)
		if err != nil {
			c.Errorf("Read content from %v error: %v", filename, err.Error())
			return
		}
		// Add the item to the memcache, if the key does not already exist
		item := &memcache.Item{
			Key:   memcacheKey,
			Value: []byte(content),
		}
		if err := memcache.Add(c, item); err == memcache.ErrNotStored {
			c.Infof("item with key %q already exists", item.Key)
		} else if err != nil {
			c.Errorf("error adding item: %v", err)
		}
	} else if err != nil {
		c.Errorf("error getting item: %v", err)
	} else {
		c.Infof("Hit the cache")
		content = string(item.Value)
	}

	c.Infof("%v", content)
	w.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprint(w, content)
}
