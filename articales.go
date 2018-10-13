package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type ArticleType int

const (
	articleTypeVide    = 0
	articleTypeYouTube = 1
)

type ArticleData struct {
	ID int `json:"-,"`

	Title     string      `json:"title"`
	Type      ArticleType `json:"type"`
	IsGenuine bool        `json:"isGenuine"`
	Src       string      `json:"src"`
	ShortDesc string      `json:"shortDescription"`
	LongDesc  string      `json:"longDescription"`
}

var articles = make([]*ArticleData, 0)

func loadArticles(dir string) (int, error) {
	i := 1
	for ; ; i++ {
		filename := path.Join(dir, fmt.Sprintf("%d.json", i))
		fileInfo, err := os.Stat(filename)
		if err != nil {
			if os.IsNotExist(err) {
				break
			}
			return 0, fmt.Errorf("error while checking article %d: %s", i, err)
		}
		dat, err := ioutil.ReadFile(path.Join(dir, fileInfo.Name()))
		if err != nil {
			return 0, fmt.Errorf("error while loading article %d: %s", i, err)
		}
		article := ArticleData{}
		err = json.Unmarshal(dat, &article)
		if err != nil {
			return 0, fmt.Errorf("error while unmarshalling article %d: %s", i, err)
		}
		article.ID = i
		articles = append(articles, &article)
	}

	return i - 1, nil
}
