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
	Token string `json:"token"`
}

func login(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusOK, utils.BuildApiResponse(1, "Invalid request", nil))
		return
	}

	existedUser, err := user.FindUser(creds.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(existedUser.Password), []byte(creds.Password)) != nil {
		c.JSON(http.StatusOK, utils.BuildApiResponse(2, "Invalid email or password", nil))
		return
	}

	expirationTime := time.Now().Add(30 * 24 * time.Hour)
	claims := &Claims{
		Email: creds.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusOK, utils.BuildApiResponse(2, "Error generating token", nil))
		return
	}

	c.SetCookie("token", tokenString, 2592000, "/", "localhost", false, true)

	c.JSON(http.StatusOK, utils.BuildApiResponse(0, "", LoginResp{Token: tokenString}))
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
	c.JSON(http.StatusOK, utils.BuildApiResponse(0, "", nil))
}

func main() {
	router := gin.Default()

	logs.InitLog()
	utils.Init()

	router.POST("cleantracks/api/webhook", payment.WebhookHandler)
	router.POST("cleantracks/api/register", register)
	router.POST("/cleantracks/api/login", login)
	protected := router.Group("cleantracks/api")
	protected.Use(authenticate)
	{
		protected.POST("checkLogin", checkLogin)
		// protected.POST("api/subscribe", subscribe)
	}

	router.Run(":9999")
}
