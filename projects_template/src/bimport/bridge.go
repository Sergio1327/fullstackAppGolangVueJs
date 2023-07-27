package bimport

import "projects_template/internal/bridge"

type Bridge struct {
	Date bridge.Date
}

type TestBridge struct {
	Date *bridge.MockDate
}
