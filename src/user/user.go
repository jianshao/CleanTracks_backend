package user

import (
	"context"
	"time"

	"github.com/jianshao/chrome-exts/CleanTracks/backend/prisma/db"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/utils/prisma"
)

type User struct {
	Id           int
	Email        string
	Password     string
	RegisterTime time.Time
}

func CreateUser(email, password string) (*User, error) {
	client := prisma.GetPrismaClient()
	user, err := client.User.CreateOne(
		db.User.Email.Set(email),
		db.User.Password.Set(password),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &User{
		Id:           user.ID,
		Email:        user.Email,
		Password:     user.Password,
		RegisterTime: user.RegisterTime,
	}, nil
}

func FindUser(email string) (*User, error) {
	client := prisma.GetPrismaClient()
	user, err := client.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(context.Background())
	if err != nil {
		return nil, err
	}
	return &User{
		Id:           user.ID,
		Email:        user.Email,
		Password:     user.Password,
		RegisterTime: user.RegisterTime,
	}, nil
}

func GetUserById(uid int) (*User, error) {
	client := prisma.GetPrismaClient()
	user, err := client.User.FindUnique(db.User.ID.Equals(uid)).Exec(context.Background())
	if err == nil {
		return &User{
			Id:           user.ID,
			Email:        user.Email,
			Password:     user.Password,
			RegisterTime: user.RegisterTime,
		}, nil
	}
	return nil, err
}
