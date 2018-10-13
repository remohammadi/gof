package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type GenuinityOpinion struct {
	gorm.Model

	ArticleID  uint16
	UserID     string `gorm:"size:15"`
	UserFP     uint32
	UserChoice bool
	UserIP     string `gorm:"size:15"`
	UserAgent  string
	Duration   uint32
	IsCorrect  bool
}
