package rpc

import (
	"mime/multipart"
	"product_storage/internal/entity/global"
	"strings"

	"github.com/gin-gonic/gin"
)

// Context контекст rpc
type Context struct {
	GinContext *gin.Context
	Request    *Request
	Response   *Response
	abort      bool
	next       bool
}

// Abort прекращение работы обработчиков в методе
func (c *Context) Abort() {
	c.abort = true
}

// Next переключиться на следующий обработчик
func (c *Context) Next() {
	c.next = true
}

// Reset сброс сигналов
func (c *Context) Reset() {
	c.abort = false
	c.next = false
}

// IsAbort сигнал о прекращении работы обработчиков в методе
func (c *Context) IsAbort() bool {
	return c.abort
}

// IsNext сигнал о переключении на следующий обработчик
func (c *Context) IsNext() bool {
	return c.next
}

// ReturnResult вернуть результат
func (c *Context) ReturnResult(result interface{}) {
	c.Response = SuccessResponse(result)
	c.Abort()
}

// ReturnSuccessNullResult вернуть успешный пустой результат
func (c *Context) ReturnSuccessNullResult() {
	c.Response = SuccessResponse("OK")
	c.Abort()
}

// ReturnStateResult вернуть результат
func (c *Context) ReturnStateResult(state string, result interface{}) {
	c.Response = SuccessResponse(gin.H{
		"data": NamedResult{Name: state, Result: result},
		"meta": MetaResponse{
			Version: 1,
		},
	})
	c.Abort()
}

// ReturnManyStateResult вернуть множественный результат
func (c *Context) ReturnManyStateResult(namedResults ...NamedResult) {
	c.Response = SuccessResponse(gin.H{
		"data": namedResults,
		"meta": MetaResponse{
			Version: 2,
		},
	})
	c.Abort()
}

// ReturnError вернуть ошибку
func (c *Context) ReturnError(err error) {
	var code int
	switch err {
	case ErrJSONParseError:
		code = RPCParseError
	case global.ErrNeedAuth:
		code = RPCAuthError
	case global.ErrParamsIncorrect:
		code = RPCInvalidParams
	case global.ErrAccessRight:
		code = RPCAccessRightError
	case global.ErrNoData500:
		code = RPCNoDataError
	default:
		code = RPCInternalError
	}

	c.Response = ErrorResponse(Error{Code: code, Message: firstUpper(err.Error())})
	c.Abort()
}

func firstUpper(s string) string {
	return strings.Replace(s, string([]rune(s)[0]), strings.ToUpper(string([]rune(s)[0])), 1)
}

// MustGetUserID получить oper_id, должен обязательно пройти через авт.middleware
func (c *Context) MustGetUserID() int {
	return c.GinContext.MustGet(global.UserIDKey).(int)
}

// MustGetUserLogin получить oper_login, должен обязательно пройти через авт.middleware
func (c *Context) MustGetUserLogin() string {
	return c.GinContext.MustGet(global.UserLoginKey).(string)
}

// GetGinFile получить gin file
func (c *Context) GetGinFile() (*multipart.FileHeader, error) {
	return c.GinContext.FormFile("file")
}

// UserAgent получить user agent
func (c *Context) UserAgent() string {
	return c.GinContext.GetHeader(headerUserAgent)
}

// MustGetToken получить token, должен обязательно пройти через авт.middleware
func (c *Context) MustGetToken() string {
	return c.GinContext.MustGet(global.TokenKey).(string)
}
