package response

import "github.com/gin-gonic/gin"

type ErrResponse struct {
	Error string // текст ошибки
}
type SuccessResponse struct {
	Data interface{}
}

// NewErrorResponse возвращает статус и ошибку
func NewErrorResponse(err error) ErrResponse {
	return ErrResponse{
		Error: err.Error(),
	}
}

// NewSuccessResponse возвращает статус и данные
func NewSuccessResponse(data interface{}, idType string) SuccessResponse {
	return SuccessResponse{
		Data: gin.H{
			idType: data,
		},
	}
}
