// Created by Kun Li(likunarmstrong@gmail.com) on 03/08/2014.
// All rights reserved.

package main

import (
	"appengine"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getPosts(r *http.Request) []*Post {
	c := appengine.NewContext(r)
	var posts []*Post
	var fname, path, link, title, dateString string
	fileInfos, err := ioutil.ReadDir(pathPrefix)
	if err != nil {
		c.Errorf("open file failed: %v", err)
		return posts
	}

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

		doc, err := goquery.NewDocumentFromReader(content)
		if err != nil {
			c.Errorf("create new goquery document failed: %v", err)
			continue
		}

		/* Title */
		title = doc.Find("h1.title").Text()
		// Invalid title
		if len(title) == 0 {
			c.Errorf("Wrong title for post: %v", fname)
			continue
		}

		/* Link */
		selection := doc.Find("body > div > h1 + p > a.external[href]")
		if selection.Length() > 0 {
			link, _ = selection.Attr("href")
		} else {
			link = fname
		}
		c.Infof("link is %v", link)

		dateString = strings.Join(strings.Split(fname, "-")[:3], "-")
		posts = append(posts, &Post{
			FileName: fname,
			Title:    title,
			Date:     dateString,
			Link:     link,
		})
		c.Infof("Process file %s:\n", fname)
		c.Infof("date is %v", dateString)
		c.Infof("title is %v", title)
	}

	// Descending order based on publish time
	return posts
}

func byte2html(b []byte) template.HTML {
	return template.HTML(string(b))
}
