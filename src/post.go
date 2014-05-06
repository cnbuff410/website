// Created by Kun Li(likunarmstrong@gmail.com) on 03/08/2014.
// All rights reserved.

package main

import (
	"fmt"
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
	pathPrefix    string = "web/html/posts"
	datePreFormat string = "2006-01-02"
	dateFormat    string = "Mon, Jan _2 2006"
)

// Post represent an article
type Post struct {
	Key      string
	Title    string
	Date     string
	Link     string
	VisitCnt int64
	Tags     []string
}

func getPosts(c appengine.Context, r *http.Request) ([]*Post, error) {
	fileInfos, err := ioutil.ReadDir(pathPrefix)
	if err != nil {
		return nil, err
	}
	var posts []*Post
	var fname, path, link, title, datePreString string
	var datePre time.Time
	var size int64
	var attr html.Attribute
	var f func(*html.Node)
	for i := len(fileInfos) - 1; i >= 0; i-- {
		if fileInfos[i].IsDir() {
			continue
		}
		fname = fileInfos[i].Name()
		path = filepath.Join(pathPrefix, fname)
		content, _ := os.Open(path)
		defer content.Close()
		doc, err := html.Parse(content)
		if err != nil {
			// Invalid post content, simply ignore now
			continue
		}
		c.Infof("File %s:\n", fname)

		// Extract title
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && strings.EqualFold(n.Data, "h1") && len(n.Attr) > 0 {
				for i := 0; i < len(n.Attr); i++ {
					attr = n.Attr[i]
					if strings.EqualFold(attr.Key, "class") && strings.EqualFold(attr.Val, "title") {
						title = n.FirstChild.Data
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)

		// Extract link
		link = fmt.Sprintf("%v/%v", r.URL.Path, fname)
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
