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
