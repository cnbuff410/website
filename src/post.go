/*
 * Put your project description here
 *
 * Author: likunarmstrong@gmail.com
 */

package main

import (
	"strings"
	"time"

	"appengine"
	"appengine/datastore"
)

type Post struct {
	Key      string
	Content  string
	Title    string
	Author   string
	VisitCnt int64
	Tags     []string
	Created  time.Time
	Last_Mod time.Time
}

/*type Comments struct {
	Comment []string
}*/

func New(title, content string) Post {
	return Post{
		Key:      generate_key(title),
		Content:  content,
		Title:    title,
		Author:   "Kun Li",
		VisitCnt: 0,
		Created:  time.Now(),
		Last_Mod: time.Now()}
}

func (p Post) Save(c appengine.Context) {
	key := datastore.NewKey(c, "Post", p.Key, 0, GlobalKey(c))
	datastore.Put(c, key, &p)
	// Check if it's modify or create
}

func Get(c appengine.Context, key string) (Post, error) {
	q := datastore.NewQuery("Entity").
		Ancestor(GlobalKey(c)).
		Filter("Key =", key)

	t := q.Run(c)
	var p Post
	_, err := t.Next(&p)

	if err != nil {
		return p, err
	}

	return p, nil
}

func GetLatest(c appengine.Context) ([]Post, error) {
	q := datastore.NewQuery("Entity").
		Ancestor(GlobalKey(c)).
		Order("-Created").
		Limit(10)

	var p []Post
	for t := q.Run(c); ; {
		var ne Post
		_, err := t.Next(&ne)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return p, err
		}
		p = append(p, ne)
	}
	return p, nil
}

func generate_key(p string) string {
	key := strings.Replace(p, " ", "-", -1)
	key = strings.ToLower(key)
	return key
}
