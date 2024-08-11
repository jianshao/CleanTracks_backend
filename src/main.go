package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/payment"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/user"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/utils"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/utils/logs"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func register(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusOK, utils.BuildApiResponse(1, "Invalid request", nil))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusOK, utils.BuildApiResponse(2, "Error creating user", nil))
		return
	}

	_, err = user.CreateUser(creds.Email, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusOK, utils.BuildApiResponse(3, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.BuildApiResponse(0, "", nil))
}

type LoginResp struct {
	Token        string `json:"token"`
	Uid          int    `json:"uid"`
	Subscription int    `json:"subscription"`
	Email        string `json:"email"`
}

func generateToken(email string) (string, error) {

	expirationTime := time.Now().Add(30 * 24 * time.Hour)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func login(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusOK, utils.BuildApiResponse(1, "Invalid request", nil))
		return
	}

	existedUser, err := user.FindUser(creds.Email)
	if err != nil {
		logs.WriteLog(logrus.ErrorLevel, nil, "find user: "+err.Error())
		c.JSON(http.StatusOK, utils.BuildApiResponse(2, "Invalid email or password", nil))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(creds.Password))
	if err != nil {
		logs.WriteLog(logrus.ErrorLevel, nil, "CompareHashAndPassword: "+err.Error())
		c.JSON(http.StatusOK, utils.BuildApiResponse(2, "Invalid email or password", nil))
		return
	}

	tokenString, err := generateToken(creds.Email)
	if err != nil {
		c.JSON(http.StatusOK, utils.BuildApiResponse(2, "Error generating token", nil))
		return
	}

	resp := LoginResp{Token: tokenString, Uid: existedUser.Id, Email: existedUser.Email}
	sub, err := user.GetCurrSubscribe(existedUser.Id)
	if err == nil {
		resp.Subscription = sub
	}

	c.JSON(http.StatusOK, utils.BuildApiResponse(0, "", resp))
}

func authenticate(c *gin.Context) {
	tokenStr := c.GetHeader("token")
	if tokenStr == "" {
		c.JSON(http.StatusOK, utils.BuildApiResponse(1, "Not authorized", nil))
		c.Abort()
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusOK, utils.BuildApiResponse(1, "Not authorized", nil))
		c.Abort()
		return
	}

	c.Set("email", claims.Email)
	c.Next()
}

func checkLogin(c *gin.Context) {
	email := c.GetString("email")
	existedUser, err := user.FindUser(email)
	if err != nil {
		c.JSON(http.StatusOK, utils.BuildApiResponse(1, "user not existed", nil))
		return
	}

	resp := LoginResp{
		Uid:   existedUser.Id,
		Token: c.GetHeader("token"),
		Email: existedUser.Email,
	}
	sub, err := user.GetCurrSubscribe(existedUser.Id)
	if err == nil {
		resp.Subscription = sub
	}
	c.JSON(http.StatusOK, utils.BuildApiResponse(0, "", resp))
}

func main() {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		// 设置 CORS 响应头
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, token")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			// 处理预检请求
			c.AbortWithStatusJSON(http.StatusOK, gin.H{})
			return
		}

		c.Next()
	})

	path := "./logs/cleantracks.log"
	if !utils.Init(path, logrus.DebugLevel) {

		return
	}

	router.POST("/cleantracks/api/webhook", payment.PaddleWebHookHandle)
	router.POST("/cleantracks/api/register", register)
	router.POST("/cleantracks/api/login", login)
	protected := router.Group("/cleantracks/api")
	protected.Use(authenticate)
	{
		protected.POST("checkLogin", checkLogin)
	}

	router.Run(":9999")
}
