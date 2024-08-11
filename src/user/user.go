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

func Subscribe(uid, status int) error {
	client := prisma.GetPrismaClient()
	_, err := client.User.FindUnique(db.User.ID.Equals(uid)).Update(db.User.Status.Set(status)).Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func GetCurrSubscribe(uid int) (int, error) {
	client := prisma.GetPrismaClient()
	user, err := client.User.FindUnique(db.User.ID.Equals(uid)).Exec(context.Background())
	if err != nil {
		return 0, nil
	}
	if user.Status == 0 {
		// 用户注册后7天内是试用期
		if !time.Now().After(user.RegisterTime.AddDate(0, 0, 7)) {
			return 1, nil
		}
	}
	return user.Status, nil
}
