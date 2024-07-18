package utils

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
