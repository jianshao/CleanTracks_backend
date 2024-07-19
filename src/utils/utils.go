package utils

import (
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/utils/logs"
	"github.com/jianshao/chrome-exts/CleanTracks/backend/src/utils/prisma"
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

func Init() {
	prisma.Init()
	logs.InitLog()
}
