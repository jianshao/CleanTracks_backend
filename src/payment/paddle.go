package payment

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/subscriptions"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/user"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/utils/logs"
	"github.com/sirupsen/logrus"
)

type PaddleBasicBillTerm struct {
	EndsAt   time.Time `json:"ends_at"`
	StartsAt time.Time `json:"starts_at"`
}

type PaddleUnitPrice struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}

type PaddleUnitPriceOverrids struct {
	CountryCodes []string        `json:"country_codes"`
	UnitPrice    PaddleUnitPrice `json:"unit_price"`
}

type PaddleQuantity struct {
	Minimum int `json:"minimum"`
	Maximum int `json:"prodmaximumuct_id"`
}

type PaddleEventDataPrice struct {
	Id                 string                  `json:"id"`
	ProductId          string                  `json:"product_id"`
	Description        string                  `json:"description"`
	Type               string                  `json:"type"`
	Name               string                  `json:"name,omitempty"`
	TaxMode            string                  `json:"tax_mode"`
	Status             string                  `json:"status"`
	CreatedAt          time.Time               `json:"created_at"`
	UpdatedAt          time.Time               `json:"updated_at"`
	BillingCycle       PaddleBillingDetailTerm `json:"billing_cycle,omitempty"`
	TrialPeriod        PaddleBillingDetailTerm `json:"trial_period,omitempty"`
	UnitPrice          PaddleUnitPrice         `json:"unit_price"`
	UnitPriceOverrides PaddleUnitPriceOverrids `json:"unit_price_overrides"`
	Quantity           PaddleQuantity          `json:"quantity"`
	ImportMeta         PaddleImportMeta        `json:"import_meta,omitempty"`
	CustomData         map[string]string       `json:"custom_data,omitempty"`
}

type PaddleEventDataProduct struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Type        string            `json:"type"`
	TaxCategory string            `json:"tax_category"`
	ImageUrl    string            `json:"image_url,omitempty"`
	Status      string            `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	ImportMeta  PaddleImportMeta  `json:"import_meta,omitempty"`
	CustomData  map[string]string `json:"custom_data,omitempty"`
}

type PaddleEventDataItem struct {
	Price              PaddleEventDataPrice   `json:"price"`
	Product            PaddleEventDataProduct `json:"product"`
	Status             string                 `json:"status"`
	Quantity           int                    `json:"quantity"`
	Recurring          bool                   `json:"recurring"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	TrialDates         PaddleBasicBillTerm    `json:"trial_dates,omitempty"`
	NextBilledAt       time.Time              `json:"next_billed_at,omitempty"`
	PreviouslyBilledAt time.Time              `json:"previously_billed_at,omitempty"`
}

type PaddleDiscount struct {
	Id       string    `json:"id"`
	EndsAt   time.Time `json:"ends_at,omitempty"`
	StartsAt time.Time `json:"starts_at"`
}

type PaddleBillingDetailTerm struct {
	Interval  string `json:"interval"`
	Frequency int    `json:"frequency"`
}

type PaddleBillingDetail struct {
	EnableCheckout        bool                `json:"enable_checkout"`
	PurchaseOrderNumber   string              `json:"purchase_order_number"`
	AdditionalInformation string              `json:"additional_information,omitempty"`
	PaymentTerms          PaddleBasicBillTerm `json:"payment_terms"`
}

type PaddleBillingCycle struct {
	Interval  string `json:"interval"`
	Frequency int    `json:"frequency"`
}

type PaddleScheduledChange struct {
	Action      string    `json:"action"`
	EffectiveAt time.Time `json:"effective_at"`
	ResumeAt    time.Time `json:"resume_at,omitempty"`
}

type PaddleImportMeta struct {
	ImportedFrom string `json:"imported_from"`
	ExternalId   string `json:"external_id,omitempty"`
}

type PaddleEventData struct {
	Id                   string                `json:"id"`
	TransactionId        string                `json:"transaction_id"`
	Items                []PaddleEventDataItem `json:"items"`
	Status               string                `json:"status"`
	Discount             PaddleDiscount        `json:"discount,omitempty"`
	PausedAt             time.Time             `json:"paused_at,omitempty"`
	AddressId            string                `json:"address_id"`
	CreatedAt            time.Time             `json:"created_at"`
	StartedAt            time.Time             `json:"started_at,omitempty"`
	UpdatedAt            time.Time             `json:"updated_at"`
	BusinessId           string                `json:"business_id,omitempty"`
	CanceledAt           time.Time             `json:"canceled_at,omitempty"`
	CustomData           map[string]string     `json:"custom_data,omitempty"`
	CustomerId           string                `json:"customer_id"`
	ImportMeta           PaddleImportMeta      `json:"import_meta,omitempty"`
	CurrencyCode         string                `json:"currency_code"`
	NextBilledAt         time.Time             `json:"next_billed_at,omitempty"`
	BillingDetails       PaddleBillingDetail   `json:"billing_details,omitempty"`
	CollectionMode       string                `json:"collection_mode"`
	FirstBilledAt        time.Time             `json:"first_billed_at,omitempty"`
	ScheduledChange      PaddleScheduledChange `json:"scheduled_change,omitempty"`
	CurrentBillingPeriod PaddleBasicBillTerm   `json:"current_billing_period,omitempty"`
	BillingCycle         PaddleBasicBillTerm   `json:"billing_cycle"`
}

type PaddleEvent struct {
	EventId        string          `json:"event_id"`
	EventType      string          `json:"event_type"`
	OccurredAt     time.Time       `json:"occurred_at"`
	NotificationId string          `json:"notification_id"`
	Data           PaddleEventData `json:"data"`
}

