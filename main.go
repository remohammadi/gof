package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

type navStates struct {
	InHome  bool
	InAbout bool
}

var (
	homeNav    = &navStates{InHome: true, InAbout: false}
	aboutNav   = &navStates{InHome: false, InAbout: true}
	defaultNav = &navStates{InHome: false, InAbout: false}
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	loadedArticles, err := loadArticles("articles")
	if err != nil {
		log.Fatal("failed to load articles: ", err)
	}
	log.Printf("%d article(s) are loaded.", loadedArticles)

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{"navStates": homeNav})
	})

	router.GET("/about", func(c *gin.Context) {
		c.HTML(http.StatusOK, "about.tmpl.html", gin.H{"navStates": aboutNav})
	})

	router.GET("/articles/:id/", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id < 1 || id > len(articles) {
			c.HTML(http.StatusBadRequest, "error.tmpl.html", gin.H{"navStates": defaultNav})
			return
		}
		article := articles[id-1]
		obj := gin.H{
			"navStates": defaultNav,
			"article":   article,
			"progress":  fmt.Sprintf("%d%%", (id*100)/len(articles)),
		}
		c.HTML(http.StatusOK, "article.tmpl.html", obj)
	})

	router.GET("/articles/:id/details/", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id < 1 || id > len(articles) {
			c.HTML(http.StatusBadRequest, "error.tmpl.html", gin.H{"navStates": defaultNav})
			return
		}
		article := articles[id-1]
		obj := gin.H{
			"navStates": defaultNav,
			"article":   article,
		}
		if id < len(articles) {
			obj["next"] = fmt.Sprintf("/articles/%d", id+1)
		}
		c.HTML(http.StatusOK, "article-details.tmpl.html", obj)
	})

	router.POST("/submit/", func(c *gin.Context) {
		articleID, err := strconv.Atoi(c.PostForm("article-id"))
		if err != nil || articleID < 1 || articleID > len(articles) {
			c.HTML(http.StatusBadRequest, "error.tmpl.html", gin.H{"navStates": defaultNav})
			return
		}
		if articleID < len(articles) {
			c.Redirect(303, fmt.Sprintf("/articles/%d/", articleID+1))
		} else {
			c.Redirect(303, "/scores/")
		}
	})

	router.GET("/scores/", func(c *gin.Context) {
		c.HTML(http.StatusBadRequest, "error.tmpl.html", gin.H{"navStates": defaultNav})
	})

	router.Run(":" + port)
}
