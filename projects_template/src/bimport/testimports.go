package bimport

import (
	"projects_template/internal/bridge"

	"github.com/golang/mock/gomock"
)

type TestBridgeImports struct {
	ctrl       *gomock.Controller
	TestBridge TestBridge
}

func NewTestBridgeImports(
	ctrl *gomock.Controller,
) *TestBridgeImports {
	return &TestBridgeImports{
		ctrl: ctrl,
		TestBridge: TestBridge{
			Date: bridge.NewMockDate(ctrl),
		},
	}
}

func (t *TestBridgeImports) BridgeImports() *BridgeImports {
	return &BridgeImports{
		Bridge: Bridge{
			Date: t.TestBridge.Date,
		},
	}
}
