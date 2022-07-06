package main

import (
	"context"
	"hexagonal/caching"
	"hexagonal/handler"
	"hexagonal/middleware"
	"hexagonal/repository"
	"hexagonal/service"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

// # Initial Environment /////////////////////////////////////////////
func initDatabaseMongoDB() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:ce19886d86cc2d57c835f38d479411ccc664df7d468b93c4b2239ff276ef519a@203.151.199.184:57030"))
	if err != nil {
		panic(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		panic(err)
	}

	return client.Database("MongoDB")
}

func initCacheRedis() *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     "203.151.199.184:9463",
		Password: "admin1234", // no password set
		DB:       0,           // use default DB
	})

	// Check the connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return client
}

//// # ///////////////////////////////////////////////////////////////

func main() {
	db := initDatabaseMongoDB()
	cacheRedis := initCacheRedis()

	// # User
	// userRepositoryMock := repository.NewUserRepositoryMock()

	cache := caching.NewAppCache(cacheRedis) // # Data Layer

	userRepositoryDB := repository.NewUserRepositoryDB(db)         // # Data Layer
	userService := service.NewUserService(userRepositoryDB, cache) // # Business Layer
	userHandler := handler.NewUserHandler(userService)             // # Presentation layer

	// for i := 20198; i < 20200; i++ {
	// 	userService.CreateUser(service.AddUserReq{
	// 		Name:     "User-" + strconv.Itoa(i),
	// 		Email:    "user-" + strconv.Itoa(i) + "@gmail.com",
	// 		Password: "User-" + strconv.Itoa(i),
	// 	})
	// }

	// # Create Api Service
	app := fiber.New()
	app.Use(cors.New(middleware.Cors))
	app.Get("/user", userHandler.GetUsers)                //! Authen
	app.Get(`/user/:userid/account`, userHandler.GetUser) //! Authen
	app.Post(`/user`, userHandler.CreateUser)             //# No Authen

	app.Post("/login", userHandler.Login) //# No Authen

	// # Start Api Service
	app.Listen("localhost:3000")
}
