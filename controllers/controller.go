package controllers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"purnur/pjt/helpers"
	"purnur/pjt/structs"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Signup(c *gin.Context) {
	var SignUpTemp structs.SignUpStruct
	c.ShouldBindJSON(&SignUpTemp)

	if SignUpTemp.Name == "" || SignUpTemp.Surname == "" || SignUpTemp.Login == "" || SignUpTemp.Password == "" {
		c.JSON(404, "Error...empty field")
	} else {
		client, ctx := helpers.DBConnection()

		DBConnect := client.Database("SocialMedia").Collection("Users")

		Hashed, _ := helpers.HashPassword(SignUpTemp.Password)

		id := primitive.NewObjectID().Hex()

		DBConnect.InsertOne(ctx, bson.M{
			"_id":      id,
			"name":     SignUpTemp.Name,
			"surname":  SignUpTemp.Surname,
			"login":    SignUpTemp.Login,
			"password": Hashed,
		})

		c.JSON(200, "Success")
	}

}

func Login(c *gin.Context) {
	var LoginTemp structs.SignUpStruct
	c.ShouldBindJSON(&LoginTemp)

	if LoginTemp.Login == "" || LoginTemp.Password == "" {
		c.JSON(404, "Error...empty field")
	} else {
		client, ctx := helpers.DBConnection()

		DBConnect := client.Database("SocialMedia").Collection("Users")

		result := DBConnect.FindOne(ctx, bson.M{
			"login": LoginTemp.Login,
		})

		var userdata structs.SignUpStruct
		result.Decode(&userdata)
		isValidPass := helpers.CompareHashPasswords(userdata.Password, LoginTemp.Password)
		fmt.Println(isValidPass)
		fmt.Printf("userdata: %v\n", userdata)

		if isValidPass {
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "media",
				Value:    userdata.Id,
				Expires:  time.Now().Add(60 * time.Second),
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
				Path:     "/",
				Domain:   "",
			})
			c.JSON(200, "Success")
		} else {
			c.JSON(404, "Wrong login or password")
		}
	}
}

func CreateFunc(c *gin.Context) {
	CookieData, err := c.Request.Cookie("media")
	fmt.Printf("CookieData: %v\n", CookieData)

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	if CookieData.Value != "" {
		var CreateTemp structs.PostStruct
		c.ShouldBindJSON(&CreateTemp)

		bytedata, Error := base64.StdEncoding.DecodeString(CreateTemp.Image)

		if Error != nil {
			fmt.Printf("Error: %v\n", Error)
		}

		rand := rand.Intn(10000)
		name := fmt.Sprintf("./static/%v.jpg", rand)
		error := ioutil.WriteFile(name, bytedata, 0644)

		if error != nil {
			fmt.Printf("error: %v\n", error)
		}
		//? ======================
		client, ctx := helpers.DBConnection()

		DBConnect := client.Database("SocialMedia").Collection("AllPosts")

		id := primitive.NewObjectID().Hex()

		img := fmt.Sprintf("%v.jpg", rand)

		_, Error2 := DBConnect.InsertOne(ctx, bson.M{
			"_id":         id,
			"owner_id":    CookieData.Value,
			"title":       CreateTemp.Title,
			"image":       img,
			"description": CreateTemp.Description,
		})

		if Error2 != nil {
			fmt.Printf("Error: %v\n", Error)
		}

	}
}

func AllPosts(c *gin.Context) {
	var AllPosts = []structs.PostStruct{}

	CookieData, err := c.Request.Cookie("media")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	if CookieData.Value != "" {
		client, ctx := helpers.DBConnection()

		connect := client.Database("SocialMedia").Collection("AllPosts")

		Result, _ := connect.Find(ctx, bson.M{})

		for Result.Next(ctx) {
			var shablon structs.PostStruct
			Result.Decode(&shablon)

			AllPosts = append(AllPosts, shablon)
		}

		c.JSON(200, AllPosts)
	}
}

func MyPost(c *gin.Context) {
	var MyPosts = []structs.PostStruct{}

	CookieData, err := c.Request.Cookie("media")
	fmt.Printf("CookieData.Value: %v\n", CookieData.Value)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	if CookieData.Value != "" {
		client, ctx := helpers.DBConnection()

		connect := client.Database("SocialMedia").Collection("AllPosts")

		Result, _ := connect.Find(ctx, bson.M{
			"owner_id": CookieData.Value,
		})
		fmt.Printf("Result: %v\n", Result)
		for Result.Next(ctx) {
			var shablon structs.PostStruct
			Result.Decode(&shablon)

			MyPosts = append(MyPosts, shablon)
		}

		c.JSON(200, MyPosts)
	}
}

func Like(c *gin.Context) {
	var LikeTemp structs.PostStruct
	c.ShouldBindJSON(&LikeTemp)

	CookieData, CookieError := c.Request.Cookie("RGCookie")

	if CookieError != nil {
		fmt.Printf("CookieError: %v\n", CookieError)
	}

	if CookieData.Value != "" {
		client, ctx := helpers.DBConnection()

		connect := client.Database("SocialMedia").Collection("AllPost")

		SingleResult := connect.FindOne(ctx, bson.M{
			"owner_id": LikeTemp.Owner_id,
		})

		var likedata structs.PostStruct
		SingleResult.Decode(&likedata)

		connection2 := client.Database("SocialMedia").Collection("Users")

		SingleResult2 := connection2.FindOne(ctx, bson.M{
			"_id": LikeTemp.Id,
		})

		var likedata2 structs.PostStruct
		SingleResult2.Decode(&likedata2)

		if likedata.Id == "" || likedata.Owner_id == "" || likedata2.Id == "" || likedata2.Owner_id == "" {
			c.JSON(404, "Error...empty field")
		} else {
			connect3 := client.Database("SocialMedia").Collection("AllPost")
			connect3.UpdateOne(ctx, bson.M{
				"like": 0,
			},
				bson.D{
					{
						"$inc", bson.D{
							{"like", 1},
						},
					},
				})
		}
	}
}
