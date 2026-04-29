package service

import (
	"strings"

	"github.com/GMWalletApp/epusdt/model/mdb"
)

func checkoutOrderIsSelected(order *mdb.Orders) bool {
	if order == nil {
		return false
	}
	return order.IsSelected || strings.TrimSpace(order.ReceiveAddress) != ""
}
