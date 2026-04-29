package service

import (
	"strings"

	"github.com/GMWalletApp/epusdt/model/data"
)

const SupportedPaymentToken = "USDT"

func IsSupportedPaymentToken(token string) bool {
	return strings.EqualFold(strings.TrimSpace(token), SupportedPaymentToken)
}

func IsSupportedPaymentAsset(network, token string) bool {
	if !IsSupportedPaymentToken(token) {
		return false
	}
	row, err := data.GetEnabledChainTokenBySymbol(network, token)
	return err == nil && row.ID > 0
}
