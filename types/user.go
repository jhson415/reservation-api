package types

import (
	"fmt"
	"net/mail"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost  = 10
	minName     = 2
	milPassword = 8
)

type UserPostRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"firstName"`
	LastName          string             `bson:"lastName" json:"lastName"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword,omitempty" json:"-"`
}

func ValidateUserRequest(user UserPostRequest) []string {
	var errorList = []string{}
	if len(user.FirstName) < minName {
		errorList = append(errorList, fmt.Sprintf("First Name should be longer than %d", minName))
	}
	if len(user.LastName) < minName {
		errorList = append(errorList, fmt.Sprintf("Last Name should be longer than %d", minName))
	}
	if len(user.Password) < milPassword {
		errorList = append(errorList, fmt.Sprintf("Password should be longer than %d", milPassword))
	}
	if _, err := mail.ParseAddress(user.Email); err != nil {
		errorList = append(errorList, fmt.Sprintf("invalid email address"))
	}
	return errorList
}

func CreateUserFromRequest(req *UserPostRequest) (*User, error) {
	var user = User{}
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	user.EncryptedPassword = string(encryptedPassword)
	fmt.Printf("%+v", user)

	return &user, nil

}
