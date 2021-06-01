package main

import (
	"encoding/base64"
	"os"
	"strconv"

	"github.com/speps/go-hashids/v2"
)

var hash *hashids.HashID

func init() {

	var config config
	db.Find(&config)
	if config.Id == 0 {

		salt := os.Getenv("salt")
		if len(salt) == 0 {
			file, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
			bytes := make([]byte, 16)
			file.Read(bytes)
			salt = base64.StdEncoding.EncodeToString(bytes)
		}
		config.Salt = salt

		min, err := strconv.Atoi(os.Getenv("minlength"))
		if err != nil {
			min = 5
		}
		config.Min = min

		db.Save(&config)
	}

	hd := hashids.NewData()
	hd.Salt = config.Salt
	hd.MinLength = config.Min

	var err error
	hash, err = hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}
}
