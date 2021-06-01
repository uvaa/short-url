package main

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {

	info, _ := os.Stat("./data/")
	if info == nil {
		os.Mkdir("./data/", os.ModePerm)
	}

	var err error
	db, err = gorm.Open(sqlite.Open("./data/url.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&shortUrl{})
	db.AutoMigrate(&config{})

}
