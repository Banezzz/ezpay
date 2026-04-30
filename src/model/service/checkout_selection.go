package service

import "github.com/GMWalletApp/ezpay/model/mdb"

func checkoutOrderIsSelected(order *mdb.Orders) bool {
	if order == nil {
		return false
	}
	return order.IsSelected
}
