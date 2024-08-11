package utils

import (
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/utils/logs"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/utils/prisma"
	"github.com/sirupsen/logrus"
)

type ApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func BuildApiResponse(code int, message string, data any) ApiResponse {
	return ApiResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func Init(path string, level logrus.Level) bool {
	prisma.Init()
	return logs.InitLog(path, level)
}

func Close() {
	prisma.Close()
}
