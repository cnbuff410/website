package main

import (
	"appengine"
	"appengine/datastore"
)

func GlobalKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Pages", "PageTable", 0, nil)
}
