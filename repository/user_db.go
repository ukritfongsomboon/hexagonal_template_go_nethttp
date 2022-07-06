package repository

import (
	"context"
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

func (r userRepositoryDB) GetAll(p PaginationUser) ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if p.Perpage == 0 {
		p.Perpage = 20
	}
	query := bson.A{
		bson.D{{"$skip", p.Page}},
		bson.D{{"$limit", p.Perpage}},
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

func (r userRepositoryDB) Create(payload UserRecive) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	new_user_id := uuid.New().String()

	filter := bson.D{{"user_id", new_user_id}}
	update := bson.D{{"$set", bson.D{{"user_id", new_user_id}, {"role", 0}, {"status", 0}, {"email", payload.Email}, {"name", payload.Name}, {"password", payload.Password}, {"create_date", time.Now()}, {"update_date", time.Now()}}}}
	opts := options.Update().SetUpsert(true)

	cursor, err := r.db.Collection("hexagonal_users").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	return cursor, nil
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
