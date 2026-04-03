package route

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/assimon/luuu/controller/comm"
	"github.com/assimon/luuu/middleware"
	"github.com/labstack/echo/v4"
)

// RegisterRoute 路由注册
func RegisterRoute(e *echo.Echo) {
	e.Any("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello epusdt, https://github.com/GMwalletApp/epusdt")
	})

	payRoute := e.Group("/pay")
	payRoute.GET("/checkout-counter/:trade_id", comm.Ctrl.CheckoutCounter)
	payRoute.GET("/check-status/:trade_id", comm.Ctrl.CheckStatus)

	// payment routes
	paymentRoute := e.Group("/payments")

	// for epusdt
	epusdtV1 := paymentRoute.Group("/epusdt/v1")
	epusdtV1.POST("/order/create-transaction", func(ctx echo.Context) error {
		// add default token and currency for old plugin

		body := make(map[string]interface{})
		if err := ctx.Bind(&body); err != nil {
			return comm.Ctrl.FailJson(ctx, err)
		}
		if _, ok := body["token"]; !ok {
			body["token"] = "usdt"
		}
		if _, ok := body["currency"]; !ok {
			body["currency"] = "cny"
		}
		if _, ok := body["network"]; !ok {
			body["network"] = "TRON"
		}
		ctx.Set("request_body", body)

		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return comm.Ctrl.FailJson(ctx, err)
		}
		ctx.Request().Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

		return comm.Ctrl.CreateTransaction(ctx)
	}, middleware.CheckApiSign())

	// gmpay v1 routes
	gmpayV1 := paymentRoute.Group("/gmpay/v1")
	gmpayV1.POST("/order/create-transaction", comm.Ctrl.CreateTransaction, middleware.CheckApiSign())

}
