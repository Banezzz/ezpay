package mdb

type SupportedAsset struct {
	Network string `gorm:"column:network;uniqueIndex:supported_asset_network_token_uindex,priority:1" json:"network"`
	Token   string `gorm:"column:token;uniqueIndex:supported_asset_network_token_uindex,priority:2" json:"token"`
	Status  int64  `gorm:"column:status;default:1" json:"status"`
	BaseModel
}

func (s *SupportedAsset) TableName() string {
	return "supported_asset"
}
