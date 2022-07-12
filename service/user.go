package service

import "hexagonal/repository"

type Oauth struct {
	Provider string `json:"provider" bson:"provider" db:"provider"`
	Id       string `json:"id" bson:"id" db:"id"`
	Email    string `json:"email" bson:"email" db:"email"`
	Password string `json:"password" bson:"password"`
}

type UserRepository struct {
	UserID   string  `json:"user_id" bson:"user_id" db:"user_id"`
	Email    string  `json:"email" bson:"email" db:"email"`
	Password string  `json:"password" bson:"password"`
	Name     string  `json:"name" bson:"name" db:"name"`
	Role     int     `json:"role" bson:"role" db:"role"`
	Status   bool    `json:"status" bson:"status" db:"status"`
	Oauth    []Oauth `json:"oauth" bson:"oauth" db:"oauth"`
}

type AuthenReq struct {
	Email    string `json:"email" bson:"email" db:"email"`
	Password string `json:"password" bson:"password"`
}

type AuthenRes struct {
	Accesstoken string `json:"accesstoken" bson:"accesstoken" db:"accesstoken"`
	Email       string `json:"email" bson:"email" db:"email"`
	Name        string `json:"name" bson:"name" db:"name"`
	Status      bool   `json:"status" bson:"status" db:"status"`
	Role        int    `json:"role" bson:"role" db:"role"`
}

type UserResponse struct {
	UserID   string `json:"user_id" bson:"user_id" db:"user_id"`
	Email    string `json:"email" bson:"email" db:"email"`
	Password string `json:"password" bson:"password"`
	Name     string `json:"name" bson:"name" db:"name"`
	Role     int    `json:"role" bson:"role" db:"role"`
	Status   bool   `json:"status" bson:"status" db:"status"`
}

type UserRes struct {
	UserID string `json:"user_id" bson:"user_id" db:"user_id"`
	Email  string `json:"email" bson:"email" db:"email"`
	Name   string `json:"name" bson:"name" db:"name"`
	Role   int    `json:"role" bson:"role" db:"role"`
}

type AddUserReq struct {
	Name     string `json:"name" bson:"name" db:"name"`
	Email    string `json:"email" bson:"email" db:"email"`
	Password string `json:"password" bson:"password"`
}

type Pagination struct {
	Page  int `json:"page" bson:"page"`
	Row   int `json:"row" bson:"row"`
	Total int `json:"total" bson:"total"`
}

type DataResponsePagination struct {
	Items      interface{} `json:"item" bson:"item"`
	Pagination Pagination  `json:"pagination" bson:"pagination"`
}

type UserService interface {
	GetUsers(repository.PaginationUser) (*DataResponsePagination, error)
	GetUser(string) (*UserRes, error)
	CreateUser(AddUserReq) (*UserRes, error)
	EditUser(string) (*UserRepository, error)
	Authentication(*AuthenReq) (*AuthenRes, error)
}
