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
	"appengine/datastore"

	"code.google.com/p/go.net/html"
	"github.com/PuerkitoBio/goquery"
)

const (
	pathPrefix    string = "web/posts"
	datePreFormat string = "2006-01-02"
	dateFormat    string = "Mon, Jan _2 2006"
	previewLength int    = 200
)

// Post represent a post
type Post struct {
	Key     string
	Title   string
	Date    string
	Link    string
	Preview string
	Tags    []string
}

// PostContent represent the html content of a post
type PostContent struct {
	Content template.HTML
}

func getPosts(r *http.Request) ([]*Post, error) {
	c := appengine.NewContext(r)
	fileInfos, err := ioutil.ReadDir(pathPrefix)
	if err != nil {
		return nil, err
	}
	var posts []*Post
	var fname, path, link, title, datePreString, preview string
	var datePre time.Time
	var size int64

	for i := len(fileInfos) - 1; i >= 0; i-- {
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
		node, err := html.Parse(content)
		if err != nil {
			// Invalid post content, simply ignore now
			c.Errorf("error parsing the content of post: %v", err.Error())
			continue
		}
		c.Infof("Process file %s:\n", fname)

		doc := goquery.NewDocumentFromNode(node)

		/* Title */
		title = doc.Find("h1.title").Text()
		// Invalid title
		if len(title) == 0 {
			c.Errorf("Wrong title for post: %v", fname)
			continue
		}

		/* Link and preview */
		preview = ""
		selection := doc.Find("body > div > h1 + p > a.external[href]")
		if selection.Length() > 0 {
			link, _ = selection.Attr("href")
			preview = "PDF"
		} else {
			link = fmt.Sprintf("%v/%v", r.URL.Path, fname)
			previewSel := doc.Find("body > div > p")
			for i := 0; i < previewSel.Length(); i++ {
				preview += previewSel.Get(i).FirstChild.Data + "\n"
				if len(preview) > previewLength {
					break
				}
			}
		}
		c.Infof("preview is %v", preview)

		size = fileInfos[i].Size()
		datePreString = strings.Join(strings.Split(fname, "-")[:3], "-")
		datePre, _ = time.Parse(datePreFormat, datePreString)
		posts = append(posts, &Post{
			Key:     fmt.Sprintf("%s-%d", fname, size),
			Title:   title,
			Preview: preview,
			Date:    datePre.Format(dateFormat),
			Link:    link,
		})
	}

	// Descending order based on publish time
	return posts, nil
}

func globalKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Pages", "PageTable", 0, nil)
}

func byte2html(b []byte) template.HTML {
	return template.HTML(string(b))
}
