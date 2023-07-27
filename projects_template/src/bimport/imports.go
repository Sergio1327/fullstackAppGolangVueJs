package bimport

import "projects_template/tools/datefunctions"

type BridgeImports struct {
	Bridge Bridge
}

func (b *BridgeImports) InitBridge() {
	b.Bridge = Bridge{
		Date: datefunctions.NewDateTool(),
	}
}

func NewEmptyBridge() *BridgeImports {
	return &BridgeImports{}
}
