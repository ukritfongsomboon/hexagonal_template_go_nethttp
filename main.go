package main

import (
	"context"
	"fmt"
	"hexagonal/caching"
	"hexagonal/handler"
	"hexagonal/middleware"
	"hexagonal/repository"
	"hexagonal/service"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// # Initial Config
func init() {
	initConfig()
	initTimeZone()
}

// func initDatabase() *sqlx.DB {
// 	db, err := sqlx.Open("mysql", "root:admin1234@tcp(203.151.199.184:6033)/mysqltest")
// 	if err != nil {
// 		panic(err)
// 	}
// 	return db
// }

// # Initial Environment
func initDatabaseMongoDB() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(viper.GetString("mongodb.uri")))
	if err != nil {
		panic(err)
	}

	// # Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		panic(err)
	}

	return client.Database(viper.GetString("mongodb.dbname"))
}

func initCacheRedis() *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Refs https://pkg.go.dev/github.com/go-redis/redis#ParseURL
	opt, err := redis.ParseURL(viper.GetString("redis.uri"))
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)

	// # Check the connection
	_, err = client.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	return client
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// # Default value ////////////////////////////
	viper.SetDefault("app.port", 3000)
	viper.SetDefault("app.env", "production")
	// # //////////////////////////////////////////

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}

	time.Local = ict
}

func xx(c *fiber.Ctx) error {
	// fmt.Println(c.Locals("user_id"))
	return c.Next()
}

// # Main Function
func main() {
	db := initDatabaseMongoDB()
	cacheRedis := initCacheRedis()

	// # User
	// userRepositoryMock := repository.NewUserRepositoryMock()

	cache := caching.NewAppCache(cacheRedis) // # Data Layer

	userRepositoryDB := repository.NewUserRepositoryDB(db)         // # Data Layer
	userService := service.NewUserService(userRepositoryDB, cache) // # Business Layer
	userHandler := handler.NewUserHandler(userService)             // # Presentation layer

	// for i := 901; i <= 2000; i++ {
	// 	userService.CreateUser(service.AddUserReq{
	// 		Name:     "User-" + strconv.Itoa(i),
	// 		Email:    "user-" + strconv.Itoa(i) + "@gmail.com",
	// 		Password: "User-" + strconv.Itoa(i),
	// 	})
	// }

	// # Create Api Service
	app := fiber.New()
	app.Use(cors.New(middleware.Cors))
	app.Get("/user", middleware.ValidToken, xx, userHandler.GetUsers)            //! Authen
	app.Get(`/user/:userid/account`, middleware.ValidToken, userHandler.GetUser) //! Authen
	app.Post(`/signup`, userHandler.SignUp)                                      //# No Authen

	app.Post("/signin", userHandler.SignIn) //# No Authen

	// # Set Production Mode or Developer Mode
	host := "0.0.0.0"
	mode := strings.ToLower(viper.GetString("app.env"))
	if mode != "production" && mode != "prd" {
		host = "localhost"
	}

	// # Start Api Service
	app.Listen(fmt.Sprintf("%v:%v", host, viper.GetInt("app.port")))
}
