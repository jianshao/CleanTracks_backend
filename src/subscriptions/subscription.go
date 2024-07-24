package subscriptions

import (
	"context"
	"time"

	"github.com/jianshao/chrome-exts/CleanTracks/backend/prisma/db"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/user"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/utils/prisma"
)

var (
	SubMap = map[int]int{}
)

type Subscription struct {
	Uid     int
	SubType int
}

func GetCurrSubscribe(uid int) (*Subscription, error) {
	client := prisma.GetPrismaClient()
	user, err := user.GetUserById(uid)
	if err != nil {
		return nil, err
	}

	subscribe := &Subscription{
		Uid:     user.Id,
		SubType: 0,
	}
	// 用户注册后7天内是试用期
	if !time.Now().After(user.RegisterTime.AddDate(0, 0, 7)) {
		subscribe.SubType = 1
	}

	// 如果在7天内有订阅，则用订阅的信息
	sub, err := client.Subscriptions.FindFirst(
		db.Subscriptions.UID.Equals(uid),
		db.Subscriptions.Status.Equals(db.SubStatusPaid),
	).Exec(context.Background())
	if err == db.ErrNotFound {
		return subscribe, nil
	} else if err != nil {
		return nil, err
	} else {
		if data, ok := SubMap[sub.VariantID]; ok {
			subscribe.SubType = data
		}
	}
	return subscribe, nil
}
