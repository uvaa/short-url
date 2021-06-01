package main

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
