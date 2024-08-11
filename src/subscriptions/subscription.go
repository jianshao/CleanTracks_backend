package subscriptions

import (
	"context"
	"time"

	"github.com/jianshao/chrome-exts/CleanTracks/backend/prisma/db"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/utils/prisma"
)

var (
	SubMap = map[int]int{}
)

type Subscription struct {
	Uid     int
	SubType int
}

type SubInfo struct {
	Uid           int
	Platform      string
	Action        string
	Occurred_time time.Time
	Detail        string
}

func BuildSubInfo(uid int, platform, action, detail string, occurredAt time.Time) *SubInfo {
	return &SubInfo{
		Uid:           uid,
		Platform:      platform,
		Action:        action,
		Detail:        detail,
		Occurred_time: occurredAt,
	}
}

func SaveRecord(uid int, platform, action, detail string, occurredAt time.Time) error {
	client := prisma.GetPrismaClient()
	_, err := client.Subscriptions.CreateOne(
		db.Subscriptions.UID.Set(uid),
		db.Subscriptions.Platform.Set(platform),
		db.Subscriptions.Action.Set(action),
		db.Subscriptions.Details.Set(detail),
		db.Subscriptions.OccurredTime.Set(occurredAt),
	).Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}
