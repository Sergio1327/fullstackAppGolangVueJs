package response

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        // успешно ли выполнилась операция
	Error   string      //	текст ошибки
	Data    interface{} //	данные которые нужно вывести
}

// NewErrorResponse возвращает статус и ошибку
func NewErrorResponse(err error) Response {
	return Response{
		Success: false,
		Error:   err.Error(),
		Data:    nil,
	}
}

// NewSuccessResponse возвращает статус и данные
func NewSuccessResponse(data interface{}, idType string) Response {
	return Response{
		Success: true,
		Error:   "null",
		Data: gin.H{
			idType: data,
		},
	}
}
