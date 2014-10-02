package main

import "html/template"

// Post represent a post
type Post struct {
	FileName string
	Title    string
	Date     string
	Link     string
}

// PostContent represent the html content of a post
type PostContent struct {
	Content template.HTML
}
