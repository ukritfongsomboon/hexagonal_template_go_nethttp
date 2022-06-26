package main

import (
	"context"
	"hexagonal/handler"
	"hexagonal/repository"
	"hexagonal/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	// fmt.Println("print init")
}

// func initDatabase() *sqlx.DB {
// 	db, err := sqlx.Open("mysql", "root:admin1234@tcp(203.151.199.184:6033)/mysqltest")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return db
// }

func initDatabaseMongoDB() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:ce19886d86cc2d57c835f38d479411ccc664df7d468b93c4b2239ff276ef519a@203.151.199.184:57030"))
	if err != nil {
		panic(err)
	}
	return client.Database("MongoDB")
}

func main() {
	// db := initDatabase()
	db := initDatabaseMongoDB()

	// customerRepositoryDB := repository.NewCustomerRepositoryDB(db)      // Data layer
	// customerService := service.NewCustomerService(customerRepositoryDB) // Business logic
	// fmt.Println(customerService.GetCustomers())

	// # User
	// userRepositoryMock := repository.NewUserRepositoryMock()
	userRepositoryDB := repository.NewUserRepositoryDB(db)  // # Database Layer
	userService := service.NewUserService(userRepositoryDB) // # Business Layer
	userHandler := handler.NewUserHandler(userService)      // # Presentation layer

	// # Create Api Service
	app := fiber.New()

	app.Get("/user", userHandler.GetUsers)
	app.Get(`/user/:userid/account`, userHandler.GetUser)

	// # Start Api Service
	app.Listen("localhost:3000")
}
