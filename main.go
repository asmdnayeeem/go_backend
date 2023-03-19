package main

import (
	"example/api/models"
	"example/api/storage"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"

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

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func CheckPasswordH(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
	api.Post(("/login"), r.Login)

}

func (r *Respository) Login(context *fiber.Ctx) error {
	user := User{}
	err := context.BodyParser(&user)
	if err != nil {
		return err
	}
	var muser User
	r.DB.Where("username = ?", user.Username).First(&muser)
	fmt.Println(muser.Password)
	if CheckPasswordH(user.Password, muser.Password) {
		context.Status(http.StatusOK)
		return context.JSON(muser)
	}
	context.Status(http.StatusUnauthorized)
	return context.JSON(fiber.Map{"message": "unauthorized"})
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
	user.Password, _ = HashPassword(user.Password)
	err = r.DB.Create((&user)).Error
	if err != nil {
		log.Fatal(err)
		return err
	}
	context.Status(http.StatusOK)
	return context.JSON(user)
}
