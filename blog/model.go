package main

import "html/template"

// Post represent a post
type Post struct {
	FileName string `json:"filename"`
	Title    string `json:"title"`
	Date     string `json:"date"`
	Link     string `json:"link"`
}

// PostContent represent the html content of a post
type PostContent struct {
	Content template.HTML
}
