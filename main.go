package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/jinzhu/gorm"
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
	var err error

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("$DATABASE_URL must be set")
	}
	var db *gorm.DB
	if strings.HasPrefix(dbURL, "sqlite://") {
		db, err = gorm.Open("sqlite3", dbURL[9:])
	} else {
		var rawDB gorm.SQLCommon
		rawDB, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatalf("Error opening database: %q", err)
		}
		db, err = gorm.Open("postgres", rawDB)
	}
	if err != nil {
		panic("failed to initiate database")
	}
	db.AutoMigrate(&GenuinityOpinion{})
	defer db.Close()

	loadedArticles, err := loadArticles("articles")
	if err != nil {
		log.Fatal("failed to load articles: ", err)
	}
	log.Printf("%d article(s) are loaded.", loadedArticles)

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
	router.StaticFile("favicon.ico", "static/favicon.ico")

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
		userFP, err := strconv.Atoi(c.PostForm("user-fp"))
		if err != nil || userFP < 0 {
			c.HTML(http.StatusBadRequest, "error.tmpl.html", gin.H{"navStates": defaultNav})
			return
		}
		userChoice, err := strconv.ParseBool(c.PostForm("user-choice"))
		if err != nil {
			c.HTML(http.StatusBadRequest, "error.tmpl.html", gin.H{"navStates": defaultNav})
			return
		}
		duration, err := strconv.Atoi(c.PostForm("duration"))
		if err != nil || duration < 0 {
			c.HTML(http.StatusBadRequest, "error.tmpl.html", gin.H{"navStates": defaultNav})
			return
		}

		ipAddr := c.Request.Header.Get("x-forwarded-for")
		if ipAddr != "" {
			list := strings.Split(ipAddr, ",")
			ipAddr = list[len(list)-1]
		} else {
			ipAddr = c.Request.RemoteAddr
		}

		opinion := GenuinityOpinion{
			ArticleID:  uint16(articleID),
			UserID:     c.PostForm("user-id"),
			UserFP:     uint32(userFP),
			UserChoice: userChoice,
			UserIP:     ipAddr,
			UserAgent:  c.Request.UserAgent(),
			Duration:   uint32(duration),
			IsCorrect:  userChoice == articles[articleID-1].IsGenuine,
		}

		q := db.Model(&GenuinityOpinion{}).
			Where("article_id = ? AND user_id = ?", opinion.ArticleID, opinion.UserID)
		if opinion.UserID == "" {
			q = q.Where("user_fp = ?")
		}
		var count int
		q.Count(&count)
		if count > 0 {
			now := gorm.NowFunc()
			opinion.DeletedAt = &now
		}
		db.Create(&opinion)

		if articleID < len(articles) {
			c.Redirect(303, fmt.Sprintf("/articles/%d/", articleID+1))
		} else {
			c.Redirect(303, fmt.Sprintf("/scores/%s/", opinion.UserID))
		}
	})

	router.GET("/scores/:user-id/", func(c *gin.Context) {
		userID := c.Param("user-id")
		opinions := make([]GenuinityOpinion, 0, len(articles))
		db.Find(&opinions, "user_id = ?", userID)
		total := len(opinions)
		if total == 0 {
			c.HTML(http.StatusBadRequest, "error.tmpl.html", gin.H{"navStates": defaultNav})
		}
		correct := 0
		for _, o := range opinions {
			if o.IsCorrect {
				correct++
			}
		}
		colors := make([]int, 11)
		for i := range colors {
			if i < correct {
				colors[i] = 0
			} else if i == correct {
				colors[i] = 1
			} else if i < 10 {
				colors[i] = 2
			} else {
				colors[i] = 3
			}
		}
		obj := gin.H{
			"navStates": defaultNav,
			"articles":  articles,
			"total":     total,
			"correct":   correct,
			"colors":    colors,
		}
		c.HTML(http.StatusOK, "scores.tmpl.html", obj)
	})

	router.Run(":" + port)
}
