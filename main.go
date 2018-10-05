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

	router.GET("/articles/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.tmpl.html", gin.H{"navStates": defaultNav})
		}
		c.HTML(http.StatusOK, "article.tmpl.html", gin.H{
			"navStates": defaultNav,
			"article": gin.H{
				"referenceGenuine": true,
				"src":              fmt.Sprintf("/static/videos/%d.mp4", id),
				"next":             fmt.Sprintf("/articles/%d", id+1),
			},
		})
	})

	router.Run(":" + port)
}
