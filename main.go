package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/speps/go-hashids/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type jsonResult struct {
	Msg  string `json:"msg"`
	Succ bool   `json:"succ"`
	Url  string `json:"url"`
}

type shortUrl struct {
	Id  int64 `gorm:"primarykey"`
	Url string
}

type config struct {
	Id   uint `gorm:"primarykey"`
	Salt string
	Min  int
}

const IDKEY = "__short_url_id_key__"
const URLKEY = "id:%s"

func main() {

	origin := os.Getenv("origin")
	if len(origin) == 0 {
		panic(fmt.Errorf("origin url can not be empty"))
	}

	info, _ := os.Stat("./data/")
	if info == nil {
		os.Mkdir("./data/", os.ModePerm)
	}

	db, err := gorm.Open(sqlite.Open("./data/url.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&shortUrl{})
	db.AutoMigrate(&config{})

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
	h, _ := hashids.NewWithData(hd)

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			url := r.FormValue("url")
			if len(url) == 0 {
				result(rw, false, "url error", "")
				return
			}
			data := shortUrl{
				Url: url,
			}
			db.Create(&data)
			shortId, err := h.EncodeInt64([]int64{data.Id})
			if err != nil {
				result(rw, false, "get short id error, try again", "")
				return
			}
			result(rw, true, "success", fmt.Sprintf("%s/%s", origin, shortId))
		} else {
			shortId := strings.TrimLeft(r.URL.Path, "/")
			if len(shortId) > 1 {
				ids := h.DecodeInt64(shortId)
				var shortUrl shortUrl
				db.Find(&shortUrl, ids)
				if len(shortUrl.Url) > 0 {
					http.Redirect(rw, r, shortUrl.Url, http.StatusFound)
				}
			}
			rw.WriteHeader(http.StatusNotFound)
		}
	})

	http.ListenAndServe(":80", nil)
}

func result(rw http.ResponseWriter, succ bool, msg string, url string) {
	json, _ := json.Marshal(jsonResult{
		Succ: succ,
		Msg:  msg,
		Url:  url,
	})
	rw.Write(json)
}
