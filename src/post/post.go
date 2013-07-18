/*
 * Put your project description here
 *
 * Author: likunarmstrong@gmail.com
 */

package posts

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
	author   string
	visitCnt int64
	Tags     []string
	Created  time.Time
	Last_Mod time.Time
}

func generate_key(p string) string {
	key := strings.Replace(p, " ", "_", -1)
	key = strings.ToLower(key)

	// TODO:
	// I really want to limit this to like 80 chars or somthing but I
	// havent found anything yet to do  that easily

	return key
}

func New(subject, content string) Entity {
	return Entity{Subject: subject,
		Content:  content,
		Key:      generate_key(subject),
		Created:  time.Now(),
		Last_Mod: time.Now()}
}

func parentKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Posts", "postsTable", 0, nil)
}

func (e Entity) Put(c appengine.Context) {
	key := datastore.NewKey(c, "Entity", e.Key, 0, parentKey(c))
	datastore.Put(c, key, &e)
}

func Get(c appengine.Context, id string) (Entity, error) {
	q := datastore.NewQuery("Entity").
		Ancestor(parentKey(c)).
		Filter("Key =", id)

	t := q.Run(c)
	var e Entity
	_, err := t.Next(&e)

	if err != nil {
		return e, err
	}

	return e, nil
}

func GetLatest(c appengine.Context) ([]Entity, error) {
	q := datastore.NewQuery("Entity").
		Ancestor(parentKey(c)).
		Order("-Created").
		Limit(10)

	var e []Entity
	for t := q.Run(c); ; {
		var ne Entity
		_, err := t.Next(&ne)
		if err == datastore.Done {
			break
		}
		if err != nil {
			return e, err
		}
		e = append(e, ne)
	}
	return e, nil
}
