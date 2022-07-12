package service

import (
	"encoding/json"
	"fmt"
	"hexagonal/caching"
	"hexagonal/repository"
	"hexagonal/utils"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	cache    caching.AppCache
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository, cache caching.AppCache) UserService {
	return userService{userRepo: userRepo, cache: cache}
}

func (s userService) GetUsers(p repository.PaginationUser) (*DataResponsePagination, error) {
	// TODO Query Redis
	data, err := s.cache.Get(fmt.Sprintf("user:*:page:%v:row:%v", p.Page, p.Row))
	if err != nil {
		if err.Error() != "cache: no documents in result" {
			utils.LogError(err)
			return nil, utils.HandlerError{
				Code:    500,
				Message: "unexpected error",
			}
		}
	}

	var toHandler DataResponsePagination

	toHandler.Pagination.Page = p.Page
	toHandler.Pagination.Row = p.Row
	toHandler.Pagination.Total = 0

	usersResponses := []UserRes{}
	var users []repository.User

	if data == nil {
		users, err = s.userRepo.GetAll(repository.Pagination{
			Page: toHandler.Pagination.Page,
			Row:  toHandler.Pagination.Row,
		})
		if err != nil {
			utils.LogError(err)
			return nil, utils.HandlerError{
				Code:    500,
				Message: "unexpected error",
			}
		}

		// # DTO Data Tranfer Object
		for _, user := range users {
			userRes := UserRes{
				UserID: user.UserID,
				Email:  user.Email,

				Name: user.Name,
				Role: user.Role,
			}
			usersResponses = append(usersResponses, userRes)

		}
		// TODO Set Cache
		json, err := json.Marshal(usersResponses)
		if err != nil {
			utils.LogError(err)
			return nil, utils.HandlerError{
				Code:    500,
				Message: "unexpected error",
			}
		}

		err = s.cache.Set(fmt.Sprintf("user:*:page:%v:row:%v", p.Page, p.Row), string(json))
		if err != nil {
			utils.LogError(err)
			return nil, utils.HandlerError{
				Code:    500,
				Message: "unexpected error",
			}
		}

	} else {
		json.Unmarshal([]byte(*data), &usersResponses)
	}

	toHandler.Items = usersResponses

	countUser, err := s.userRepo.CountAll()
	if err != nil {
		utils.LogError(err)
		return nil, utils.HandlerError{
			Code:    500,
			Message: "unexpected error",
		}
	}

	toHandler.Pagination.Total = int(countUser)

	return &toHandler, nil
}

func (s userService) GetUser(userid string) (*UserRes, error) {
	// TODO Query Redis
	data, err := s.cache.Get("user:" + userid)
	if err != nil {
		if err.Error() != "cache: no documents in result" {
			utils.LogError(err)
			return nil, utils.HandlerError{
				Code:    500,
				Message: "unexpected error",
			}
		}
	}
	var user *repository.User
	var usersResponses UserRes
	if data == nil {
		user, err = s.userRepo.GetById(userid)
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

		usersResponses = UserRes{
			UserID: user.UserID,
			Email:  user.Email,
			// Password: user.Password,
			Name: user.Name,
			Role: user.Role,
			// Status:   user.Status,
		}

		// TODO Set Cache
		json, err := json.Marshal(usersResponses)
		if err != nil {
			utils.LogError(err)
			return nil, utils.HandlerError{
				Code:    500,
				Message: "unexpected error",
			}
		}

		err = s.cache.Set("user:"+userid, string(json))
		if err != nil {
			utils.LogError(err)
			return nil, utils.HandlerError{
				Code:    500,
				Message: "unexpected error",
			}
		}
	} else {
		json.Unmarshal([]byte(*data), &usersResponses)
	}

	return &usersResponses, nil
}

