package response

import (
	"product_storage/internal/entity/global"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrResponse struct {
	Error string // текст ошибки
}
type SuccessResponse struct {
	Data interface{}
}

// NewErrorResponse возвращает статус и ошибку
func NewErrorResponse(err error) ErrResponse {
	log := logrus.New()
	log.Error(err)

	return ErrResponse{
		Error: global.ErrInternalError.Error(),
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
