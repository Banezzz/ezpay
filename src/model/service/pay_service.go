package service

import (
	"github.com/Banezzz/ezpay/config"
	"github.com/Banezzz/ezpay/model/data"
	"github.com/Banezzz/ezpay/model/response"
	"github.com/Banezzz/ezpay/util/constant"
)

// ErrOrder is returned when checkout initialization cannot find the trade id.
var ErrOrder = constant.OrderNotExists

// GetCheckoutCounterByTradeId returns checkout initialization data for an existing order.
// It does not decide the payment state; callers should use CheckStatus for that.
func GetCheckoutCounterByTradeId(tradeId string) (*response.CheckoutCounterResponse, error) {
	orderInfo, err := data.GetOrderInfoByTradeId(tradeId)
	if err != nil {
		return nil, err
	}
	if orderInfo.ID <= 0 {
		return nil, ErrOrder
	}

	resp := &response.CheckoutCounterResponse{
		TradeId:        orderInfo.TradeId,
		Amount:         orderInfo.Amount,
		ActualAmount:   orderInfo.ActualAmount,
		Token:          orderInfo.Token,
		Currency:       orderInfo.Currency,
		ReceiveAddress: orderInfo.ReceiveAddress,
		Network:        orderInfo.Network,
		ExpirationTime: orderInfo.CreatedAt.AddMinutes(config.GetOrderExpirationTime()).TimestampMilli(),
		RedirectUrl:    orderInfo.RedirectUrl,
		CreatedAt:      orderInfo.CreatedAt.TimestampMilli(),
		IsSelected:     checkoutOrderIsSelected(orderInfo),
	}
	return resp, nil
}
