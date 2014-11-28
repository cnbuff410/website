// Created by Kun Li(likunarmstrong@gmail.com) on 03/08/2014.
// All rights reserved.

package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"appengine"

	"github.com/PuerkitoBio/goquery"
)

func getPostsMeta(r *http.Request) []*Post {
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
		selection := doc.Find("body > div > h1 + h2 > a.external[href]")
		if selection.Length() > 0 {
			link, _ = selection.Attr("href")
		} else {
			link = fname
		}

		dateString = strings.Join(strings.Split(fname, "-")[:3], "-")
		posts = append(posts, &Post{
			FileName: fname[0:(len(fname) - len(filepath.Ext(fname)))],
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

func getPostContent(title string) (string, error) {
	fname := title + ".html"
	filePath := filepath.Join(pathPrefix, fname)
	content, _ := os.Open(filePath)
	defer content.Close()

	doc, err := goquery.NewDocumentFromReader(content)
	if err != nil {
		return "", err
	}

	body, err := doc.Find("body").Html()
	if err != nil {
		return body, err
	}
	return body, nil
}

func byte2html(b []byte) template.HTML {
	return template.HTML(string(b))
}
