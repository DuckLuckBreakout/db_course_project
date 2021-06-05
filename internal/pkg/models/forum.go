package models

type Forum struct {
	Title   string `json:"title" valid:"notnull"`
	User    string `json:"user" valid:"notnull"`
	Slug    string `json:"slug" valid:"notnull"`
	Posts   int64  `json:"posts"`
	Threads int32  `json:"treads"`
}