type ProcFunc func() error

var (
	gEventType = map[string]ProcFunc{
		// activated是试用期结束，trialing是试用期开始
		// "subscription.activated": gPaddle.Create,
		"subscription.canceled": gPaddle.Cancel,
		"subscription.created":  gPaddle.Create,
		// "subscription.imported":  true,
		"subscription.past_due": gPaddle.Expired,
		"subscription.paused":   gPaddle.Pause,
		"subscription.resumed":  gPaddle.Resume,
		// "subscription.trialing":  gPaddle.Create,
		"subscription.updated": gPaddle.Update,
	}
	gEventStatus = map[string]bool{
		"active":   true,
		"canceled": true,
		"past_due": true,
		"paused":   true,
		"trialing": true,
	}
	gMoneyCode = map[string]bool{
		"USD": true,
		"EUR": true,
		"GBP": true,
		"JPY": true,
		"AUD": true,
		"CAD": true,
		"CHF": true,
		"HKD": true,
		"SGD": true,
	}
	gCollectMode = map[string]bool{
		"automatic": true,
		"manual":    true,
	}
	gInterval = map[string]bool{
		"day":   true,
		"week":  true,
		"month": true,
		"year":  true,
	}
	gAction = map[string]bool{
		"cancel": true,
		"pause":  true,
		"resume": true,
	}
	gItemStatus = map[string]bool{
		"active":   true,
		"inactive": true,
		"trialing": true,
	}
	gPriceType = map[string]bool{
		"custom":   true,
		"standard": true,
	}
	gPriceTaxMode = map[string]bool{
		"account_setting": true,
		"external":        true,
		"internal":        true,
	}
	gPriceStatus = map[string]bool{
		"active":   true,
		"archived": true,
	}
	gProductTaxCate = map[string]bool{
		"digital-goods":                 true,
		"ebooks":                        true,
		"implementation-services":       true,
		"professional-services":         true,
		"saas":                          true,
		"software-programming-services": true,
		"standard":                      true,
		"training-services":             true,
		"website-hosting":               true,
	}
)

var (
	gPaddle = &Paddle{}
)

type Paddle struct {
	paddle *PaddleEvent
	uid    int
}

func (p *Paddle) GetPaymentURL(paymentID string) string {
	return "https://checkout.paddle.com/checkout/oneoff/package/" + paymentID
}

func checkData(event *PaddleEvent) error {
	return nil
}

func paddleCheckAuth(event *PaddleEvent) error {
	return nil
}

func paddleFindUser(event *PaddleEvent) (int, error) {
	email, ok := event.Data.CustomData["email"]
	if !ok {
		return 0, errors.New("")
	}
	userInfo, err := user.FindUser(email)
	if err != nil {
		return 0, err
	}
	return userInfo.Id, nil
}

func (p *Paddle) Prepare(data any) error {
	event, ok := data.(*PaddleEvent)
	if !ok {
		return errors.New("")
	}
	// 检查数据
	if err := checkData(event); err != nil {
		return errors.New("")
	}
	// 检查权限
	if err := paddleCheckAuth(event); err != nil {
		return errors.New("")
	}
	// 查找用户
	uid, err := paddleFindUser(event)
	if err != nil {
		return errors.New("")
	}

	p.uid = uid
	p.paddle = event
	return nil
}

type PaddleSubDetail struct {
	PriceId   string `json:"price_id"`
	ProductId string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func buildPaddleSubDetail(item PaddleEventDataItem) ([]byte, error) {
	detail := PaddleSubDetail{
		PriceId:   item.Price.Id,
		ProductId: item.Price.ProductId,
		Quantity:  item.Quantity,
	}
	return json.Marshal(detail)
}

var (
	gSubTypeMap = map[string]int{
		"pri_01j4nt2jb3brrha2pxa8ma9c8f": 1,
		"pri_01j4nt5zw06qtpctxtk0vjp6xk": 2,
		"pri_01j4r9hzer7j9bws375kj1557w": 2,
		"pri_01j4r9h4em2ak1tkv789evmcbb": 1,
	}
)

func (p *Paddle) Create() error {
	for _, item := range p.paddle.Data.Items {
		details, err := buildPaddleSubDetail(item)
		if err != nil {
			return nil
		}

		// 保存订阅历史记录
		err = subscriptions.SaveRecord(p.uid, "paddle", p.paddle.EventType, string(details), p.paddle.OccurredAt)
		if err != nil {
			return err
		}

		status, ok := gSubTypeMap[item.Price.Id]
		if !ok {
			logs.WriteLog(logrus.ErrorLevel, nil, "")
		}
		// 更新用户当前订阅状态
		err = user.Subscribe(p.uid, status)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Paddle) Update() error {
	return nil
}

func (p *Paddle) Cancel() error {
	return nil
}

func (p *Paddle) Pause() error {
	return nil
}

func (p *Paddle) Resume() error {
	return nil
}

func (p *Paddle) Expired() error {
	return nil
}

func PaddleWebHookHandle(c *gin.Context) {
	var request PaddleEvent
	if err := c.BindJSON(&request); err != nil {
		logs.WriteLog(logrus.ErrorLevel, nil, err.Error())
		return
	}

	var lock sync.Mutex
	lock.Lock()
	if err := gPaddle.Prepare(request); err != nil {
		logs.WriteLog(logrus.ErrorLevel, nil, err.Error())
	}

	if handler, ok := gEventType[request.EventType]; ok {
		if err := handler(); err != nil {
			logs.WriteLog(logrus.ErrorLevel, nil, err.Error())
		}
	}
	lock.Unlock()

	// 最后返回一个200的responsse
	c.JSON(http.StatusOK, nil)
}
