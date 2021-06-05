package models

type Tread struct {
	Id int32 `json:"id"`
	Title string `json:"title" valid:"notnull"`
	Author string `json:"author" valid:"notnull"`
	Forum string `json:"forum"`
	Message string `json:"message" valid:"notnull"`
	Votes int32 `json:"votes"`
	Slug string `json:"slug"`
	Created string `json:"created" valid:"notnull"`
}

