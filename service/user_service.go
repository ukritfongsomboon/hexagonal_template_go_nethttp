package service

import (
	"hexagonal/repository"
	"hexagonal/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
		utils.LogError(err)
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
			return nil, utils.HandlerError{
				Code:    200,
				Message: "user not found",
			}
		}
		// # Tech Error
		utils.LogError(err)
		return nil, utils.HandlerError{
			Code:    500,
			Message: "unexpected error",
		}
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

func (s userService) CreateUser(r AddUserReq) error {
	// TODO 1.validate data
	// TODO 1.1 Validate Data format
	// TODO 1.2 ValiData exists
	emailStatus, err := s.userRepo.CheckEmial(r.Email)
	if err != nil {
		return utils.HandlerError{
			Code:    fiber.StatusBadRequest,
			Message: "unexpected error",
		}
	}

	if *emailStatus {
		return utils.HandlerError{
			Code:    fiber.StatusBadRequest,
			Message: "email is exists",
		}
	}

	nameStatus, err := s.userRepo.CheckName(r.Name)
	if err != nil {
		return utils.HandlerError{
			Code:    fiber.StatusBadRequest,
			Message: "unexpected error",
		}
	}

	if *nameStatus {
		return utils.HandlerError{
			Code:    fiber.StatusBadRequest,
			Message: "name is exists",
		}
	}

	// TODO 2.Generate new password use bcryp
	newPass, err := bcrypt.GenerateFromPassword([]byte(r.Password), 10)
	if err != nil {
		return utils.HandlerError{
			Code:    500,
			Message: "unexpected error",
		}
	}

	// TODO 3.make payload to repositiry
	data := repository.UserRecive{
		Name:     r.Name,
		Email:    r.Email,
		Password: string(newPass),
		Status:   false,
		Role:     1,
	}

	// TODO 4.insert to db
	_, err = s.userRepo.Create(data)

	// TODO 5.response
	if err != nil {
		utils.LogError(err)
		return utils.HandlerError{
			Code:    500,
			Message: "unexpected error",
		}
	}
	return nil
}

func (s userService) EditUser(string) (*UserRepository, error) {
	return nil, nil
}

func (s userService) Authentication(AuthenReq) (*AuthenRes, error) {
	return nil, nil
}
