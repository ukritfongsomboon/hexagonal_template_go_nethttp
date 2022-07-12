package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// # Adapter
type userRepositoryDB struct {
	db *mongo.Database
}

// # Contructor Adapter
func NewUserRepositoryDB(db *mongo.Database) userRepositoryDB {
	// // TODO Migrate the schema
	// db.AutoMigrate(&User{})
	return userRepositoryDB{db: db}
}

func (r userRepositoryDB) GetAll(p Pagination) ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// if p.Page > 1 {
	// 	p.Page = p.Page - 1
	// }

	if p.Row == 0 {
		p.Row = 10
	}

	// refs https://www.codementor.io/@arpitbhayani/fast-and-efficient-pagination-in-mongodb-9095flbqr
	query := bson.A{
		bson.D{{"$skip", p.Row * (p.Page - 1)}},
		bson.D{{"$limit", p.Row}},
		bson.D{{"$sort", bson.D{{"create_date", 1}}}},
	}
	result := []User{}
	cursor, err := r.db.Collection("hexagonal_users").Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r userRepositoryDB) GetById(UserID string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := bson.D{{"user_id", UserID}}
	result := User{}
	err := r.db.Collection("hexagonal_users").FindOne(ctx, query).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r userRepositoryDB) Create(payload UserRecive) (*User, error) {

	// # version 1

	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()

	// new_user_id := uuid.New().String()

	// filter := bson.D{{"user_id", new_user_id}}
	// update := bson.D{{"$set", bson.D{{"user_id", new_user_id}, {"role", 0}, {"status", 0}, {"email", payload.Email}, {"name", payload.Name}, {"password", payload.Password}, {"create_date", time.Now()}, {"update_date", time.Now()}}}}
	// opts := options.Update().SetUpsert(true)

	// cursor, err := r.db.Collection("hexagonal_users").UpdateOne(ctx, filter, update, opts)
	// if err != nil {
	// 	return nil, err
	// }

	// return cursor, nil

	// # New version 2

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	user := User{
		CreatedDate: time.Now(),
		LastUpdate:  time.Now(),
		Email:       payload.Email,
		Name:        payload.Name,
		Password:    payload.Password,
		Status:      false,
		Role:        0,
		UserID:      uuid.New().String(),
		Oauth:       []Oauth{},
	}

	// db.hexagonal_users_v2.createIndex({"email":1,"user_id":1},{unique:true})
	// https://codevoweb.com/golang-mongodb-jwt-authentication-authorization/
	// https://kb.objectrocket.com/mongo-db/how-to-create-an-index-using-the-golang-driver-for-mongodb-455
	// https://hafizhabdurrachman.medium.com/how-to-improve-api-part-2-caching-redis-written-in-golang-e9644691b931

	res, err := r.db.Collection("hexagonal_users").InsertOne(ctx, &user)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("email already exist")
		}
		return nil, err
	}

	// TODO Create a unique index for the email field

	// refs https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#IndexView.CreateMany
	// TODO Declare an array of bsonx models for the indexes
	models := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "user_id", Value: 1}}, Options: options.Index().SetUnique(true),
		},
	}

	if _, err := r.db.Collection("hexagonal_users").Indexes().CreateMany(ctx, models); err != nil {
		// TODO แก้ไข เมื่อ setindex ไม่ได้ ไม่ต้อง errorออกมา หรือแก้เป็น log warning
		return nil, errors.New("could not create index")
	}

	var newUser User
	query := bson.M{"_id": res.InsertedID}

	err = r.db.Collection("hexagonal_users").FindOne(ctx, query).Decode(&newUser)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (r userRepositoryDB) Update(UserID string, UserRecive UserRecive) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	filter := bson.D{{"user_id", UserID}}
	update := bson.D{{"$set", bson.D{{"email", UserRecive.Email}, {"role", UserRecive.Role}, {"status", UserRecive.Status}, {"name", UserRecive.Name}, {"password", UserRecive.Password}, {"update_date", time.Now()}}}}

	cursor, err := r.db.Collection("hexagonal_users").UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

func (r userRepositoryDB) Delete(string) error {
	return nil
}

func (r userRepositoryDB) CheckEmial(email string) (*bool, error) {
	// TODO ใช้สำหรับตรวจสอบ Email ในระบบว่ามีอยู่หรือไม่ Return True เมื่อ มีในระบบ False เมื่อไม่มีในระบบ
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	filter := bson.D{{"email", email}}
	count, err := r.db.Collection("hexagonal_users").CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	var status bool
	if count > 0 {
		status = true
		return &status, nil
	}
	return &status, nil
}

func (r userRepositoryDB) CheckName(name string) (*bool, error) {
	// TODO ใช้สำหรับตรวจสอบ Username ในระบบว่ามีอยู่หรือไม่ Return True เมื่อ มีในระบบ False เมื่อไม่มีในระบบ
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	filter := bson.D{{"name", name}}
	count, err := r.db.Collection("hexagonal_users").CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	var status bool
	if count > 0 {
		status = true
		return &status, nil
	}
	return &status, nil
}

func (r userRepositoryDB) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	query := bson.D{{"email", email}}
	result := User{}
	err := r.db.Collection("hexagonal_users").FindOne(ctx, query).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r userRepositoryDB) CountAll() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	filter := bson.D{{"status", true}}
	count, err := r.db.Collection("hexagonal_users").CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, err
}
