package bimport

import "product_storage/tools/datefunctions"

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
