package dao

import (
	"sync"

	"github.com/assimon/luuu/model/mdb"
	"github.com/gookit/color"
)

var once sync.Once

// 自动建表
func MdbTableInit() {
	once.Do(func() {
		if err := Mdb.AutoMigrate(&mdb.Orders{}); err != nil {
			color.Red.Printf("[store_db] AutoMigrate DB(Orders),err=%s\n", err)
			// panic(err)
			return
		}
		if err := Mdb.AutoMigrate(&mdb.WalletAddress{}); err != nil {
			color.Red.Printf("[store_db] AutoMigrate DB(Orders),err=%s\n", err)
			// panic(err)
			return
		}
		if err := Mdb.AutoMigrate(&mdb.SupportedAsset{}); err != nil {
			color.Red.Printf("[store_db] AutoMigrate DB(SupportedAsset),err=%s\n", err)
			return
		}

		seedSupportedAssets()
	})
}

func seedSupportedAssets() {
	var cnt int64
	if err := Mdb.Model(&mdb.SupportedAsset{}).Count(&cnt).Error; err != nil {
		color.Red.Printf("[store_db] Count SupportedAsset,err=%s\n", err)
		return
	}
	if cnt > 0 {
		return
	}
	defaults := []mdb.SupportedAsset{
		{Network: mdb.NetworkTron, Token: "TRX", Status: mdb.TokenStatusEnable},
		{Network: mdb.NetworkTron, Token: "USDT", Status: mdb.TokenStatusEnable},
		{Network: mdb.NetworkSolana, Token: "SOL", Status: mdb.TokenStatusEnable},
		{Network: mdb.NetworkSolana, Token: "USDT", Status: mdb.TokenStatusEnable},
		{Network: mdb.NetworkSolana, Token: "USDC", Status: mdb.TokenStatusEnable},
		{Network: mdb.NetworkEthereum, Token: "USDT", Status: mdb.TokenStatusEnable},
		{Network: mdb.NetworkEthereum, Token: "USDC", Status: mdb.TokenStatusEnable},
		{Network: mdb.NetworkBsc, Token: "USDT", Status: mdb.TokenStatusEnable},
		{Network: mdb.NetworkBsc, Token: "USDC", Status: mdb.TokenStatusEnable},
		{Network: mdb.NetworkPolygon, Token: "USDT", Status: mdb.TokenStatusEnable},
		{Network: mdb.NetworkPolygon, Token: "USDC", Status: mdb.TokenStatusEnable},
		{Network: mdb.NetworkPlasma, Token: "USDT", Status: mdb.TokenStatusEnable},
	}
	if err := Mdb.Create(&defaults).Error; err != nil {
		color.Red.Printf("[store_db] Seed SupportedAsset,err=%s\n", err)
	}
}
