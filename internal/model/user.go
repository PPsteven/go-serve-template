package model

type User struct {
	ID uint `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique" binding:"required"`
	Password string `json:"password"`
}
