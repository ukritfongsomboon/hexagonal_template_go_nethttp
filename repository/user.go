package repository

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Oauth struct {
	Provider string `json:"provider" bson:"provider" db:"provider"`
	Id       string `json:"id" bson:"id" db:"id"`
	Email    string `json:"email" bson:"email" db:"email"`
	Password string `json:"password" bson:"password"`
}

type User struct {
	UserID      string    `json:"user_id" bson:"user_id" db:"user_id"`
	Email       string    `json:"email" bson:"email" db:"email"`
	Password    string    `json:"password" bson:"password"`
	Name        string    `json:"name" bson:"name" db:"name"`
	CreatedDate time.Time `json:"create_date" bson:"create_date" db:"create_date"`
	LastUpdate  time.Time `json:"update_date" bson:"update_date" db:"update_date"`
	Role        int       `json:"role" bson:"role" db:"role"`
	Status      bool      `json:"status" bson:"status" db:"status"`
	Oauth       []Oauth   `json:"oauth" bson:"oauth" db:"oauth"`
}
type UserRecive struct {
	Email    string `json:"email" bson:"email" db:"email"`
	Password string `json:"password" bson:"password"`
	Name     string `json:"name" bson:"name" db:"name"`
	Role     int    `json:"role" bson:"role" db:"role"`
	Status   bool   `json:"status" bson:"status" db:"status"`
}

// # POD
type UserRepository interface {
	GetAll() ([]User, error)
	GetById(string) (*User, error)
	Add(UserRecive) (*mongo.UpdateResult, error)
	Edit(string, UserRecive) (*mongo.UpdateResult, error)
}
