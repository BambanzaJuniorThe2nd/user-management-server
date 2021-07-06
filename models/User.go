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
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt"`
}
