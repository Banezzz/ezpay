package mdb

import "strings"

const (
	TokenStatusEnable  = 1
	TokenStatusDisable = 2
)

const (
	NetworkTron      = "tron"
	NetworkSolana    = "solana"
	NetworkEthereum  = "ethereum"
	NetworkBsc       = "bsc"
	NetworkBscLegacy = "binance"
	NetworkPolygon   = "polygon"
	NetworkPlasma    = "plasma"
)

func NormalizeNetwork(network string) string {
	network = strings.ToLower(strings.TrimSpace(network))
	if network == NetworkBscLegacy {
		return NetworkBsc
	}
	return network
}

func NetworkAliases(network string) []string {
	network = NormalizeNetwork(network)
	if network == "" {
		return nil
	}
	if network == NetworkBsc {
		return []string{NetworkBsc, NetworkBscLegacy}
	}
	return []string{network}
}

func SameNetwork(a, b string) bool {
	return NormalizeNetwork(a) == NormalizeNetwork(b)
}

func IsEVMNetwork(network string) bool {
	switch NormalizeNetwork(network) {
	case NetworkEthereum, NetworkBsc, NetworkPolygon, NetworkPlasma:
		return true
	default:
		return false
	}
}

const (
	WalletSourceManual = "manual"
	WalletSourceImport = "import"
)

type WalletAddress struct {
	Network string `gorm:"column:network;uniqueIndex:wallet_address_network_address_uindex" json:"network" example:"tron"`
	Address string `gorm:"column:address;uniqueIndex:wallet_address_network_address_uindex" json:"address" example:"TTestTronAddress001"`
	// 状态 1=启用 2=禁用
	Status int64  `gorm:"column:status;default:1" json:"status" enums:"1,2" example:"1"`
	Remark string `gorm:"column:remark;size:255" json:"remark" example:"主钱包"`
	Source string `gorm:"column:source;size:16;default:manual" json:"source" enums:"manual,import" example:"manual"`
	BaseModel
}

func (w *WalletAddress) TableName() string {
	return "wallet_address"
}
