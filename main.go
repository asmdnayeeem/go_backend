package main

import (
	"example/api/models"
	"example/api/storage"
	"log"
	"os"

	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Respository struct {
	DB *gorm.DB
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSlMode:  os.Getenv("DB_SSLMODE"),
	}
	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	err = models.MigrateUser(db)
	if err != nil {
		log.Fatal(err)

	}
	r := Respository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)

	app.Listen(":8080")
}

func (r *Respository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create", r.CreateUser)
	api.Get(("/showuser"), r.ShowUser)
	api.Post(("/deluser"), r.DelUser)
	api.Post(("/updateuser"), r.UpdateUser)
}

func (r *Respository) UpdateUser(context *fiber.Ctx) error {
	user := User{}
	err := context.BodyParser(&user)
	if err != nil {
		return err
	}
	err = r.DB.Where("username = ?", user.Username).Updates(&user).Error
	if err != nil {
		log.Fatal(err)
		return err
	}
	context.Status(http.StatusOK)
	return context.JSON(user)
}

func (r *Respository) DelUser(context *fiber.Ctx) error {
	user := User{}
	err := context.BodyParser(&user)
	if err != nil {
		return err
	}	
	err = r.DB.Where("username = ?", user.Username).Delete(&user).Error
	if err != nil {
		log.Fatal(err)
		return err
	}
	context.Status(http.StatusOK)
	return context.JSON(user)
}

func (r *Respository) ShowUser(context *fiber.Ctx) error {
	var user []User
	r.DB.Find(&user)
	context.Status(http.StatusOK)
	return context.JSON(user)
}

func (r *Respository) CreateUser(context *fiber.Ctx) error {
	user := User{}
	err := context.BodyParser(&user)
	if err != nil {
		return err
	}

	err = r.DB.Create((&user)).Error
	if err != nil {
		log.Fatal(err)
		return err
	}
	context.Status(http.StatusOK)
	return context.JSON(user)
}
