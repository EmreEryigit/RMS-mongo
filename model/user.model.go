package model

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id"`
	First_name     *string            `json:"first_name" validate:"required,min=2,max=100"`
	Last_name      *string            `json:"last_name" validate:"required,min=2,max=100"`
	HashedPassword *string            `json:"-"`
	Email          *string            `json:"email" validate:"email,required"`
	Phone          *string            `json:"phone" validate:"required"`
	Avatar         *string            `json:"avatar"`
	Created_at     time.Time          `json:"created_at"`
	Updated_at     time.Time          `json:"updated_at"`
	User_id        *string            `json:"user_id" `
}

type UserPrivate struct {
	User
	Password *string `json:"Password" validate:"required,min=6"`
}

func (u *UserPrivate) HashPassword() {
	hash, err := bcrypt.GenerateFromPassword([]byte(*u.Password), 8)
	if err != nil {
		log.Panic(err)
	}
	str := string(hash)
	u.HashedPassword = &str
}

func (u *User) VerifyPassword(providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*u.HashedPassword), []byte(providedPassword))
	valid := true
	if err != nil {
		valid = false
	}
	return valid
}
