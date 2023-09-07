package bimport

import "product_storage/internal/bridge"

type Bridge struct {
	Date bridge.Date
}

type TestBridge struct {
	Date *bridge.MockDate
}
