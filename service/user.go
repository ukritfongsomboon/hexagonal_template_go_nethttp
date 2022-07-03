package service

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
}

type UserResponse struct {
	UserID   string `json:"user_id" bson:"user_id" db:"user_id"`
	Email    string `json:"email" bson:"email" db:"email"`
	Password string `json:"password" bson:"password"`
	Name     string `json:"name" bson:"name" db:"name"`
	Role     int    `json:"role" bson:"role" db:"role"`
	Status   bool   `json:"status" bson:"status" db:"status"`
}

type AddUserReq struct {
	Name     string `json:"name" bson:"name" db:"name"`
	Email    string `json:"email" bson:"email" db:"email"`
	Password string `json:"password" bson:"password"`
}

type UserService interface {
	GetUsers() ([]UserResponse, error)
	GetUser(string) (*UserResponse, error)
	CreateUser(AddUserReq) error
	EditUser(string) (*UserRepository, error)
	Authentication(AuthenReq) (*AuthenRes, error)
}
