package models

type Status struct {
	User   int32 `json:"user" valid:"notnull"`
	Forum  int32 `json:"forum" valid:"notnull"`
	Thread int32 `json:"thread" valid:"notnull"`
	Post   int64 `json:"post" valid:"notnull"`
}
