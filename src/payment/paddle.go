package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	Id           string                   `json:"id"`
	ProductId    string                   `json:"product_id"`
	Description  string                   `json:"description"`
	Type         string                   `json:"type"`
	Name         *string                  `json:"name" binding:"omitempty"`
	TaxMode      string                   `json:"tax_mode"`
	Status       string                   `json:"status"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
	BillingCycle *PaddleBillingDetailTerm `json:"billing_cycle" binding:"omitempty"`
	TrialPeriod  *PaddleBillingDetailTerm `json:"trial_period" binding:"omitempty"`
	UnitPrice    PaddleUnitPrice          `json:"unit_price"`
	// UnitPriceOverrides PaddleUnitPriceOverrids  `json:"unit_price_overrides"`
	Quantity   PaddleQuantity     `json:"quantity"`
	ImportMeta *PaddleImportMeta  `json:"import_meta" binding:"omitempty"`
	CustomData *map[string]string `json:"custom_data" binding:"omitempty"`
}

type PaddleEventDataProduct struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description *string            `json:"description" binding:"omitempty"`
	Type        string             `json:"type"`
	TaxCategory string             `json:"tax_category"`
	ImageUrl    *string            `json:"image_url" binding:"omitempty"`
	Status      string             `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	ImportMeta  *PaddleImportMeta  `json:"import_meta" binding:"omitempty"`
	CustomData  *map[string]string `json:"custom_data" binding:"omitempty"`
}

