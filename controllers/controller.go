package controllers

import (
	"context"
	"fmt"
	"net/http"
	"purnur/pjt/structs"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var SignUpTemp structs.SignUpStruct
	c.ShouldBindJSON(&SignUpTemp)

	if SignUpTemp.Name == "" || SignUpTemp.Surname == "" || SignUpTemp.Login == "" || SignUpTemp.Password =="" {
		c.JSON(404, "Error...empty field")
	} else {
		client, ctx := DBConnection()

		DBConnect := client.Database("SocialMedia").Collection("Users")

		id := primitive.NewObjectID().Hex()
		Hashed, _ := HashPassword(SignUpTemp.Password)

		DBConnect.InsertOne(ctx, bson.M{
			"_id":      id,
			"name":     SignUpTemp.Name,
			"surname":  SignUpTemp.Surname,
			"login":    SignUpTemp.Login,
			"password": Hashed,
		})
	}
}

func Login(c *gin.Context) {
	var LoginTemp structs.SignUpStruct
	c.ShouldBindJSON(&LoginTemp)

	if LoginTemp.Login == "" || LoginTemp.Password == "" {
		c.JSON(404, "Error...empty field")
	} else {
		client, ctx := DBConnection()

		DBConnect := client.Database("SocialMedia").Collection("Users")

		result := DBConnect.FindOne(ctx, bson.M{
			"login": LoginTemp.Login,
		})

		var userdata structs.SignUpStruct
		result.Decode(&userdata)
		isValidPass := CompareHashPasswords(userdata.Password, LoginTemp.Password)
		fmt.Println(isValidPass)

		if isValidPass {
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "RG-Cookie",
				Value:    userdata.Id,
				Expires:  time.Now().Add(60 * time.Second),
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})
			c.JSON(200, "Success")
		} else {
			c.JSON(404, "Wrong login or password")
		}
	}
}





func DBConnection() (*mongo.Client, context.Context) {
	url := options.Client().ApplyURI("mongodb://192.168.43.246:27017")
	NewCtx, _ := context.WithTimeout(context.Background(), 2*time.Second)

	Client, err := mongo.Connect(NewCtx, url)
	if err != nil {
		fmt.Printf("errors: %v\n", err)
	}
	return Client, NewCtx
}

func HashPassword(password string) (string, error) {
	var passwordBytes = []byte(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)

	return string(hashedPasswordBytes), err
}

func CompareHashPasswords(HashedPasswordFromDB, PasswordToCampare string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(HashedPasswordFromDB), []byte(PasswordToCampare))
	return err == nil
}
