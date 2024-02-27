package model

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey" example:"1"`
	Username string `json:"username" gorm:"unique" binding:"required" example:"JohnDoe"`
	Password string `json:"password" example:"xxxxxxx"`
}
