package models

type User struct {
	ID             int    `json:"id"`
	Name           string `json:"name" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
	Password       string `json:"password" validate:"required,min=6"`
	Username       string `json:"username" validate:"required"`
	ProfilePicture string `json:"profile_picture"`
}
