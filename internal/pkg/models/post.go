package models

import "time"

type Post struct {
	Id         int64     `json:"id"`
	Parent     int64     `json:"parent"`
	Thread     int32     `json:"thread"`
	Author     string    `json:"author"`
	ThreadSlug string    `json:"thread_slug"`
	Message    string    `json:"message"`
	Forum      string    `json:"forum"`
	IsEdited   bool      `json:"isEdited"`
	Created    time.Time `json:"created"`
}

type PostSearch struct {
	Thread     int32  `json:"thread"`
	ThreadSlug string `json:"slug"`
	Forum      string `json:"forum"`
	Sort       string `json:"sort"`
	Limit      int32  `json:"limit"`
	Since      int64  `json:"since"`
	Desc       bool   `json:"desc"`
}

type PostCommon struct {
	Id         int64     `json:"id"`
	Thread     int32     `json:"thread"`
	Author     string    `json:"author"`
	Message    string    `json:"message"`
	Forum      string    `json:"forum"`
	Created    time.Time `json:"created"`
}