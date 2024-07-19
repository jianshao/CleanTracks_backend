package payment

import (
	"context"
	"fmt"
	"reflect"

	"github.com/NdoleStudio/lemonsqueezy-go"
	"github.com/gin-gonic/gin"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/prisma/db"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/utils"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/utils/logs"
	"github.com/sirupsen/logrus"
)

var client *lemonsqueezy.Client

func init() {
	// client := lemonsqueezy.New(lemonsqueezy.WithAPIKey(""))
}

func checkAuth() (bool, error) {
	return true, nil
}

type handler func(sub *lemonsqueezy.WebhookRequestSubscription)

var (
	funcMap = map[string]handler{
		lemonsqueezy.WebhookEventSubscriptionCreated:   subCreated,
		lemonsqueezy.WebhookEventSubscriptionUpdated:   subUpdated,
		lemonsqueezy.WebhookEventSubscriptionCancelled: subCancell,
		lemonsqueezy.WebhookEventSubscriptionResumed:   subResumed,
		lemonsqueezy.WebhookEventSubscriptionPaused:    subPaused,
		lemonsqueezy.WebhookEventSubscriptionUnpaused:  subUnPaused,
		lemonsqueezy.WebhookEventSubscriptionExpired:   subExpired,
	}
)

func subCreated(sub *lemonsqueezy.WebhookRequestSubscription) {
	if reflect.TypeOf(sub.Meta.CustomData["uid"]).Kind() != reflect.Int {
		logs.WriteLog(logrus.ErrorLevel, nil, "no uid")
	}

	client := utils.GetPrismaClient()
	_, err := client.Subscriptions.CreateOne(
		db.Subscriptions.UID.Set(sub.Meta.CustomData["uid"].(int)),
		db.Subscriptions.StoreID.Set(sub.Data.Attributes.StoreID),
		db.Subscriptions.ProductID.Set(sub.Data.Attributes.ProductID),
		db.Subscriptions.VariantID.Set(sub.Data.Attributes.VariantID),
		db.Subscriptions.SubscriptionID.Set(sub.Data.Attributes.FirstSubscriptionItem.SubscriptionID),
	).Exec(context.Background())
	if err != nil {
		logs.WriteLog(logrus.ErrorLevel, nil, err.Error())
	}
	return
}

func subUpdated(sub *lemonsqueezy.WebhookRequestSubscription) {
}

func subPaused(sub *lemonsqueezy.WebhookRequestSubscription) {
	client := utils.GetPrismaClient()
	client.Subscriptions.UpsertOne(
		db.Subscriptions.SubscriptionID.Equals(sub.Data.Attributes.FirstSubscriptionItem.SubscriptionID),
	).Update(db.Subscriptions.Status.Set(db.SubStatusPaused)).Exec(context.Background())
}

func subUnPaused(sub *lemonsqueezy.WebhookRequestSubscription) {
	client := utils.GetPrismaClient()
	client.Subscriptions.UpsertOne(
		db.Subscriptions.SubscriptionID.Equals(sub.Data.Attributes.FirstSubscriptionItem.SubscriptionID),
	).Update(db.Subscriptions.Status.Set(db.SubStatusPaid)).Exec(context.Background())
}

func subExpired(sub *lemonsqueezy.WebhookRequestSubscription) {
	client := utils.GetPrismaClient()
	client.Subscriptions.UpsertOne(
		db.Subscriptions.SubscriptionID.Equals(sub.Data.Attributes.FirstSubscriptionItem.SubscriptionID),
	).Update(db.Subscriptions.Status.Set(db.SubStatusExpired)).Exec(context.Background())
}

func subCancell(sub *lemonsqueezy.WebhookRequestSubscription) {

}

func subResumed(sub *lemonsqueezy.WebhookRequestSubscription) {

}

func subPaid(sub *lemonsqueezy.WebhookRequestSubscriptionInvoice) {
	client := utils.GetPrismaClient()
	sid := sub.Data.Attributes.SubscriptionID
	_, err := client.Subscriptions.FindUnique(db.Subscriptions.SubscriptionID.Equals(sid)).Update(
		db.Subscriptions.Status.Set(db.SubStatusPaid),
	).Exec(context.Background())
	if err != nil {
		logs.WriteLog(logrus.ErrorLevel, nil, err.Error())
	}
}

func WebhookHandler(c *gin.Context) {

	// 1. Authenticate the webhook request from Lemon Squeezy using the `X-Signature` header
	// sign := c.Request.Header.Get("X-Signature")

	// 2. Process the payload if the request is authenticated
	eventName := c.Request.Header.Get("X-Event-Name")
	logs.WriteLog(logrus.WarnLevel, nil, fmt.Sprintf("webhook revieve event: %s", eventName))

	switch eventName {
	case lemonsqueezy.WebhookEventSubscriptionCreated:
	case lemonsqueezy.WebhookEventSubscriptionUpdated:
	case lemonsqueezy.WebhookEventSubscriptionCancelled:
	case lemonsqueezy.WebhookEventSubscriptionResumed:
	case lemonsqueezy.WebhookEventSubscriptionExpired:
	case lemonsqueezy.WebhookEventSubscriptionPaused:
	case lemonsqueezy.WebhookEventSubscriptionUnpaused:
		var request lemonsqueezy.WebhookRequestSubscription
		if err := c.BindJSON(&request); err != nil {
			logs.WriteLog(logrus.ErrorLevel, nil, err.Error())
		}
		funcMap[eventName](&request)
		//
	case lemonsqueezy.WebhookEventSubscriptionPaymentSuccess:
		var request lemonsqueezy.WebhookRequestSubscriptionInvoice
		if err := c.BindJSON(&request); err != nil {
			logs.WriteLog(logrus.ErrorLevel, nil, err.Error())
		}
		subPaid(&request)
	default:
		logs.WriteLog(logrus.ErrorLevel, nil, fmt.Sprintf("invalid event [%s] received with request", eventName))
	}
}
