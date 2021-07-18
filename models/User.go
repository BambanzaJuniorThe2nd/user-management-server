package models

import (
	"time"
)

type User struct {
	ID        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name,omitempty" bson:"name"`
	Email     string    `json:"email,omitempty" bson:"email"`
	Title     string    `json:"title,omitempty" bson:"title"`
	Birthdate time.Time `json:"birthdate,omitempty" bson:"birthdate"`
	Password  string    `json:"password,omitempty" bson:"password"`
	IsAdmin   bool      `json:"isAdmin" bson:"isAdmin"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt"`
}

type GetByTokenArgs struct {
	Token string `json:"token,omitempty"`
}

type CreateByAdminArgs struct {
	Name      string
	Email     string
	Title     string
	Birthdate string
	IsAdmin   bool
}

type UpdateByAdminArgs struct {
	CreateByAdminArgs
}

type LoginArgs struct {
	Email    string
	Password string
}

type LoginResult struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreateDefaultAdminArgs struct {
	Password string
	CreateByAdminArgs
}

type ChangePasswordArgs struct {
	Password string
}
