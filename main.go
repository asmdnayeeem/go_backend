package main

import (
	"example/api/models"
	"example/api/storage"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isadmin"`
}

type Respository struct {
	DB *gorm.DB
}

// bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(

		[]byte(password), 14)
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
	api.Post(("/login"), r.Login)
	api.Use(jwtware.New(jwtware.Config{SigningKey: []byte("secret")}))
	api.Post("/:username/createuser", r.CreateUser)
	api.Get(("/showusers"), r.ShowUser)
	api.Post(("/:usesrname/deleteuser"), r.DelUser)
	api.Post(("/:username/updateuser"), r.UpdateUser)
	api.Get("/logout", logout)

}

// Controllers
// login
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

		claims := jwt.MapClaims{
			"name": user.Username,
			"exp":  time.Now().Add(time.Minute * 5).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return context.SendStatus(fiber.StatusInternalServerError)
		}
		cookie := new(fiber.Cookie)
		cookie.Name = "token"
		cookie.Value = t
		cookie.Expires = time.Now().Add(time.Minute * 5)
		cookie.HTTPOnly = true
		cookie.Secure = true
		cookie.SameSite = "lax"
		context.Cookie(cookie)
		context.Status(http.StatusOK)
		return context.JSON(fiber.Map{"token": t})
		// return context.JSON(t)
	}
	context.Status(http.StatusUnauthorized)
	return context.JSON(fiber.Map{"message": "unauthorized"})
}

// Update users
func (r *Respository) UpdateUser(context *fiber.Ctx) error {
	var suser User
	username := context.Params("username")
	r.DB.Where("username=?", username).First(&suser)
	fmt.Println(suser.IsAdmin)
	if !suser.IsAdmin {
		return context.JSON("Not Admin")
	} else {
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
	return context.JSON(user)}
}

// Delete users
func (r *Respository) DelUser(context *fiber.Ctx) error {
	var suser User
	username := context.Params("username")
	r.DB.Where("username=?", username).First(&suser)
	fmt.Println(suser.IsAdmin)
	if !suser.IsAdmin {
		return context.JSON("Not Admin")
	} else {
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
}

// Show users
func (r *Respository) ShowUser(context *fiber.Ctx) error {
	var user []User
	r.DB.Find(&user)
	context.Status(http.StatusOK)
	return context.JSON(user)
}

// Create users
func (r *Respository) CreateUser(context *fiber.Ctx) error {
	var suser User
	username := context.Params("username")
	r.DB.Where("username=?", username).First(&suser)
	fmt.Println(suser.IsAdmin)
	if !suser.IsAdmin {
		return context.JSON("Not Admin")
	} else {

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
}
//logout
func logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name: "token",
		Expires:  time.Now().Add(-(time.Hour * 2)),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "lax",
	})
	return c.JSON("logged out successfully")

}


