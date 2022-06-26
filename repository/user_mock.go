package repository

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type userRepositoryMock struct {
	users []User
}

func NewUserRepositoryMock() userRepositoryMock {
	users := []User{
		{UserID: "136e4d95-41a7-4d50-9c9d-e25f93fa406a", Email: "kobori4268@gmail.com", Password: "4813494d137e1631bba301d5acab6e7bb7aa74ce1185d456565ef51d737677b2", Name: "user1", Role: 0, Status: false},
		{UserID: "136e4d95-41a7-4d50-9c9d-e25f93fa406b", Email: "kobori4268@gmail.com", Password: "4813494d137e1631bba301d5acab6e7bb7aa74ce1185d456565ef51d737677b2", Name: "user2", Role: 0, Status: false},
	}
	return userRepositoryMock{users: users}
}

func (r userRepositoryMock) GetAll() ([]User, error) {
	return r.users, nil
}

func (r userRepositoryMock) GetById(userid string) (*User, error) {
	for _, user := range r.users {
		if userid == user.UserID {
			return &user, nil
		}
	}
	return nil, errors.New("mongo: no documents in result")
}

func (r userRepositoryMock) Add(UserRecive) (*mongo.UpdateResult, error) {
	return nil, nil
}

func (r userRepositoryMock) Edit(string, UserRecive) (*mongo.UpdateResult, error) {
	return nil, nil
}
func (r userRepositoryMock) Delete(string) error {
	return nil
}
