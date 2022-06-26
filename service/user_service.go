package service

import (
	"errors"
	"fmt"
	"hexagonal/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return userService{userRepo: userRepo}
}

func (s userService) GetUsers() ([]UserResponse, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// # DTO Data Tranfer Object
	usersResponses := []UserResponse{}
	for _, user := range users {
		userRes := UserResponse{
			UserID:   user.UserID,
			Email:    user.Email,
			Password: user.Password,
			Name:     user.Name,
			Role:     user.Role,
			Status:   user.Status,
		}
		usersResponses = append(usersResponses, userRes)

	}
	return usersResponses, nil
}

func (s userService) GetUser(userid string) (*UserResponse, error) {
	user, err := s.userRepo.GetById(userid)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	usersResponses := UserResponse{
		UserID:   user.UserID,
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		Role:     user.Role,
		Status:   user.Status,
	}
	return &usersResponses, nil
}

func (s userService) AddUser(string) (*UserRepository, error) {
	return nil, nil
}

func (s userService) EditUser(string) (*UserRepository, error) {
	return nil, nil
}

func (s userService) Authen(AuthenReq) (*AuthenRes, error) {
	return nil, nil
}
