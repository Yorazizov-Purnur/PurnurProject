package helpers

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func DBConnection() (*mongo.Client, context.Context) {
	url := options.Client().ApplyURI("mongodb://192.168.43.246:27018")
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