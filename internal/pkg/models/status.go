package models

type Status struct {
	User  int32 `json:"user" valid:"notnull"`
	Forum int32 `json:"forum" valid:"notnull"`
	Tread int32 `json:"tread" valid:"notnull"`
	Post  int64 `json:"post" valid:"notnull"`
}
