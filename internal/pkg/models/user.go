package models

// Информация о пользователе
type User struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname" valid:"notnull"`
	About    string `json:"about"`
	Email    string `json:"email" valid:"notnull"`
}

// Информация для обновления пользователя
type UserUpdate struct {
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}
