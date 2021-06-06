package models

import "time"

type Thread struct {
	Id      int32     `json:"id"`
	Title   string    `json:"title" valid:"notnull"`
	Author  string    `json:"author" valid:"notnull"`
	Forum   string    `json:"forum"`
	Message string    `json:"message" valid:"notnull"`
	Votes   int32     `json:"votes"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"created"`
}

type ThreadSearch struct {
	Forum  string    `json:"slug"`
	Limit int32     `json:"limit"`
	Since time.Time `json:"since"`
	Desc  bool      `json:"desc"`
}
