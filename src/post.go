// Created by Kun Li(likunarmstrong@gmail.com) on 03/08/2014.
// All rights reserved.

package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"appengine"

	"code.google.com/p/go.net/html"
)

const (
	pathPrefix    string = "web/posts"
	datePreFormat string = "2006-01-02"
	dateFormat    string = "Mon, Jan _2 2006"
)

// Post represent a post
type Post struct {
	Key      string
	Title    string
	Date     string
	Link     string
	VisitCnt int64
	Tags     []string
}

// PostContent represent the html content of a post
type PostContent struct {
	Content template.HTML
}

func getPosts(c appengine.Context, r *http.Request) ([]*Post, error) {
	fileInfos, err := ioutil.ReadDir(pathPrefix)
	if err != nil {
		return nil, err
	}
	var posts []*Post
	var attr html.Attribute
	var f func(*html.Node)
	for i := len(fileInfos) - 1; i >= 0; i-- {
		var fname, path, link, title, datePreString string
		var datePre time.Time
		var size int64

		if fileInfos[i].IsDir() {
			continue
		}
		fname = fileInfos[i].Name()
		// Only process html file
		if !strings.EqualFold(filepath.Ext(fname), ".html") {
			continue
		}
		path = filepath.Join(pathPrefix, fname)
		content, _ := os.Open(path)
		defer content.Close()
		doc, err := html.Parse(content)
		if err != nil {
			// Invalid post content, simply ignore now
			c.Errorf("error parsing the content of post: %v", err.Error())
			continue
		}
		c.Infof("Process file %s:\n", fname)

		f = func(n *html.Node) {
			// Extract title
			if n.Type == html.ElementNode && strings.EqualFold(n.Data, "title") {
				if n.FirstChild != nil {
					title = n.FirstChild.Data
				}
			}
			// Extract link if it's just a pdf
			if n.Type == html.ElementNode && strings.EqualFold(n.Data, "body") {
				n1 := n.FirstChild.NextSibling
				if strings.EqualFold(n1.Data, "div") {
					n2 := n1.FirstChild.NextSibling.NextSibling.NextSibling
					if strings.EqualFold(n2.Data, "p") {
						n3 := n2.FirstChild
						if strings.EqualFold(n3.Data, "a") && len(n3.Attr) > 0 {
							n4 := n3.FirstChild
							if strings.EqualFold(n4.Data, "pdf") {
								for i := 0; i < len(n3.Attr); i++ {
									attr = n3.Attr[i]
									if strings.EqualFold(attr.Key, "href") {
										link = attr.Val
										break
									}
								}
							}
						}
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)

		// Invalid title
		if len(title) == 0 {
			c.Errorf("Wrong title for post: %v", fname)
			continue
		}
		// Extract link if it's not a pure pdf
		if len(link) == 0 {
			link = fmt.Sprintf("%v/%v", r.URL.Path, fname)
		}
		size = fileInfos[i].Size()
		datePreString = strings.Join(strings.Split(fname, "-")[:3], "-")
		datePre, _ = time.Parse(datePreFormat, datePreString)
		posts = append(posts, &Post{
			Key:      fmt.Sprintf("%s-%d", fname, size),
			Title:    title,
			Date:     datePre.Format(dateFormat),
			Link:     link,
			VisitCnt: 100})
	}

	// Descending order based on publish time
	return posts, nil
}
