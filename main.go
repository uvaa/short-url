package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {

	origin := os.Getenv("origin")
	if len(origin) == 0 {
		panic(fmt.Errorf("origin url can not be empty"))
	}

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		// post request
		// save data and return the short url
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
			shortId, err := hash.EncodeInt64([]int64{data.Id})
			if err != nil {
				result(rw, false, "get short id error, try again", "")
				return
			}
			result(rw, true, "success", fmt.Sprintf("%s/%s", origin, shortId))
		} else {
			// other request(get, put and more)
			// get the short url data, and redirect to the origin url
			shortId := strings.TrimLeft(r.URL.Path, "/")
			if len(shortId) > 1 {
				ids, err := hash.DecodeInt64WithError(shortId)
				if err == nil {
					var shortUrl shortUrl
					db.Find(&shortUrl, ids)
					if shortUrl.Id > 0 {
						http.Redirect(rw, r, shortUrl.Url, http.StatusFound)
					}
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
