package rpc

import "errors"

var (
	// ErrMethodNotFound метод не найден
	ErrMethodNotFound = errors.New("метод не найден")
	// ErrJSONParseError произошла ошибка при парсинге jrpc
	ErrJSONParseError = errors.New("произошла ошибка при парсинге jrpc")
)

const (
	// RPCParseError ошибка при парсинге основного JSON
	RPCParseError = -32700
	// RPCInvalidRequest некорректный запрос
	RPCInvalidRequest = -32600
	// RPCMethodNotFound метод не найден
	RPCMethodNotFound = -32601
	// RPCInvalidParams некорректные параметры
	RPCInvalidParams = -32602
	// RPCInternalError внутреняя ошибка
	RPCInternalError = -32603
	// RPCAuthError ошибка при проверке авторизации
	RPCAuthError = -32000
	// RPCAccessRightError ошибка при проверке прав доступа
	RPCAccessRightError = -32001
	// RPCNoDataError ошибка при отсутствии данных
	RPCNoDataError = -32002
)
