package user

import (
	"context"

	"github.com/jianshao/chrome-exts/CleanTracks/backend/prisma/db"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/utils/prisma"
)

type User struct {
	Email    string
	Password string
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
		Email:    user.Email,
		Password: user.Password,
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
		Email:    user.Email,
		Password: user.Password,
	}, nil
}