type PaddleEventDataItem struct {
	Price              PaddleEventDataPrice   `json:"price"`
	Product            PaddleEventDataProduct `json:"product"`
	Status             string                 `json:"status"`
	Quantity           int                    `json:"quantity"`
	Recurring          bool                   `json:"recurring"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	TrialDates         *PaddleBasicBillTerm   `json:"trial_dates" binding:"omitempty"`
	NextBilledAt       *time.Time             `json:"next_billed_at" binding:"omitempty"`
	PreviouslyBilledAt *time.Time             `json:"previously_billed_at" binding:"omitempty"`
}

type PaddleDiscount struct {
	Id       string     `json:"id"`
	EndsAt   *time.Time `json:"ends_at" binding:"omitempty"`
	StartsAt time.Time  `json:"starts_at"`
}

type PaddleBillingDetailTerm struct {
	Interval  string `json:"interval"`
	Frequency int    `json:"frequency"`
}

type PaddleBillingDetail struct {
	EnableCheckout        bool                `json:"enable_checkout"`
	PurchaseOrderNumber   string              `json:"purchase_order_number"`
	AdditionalInformation *string             `json:"additional_information" binding:"omitempty"`
	PaymentTerms          PaddleBasicBillTerm `json:"payment_terms"`
}

type PaddleBillingCycle struct {
	Interval  string `json:"interval"`
	Frequency int    `json:"frequency"`
}

type PaddleScheduledChange struct {
	Action      string     `json:"action"`
	EffectiveAt time.Time  `json:"effective_at"`
	ResumeAt    *time.Time `json:"resume_at" binding:"omitempty"`
}

type PaddleImportMeta struct {
	ImportedFrom string  `json:"imported_from"`
	ExternalId   *string `json:"external_id" binding:"omitempty"`
}

type PaddleEventData struct {
	Id                   string                 `json:"id"`
	TransactionId        string                 `json:"transaction_id"`
	Items                []PaddleEventDataItem  `json:"items"`
	Status               string                 `json:"status"`
	Discount             *PaddleDiscount        `json:"discount" binding:"omitempty"`
	PausedAt             *time.Time             `json:"paused_at" binding:"omitempty"`
	AddressId            string                 `json:"address_id"`
	CreatedAt            time.Time              `json:"created_at"`
	StartedAt            *time.Time             `json:"started_at" binding:"omitempty"`
	UpdatedAt            time.Time              `json:"updated_at"`
	BusinessId           *string                `json:"business_id" binding:"omitempty"`
	CanceledAt           *time.Time             `json:"canceled_at" binding:"omitempty"`
	CustomData           *map[string]string     `json:"custom_data" binding:"omitempty"`
	CustomerId           string                 `json:"customer_id"`
	ImportMeta           *PaddleImportMeta      `json:"import_meta" binding:"omitempty"`
	CurrencyCode         string                 `json:"currency_code"`
	NextBilledAt         *time.Time             `json:"next_billed_at" binding:"omitempty"`
	BillingDetails       *PaddleBillingDetail   `json:"billing_details" binding:"omitempty"`
	CollectionMode       string                 `json:"collection_mode"`
	FirstBilledAt        *time.Time             `json:"first_billed_at" binding:"omitempty"`
	ScheduledChange      *PaddleScheduledChange `json:"scheduled_change" binding:"omitempty"`
	CurrentBillingPeriod *PaddleBasicBillTerm   `json:"current_billing_period" binding:"omitempty"`
	BillingCycle         PaddleBasicBillTerm    `json:"billing_cycle"`
}

type PaddleEvent struct {
	EventId        string          `json:"event_id"`
	EventType      string          `json:"event_type"`
	OccurredAt     time.Time       `json:"occurred_at"`
	NotificationId string          `json:"notification_id"`
	Data           PaddleEventData `json:"data"`
}

type CustomerRespData struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CustomerResp struct {
	Data CustomerRespData `json:"data"`
}

type ProcFunc func() error

var (
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
	gSubTypeMap = map[string]int{
		"pri_01j4nt2jb3brrha2pxa8ma9c8f": 1,
		"pri_01j4nt5zw06qtpctxtk0vjp6xk": 2,
		"pri_01j4r9hzer7j9bws375kj1557w": 2,
		"pri_01j4r9h4em2ak1tkv789evmcbb": 1,
	}
	gEventType = map[string]ProcFunc{
		// activated是试用期结束，trialing是试用期开始
		"subscription.created": gPaddle.Create,
		// "subscription.trialing":  gPaddle.Create,
		"subscription.activated": gPaddle.Active,
		"subscription.canceled":  gPaddle.Cancel,
		// "subscription.imported":  true,
		"subscription.past_due": gPaddle.Expired,
		"subscription.paused":   gPaddle.Pause,
		"subscription.resumed":  gPaddle.Resume,
		"subscription.updated":  gPaddle.Update,
	}
)

var (
	gPaddle = &Paddle{}
)

type Paddle struct {
	paddle *PaddleEvent
	uid    int
}

func checkData(event *PaddleEvent) error {
	return nil
}

func paddleCheckAuth(event *PaddleEvent) error {
	return nil
}

func getCustomerEmail(customerId string) (string, error) {
	url := "https://sandbox-api.paddle.com/customers/" + customerId
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to create request: %s", err.Error())
	}
	// 设置请求头
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer ")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("get customer info from api failed: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read response body: %v", err.Error())
	}

	// 解析 JSON 响应
	var custom CustomerResp
	if err := json.Unmarshal(body, &custom); err != nil {
		return "", fmt.Errorf("Failed to unmarshal JSON: %v", err.Error())
	}

	return custom.Data.Email, nil
}

func paddleFindUser(event *PaddleEvent) (int, error) {
	email, err := getCustomerEmail(event.Data.CustomerId)
	if err != nil {
		return 0, err
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
		return errors.New("prepare: invalid struct type")
	}
	// 检查数据
	if err := checkData(event); err != nil {
		return err
	}
	// 检查权限
	if err := paddleCheckAuth(event); err != nil {
		return err
	}
	// 查找用户
	uid, err := paddleFindUser(event)
	if err != nil {
		return err
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

func buildPaddleSubDetail(item *PaddleEventDataItem) ([]byte, error) {
	detail := PaddleSubDetail{
		PriceId:   item.Price.Id,
		ProductId: item.Price.ProductId,
		Quantity:  item.Quantity,
	}
	return json.Marshal(detail)
}

func save2DB(uid int, item *PaddleEventDataItem, paddle *PaddleEvent) error {
	details, err := buildPaddleSubDetail(item)
	if err != nil {
		return err
	}

	// 保存订阅历史记录
	err = subscriptions.SaveRecord(uid, "paddle", paddle.EventType, string(details), paddle.OccurredAt)
	if err != nil {
		return err
	}
	return nil
}

// 订阅创建时会先经历试用期（不论是否设置了试用期）,当active时才是正式订阅（已付款）
func (p *Paddle) Create() error {
	for _, item := range p.paddle.Data.Items {
		if err := save2DB(p.uid, &item, p.paddle); err != nil {
			return fmt.Errorf("Save to db failed: %s", err.Error())
		}

		// if status, ok := gSubTypeMap[item.Price.Id]; ok {
		// 	// 更新用户当前订阅类型
		// 	if err := user.Subscribe(p.uid, status); err != nil {
		// 		return err
		// 	}
		// } else {
		// 	return errors.New(fmt.Sprintf("invalid price id: %s", item.Price.Id))
		// }
	}

	return nil
}

func (p *Paddle) Active() error {
	for _, item := range p.paddle.Data.Items {
		if err := save2DB(p.uid, &item, p.paddle); err != nil {
			return fmt.Errorf("Save to db failed: %s", err.Error())
		}

		if status, ok := gSubTypeMap[item.Price.Id]; ok {
			// 更新用户当前订阅类型
			if err := user.Subscribe(p.uid, status); err != nil {
				return err
			}
		} else {
			return errors.New(fmt.Sprintf("invalid price id: %s", item.Price.Id))
		}
	}

	return nil
}

// 更新订阅，订阅升级或降级
func (p *Paddle) Update() error {
	for _, item := range p.paddle.Data.Items {
		if err := save2DB(p.uid, &item, p.paddle); err != nil {
			return fmt.Errorf("Save to db failed: %s", err.Error())
		}

		// 更新用户当前订阅类型
		if status, ok := gSubTypeMap[item.Price.Id]; ok {
			if err := user.Subscribe(p.uid, status); err != nil {
				return err
			}
		} else {
			return errors.New(fmt.Sprintf("invalid price id: %s", item.Price.Id))
		}
	}
	return nil
}

func (p *Paddle) Cancel() error {
	for _, item := range p.paddle.Data.Items {
		if err := save2DB(p.uid, &item, p.paddle); err != nil {
			return fmt.Errorf("Save to db failed: %s", err.Error())
		}

		// 更新用户当前订阅类型
		if err := user.Subscribe(p.uid, 0); err != nil {
			return err
		}
	}
	return nil
}

func (p *Paddle) Pause() error {
	for _, item := range p.paddle.Data.Items {
		if err := save2DB(p.uid, &item, p.paddle); err != nil {
			return fmt.Errorf("Save to db failed: %s", err.Error())
		}

		// 更新用户当前订阅类型
		if err := user.Subscribe(p.uid, 0); err != nil {
			return err
		}
	}
	return nil
}

func (p *Paddle) Resume() error {
	for _, item := range p.paddle.Data.Items {
		if err := save2DB(p.uid, &item, p.paddle); err != nil {
			return fmt.Errorf("Save to db failed: %s", err.Error())
		}

		// 更新用户当前订阅类型
		if status, ok := gSubTypeMap[item.Price.Id]; ok {
			if err := user.Subscribe(p.uid, status); err != nil {
				return err
			}
		} else {
			return errors.New(fmt.Sprintf("invalid price id: %s", item.Price.Id))
		}
	}
	return nil
}

func (p *Paddle) Expired() error {
	for _, item := range p.paddle.Data.Items {
		if err := save2DB(p.uid, &item, p.paddle); err != nil {
			return fmt.Errorf("Save to db failed: %s", err.Error())
		}

		// 更新用户当前订阅类型
		if err := user.Subscribe(p.uid, 0); err != nil {
			return err
		}
	}
	return nil
}

func PaddleWebHookHandle(c *gin.Context) {
	var request PaddleEvent
	if err := c.ShouldBindJSON(&request); err == nil {
		var lock sync.Mutex
		// 在出错的情况下保存当前请求信息
		data, _ := json.Marshal(request)

		lock.Lock()
		defer lock.Unlock()
		if err := gPaddle.Prepare(&request); err == nil {
			if handler, ok := gEventType[request.EventType]; ok {
				if err := handler(); err != nil {
					logs.WriteLog(logrus.ErrorLevel, nil, fmt.Sprintf("payment request: %s", data))
					logs.WriteLog(logrus.ErrorLevel, nil, fmt.Sprintf("handlle %s failed, %s", request.EventType, err.Error()))
				}
			} else {
				logs.WriteLog(logrus.ErrorLevel, nil, fmt.Sprintf("payment request: %s", data))
				logs.WriteLog(logrus.ErrorLevel, nil, "invalid event type: "+request.EventType)
			}
		} else {
			logs.WriteLog(logrus.ErrorLevel, nil, fmt.Sprintf("payment request: %s", data))
			logs.WriteLog(logrus.ErrorLevel, nil, fmt.Sprintf("prepare failed, %s", err.Error()))
		}

	} else {
		body, _ := ioutil.ReadAll(c.Request.Body)
		logs.WriteLog(logrus.ErrorLevel, nil, fmt.Sprintf("payment request: %s", body))
		logs.WriteLog(logrus.ErrorLevel, nil, fmt.Sprintf("BindJson failed, %s", err.Error()))
	}

	// 不论是否成功处理最后都返回一个200的responsse
	c.JSON(http.StatusOK, nil)
}