func (s userService) CreateUser(r AddUserReq) (*UserRes, error) {
	// TODO 1.validate data
	// TODO 1.1 Validate Data format
	// TODO 1.2 ValiData exists

	// # version 1
	// emailStatus, err := s.userRepo.CheckEmial(strings.ToLower(r.Email))
	// if err != nil {
	// 	return utils.HandlerError{
	// 		Code:    fiber.StatusBadRequest,
	// 		Message: "unexpected error",
	// 	}
	// }

	// if *emailStatus {
	// 	return utils.HandlerError{
	// 		Code:    fiber.StatusBadRequest,
	// 		Message: "email is exists",
	// 	}
	// }

	// nameStatus, err := s.userRepo.CheckName(r.Name)
	// if err != nil {
	// 	return utils.HandlerError{
	// 		Code:    fiber.StatusBadRequest,
	// 		Message: "unexpected error",
	// 	}
	// }

	// if *nameStatus {
	// 	return utils.HandlerError{
	// 		Code:    fiber.StatusBadRequest,
	// 		Message: "name is exists",
	// 	}
	// }

	// # version 2
	// TODO 2.Generate new password use bcryp
	newPass, err := bcrypt.GenerateFromPassword([]byte(r.Password), 10)
	if err != nil {
		return nil, utils.HandlerError{
			Code:    500,
			Message: "unexpected error",
		}
	}

	// TODO 3.make payload to repositiry
	data := repository.UserRecive{
		Name:     r.Name,
		Email:    strings.ToLower(r.Email),
		Password: string(newPass),
		Status:   false,
		Role:     1,
	}

	// TODO 4.insert to db
	_, err = s.userRepo.Create(data)
	// TODO 5.response
	if err != nil {
		if err.Error() == "email already exist" {
			utils.LogError(err)
			return nil, utils.HandlerError{
				Code:    400,
				Message: "email already exist",
			}
		} else {
			utils.LogError(err)
			return nil, utils.HandlerError{
				Code:    500,
				Message: "unexpected error",
			}
		}

	}

	return nil, nil
}

func (s userService) EditUser(string) (*UserRepository, error) {
	return nil, nil
}

func (s userService) Authentication(payload *AuthenReq) (*AuthenRes, error) {
	// TODO 1.Recive Email and Password
	// TODO 1.1 Validate Email and Password

	// TODO 2.get Email From Database
	user, err := s.userRepo.GetByEmail(strings.ToLower(payload.Email))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, utils.HandlerError{
				Code:    401,
				Message: "username or password is incorrect",
			}
		}
		// # Tech Error
		utils.LogError(err)
		return nil, utils.HandlerError{
			Code:    500,
			Message: "unexpected error",
		}
	}

	// # DTO
	// usersResponses := UserResponse{
	// 	UserID:   user.UserID,
	// 	Email:    user.Email,
	// 	Password: user.Password,
	// 	Name:     user.Name,
	// 	Role:     user.Role,
	// 	Status:   user.Status,
	// }

	// TODO 3.Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return nil, utils.HandlerError{
			Code:    401,
			Message: "username or password is incorrect",
		}

	}

	// TODO 4.Generate New Jwt
	//# Create Claim JWT Token
	// claims := jwt.StandardClaims{
	// 	Issuer:    user.UserID,
	// 	ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	// }

	// //# Create Header Claim
	// jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// //# Get signature JWT
	// // signature := utils.Getenv("APP_SIGNATURE", "ukritFongsomboon")
	// signature := "xxx"

	// //# Create JWT Token
	// token, err := jwtToken.SignedString([]byte(signature))

	private := viper.GetString("app.access_token_private_key")
	token, err := utils.CreateToken(30*time.Minute, user.UserID, private)

	//# DTO
	t := AuthenRes{
		Accesstoken: token,
		Status:      user.Status,
		Name:        user.Name,
		Email:       strings.ToLower(user.Email),
		Role:        user.Role,
	}

	// TODO 5.Return To handler
	return &t, nil
}
