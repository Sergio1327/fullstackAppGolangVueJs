package rpc

// FuncHandler функция обработчика
type FuncHandler func(*Context)

// Request запрос
type Request struct {
	JSON   string      `json:"jsonrpc"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// BatchRequest запрос
type BatchRequest []Request

// Response ответ
type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *Error      `json:"error,omitempty"`
}

// BatchResponse массовый ответ
type BatchResponse []Response

// Error ошибка
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SuccessResponse успешный ответ
func SuccessResponse(result interface{}) *Response {
	return &Response{
		JSONRPC: "2.0",
		Result:  result,
		Error:   nil,
	}
}

// ErrorResponse ответ при ошибке
func ErrorResponse(err Error) *Response {
	return &Response{
		JSONRPC: "2.0",
		Result:  nil,
		Error:   &err,
	}
}

// MetaResponse мета информация в ответе
type MetaResponse struct {
	Version int `json:"version"`
}

// NamedResult именованный результат
type NamedResult struct {
	Name   string      `json:"name"`
	Result interface{} `json:"result"`
}

// MethodErrors ошибки методов
type MethodErrors map[string]*Error

// MethodNames извлечения списка методов через запятую
func MethodNames(data BatchRequest) string {
	mn := ""

	for index, row := range data {
		if index != 0 {
			mn += ", "

		}
		mn += row.Method
	}

	return mn
}

// RequestP часть запроса для декодирования
type RequestP struct {
	JSON   string `json:"jsonrpc"`
	Method string `json:"method"`
}
