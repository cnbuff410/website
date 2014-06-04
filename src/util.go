package main

import (
	"html/template"

	"appengine"
	"appengine/datastore"
)

func globalKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Pages", "PageTable", 0, nil)
}

func byte2html(b []byte) template.HTML {
	return template.HTML(string(b))
}
