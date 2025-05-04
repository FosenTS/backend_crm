package fasthttp

import (
	"backend_crm/internal/controller/http/fasthttp/app"
	"backend_crm/internal/controller/http/fasthttp/authorization"
	"backend_crm/internal/controller/http/fasthttp/orders"
	"context"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type controller struct {
	authorization authorization.Controller
	orders        orders.Contoller
	app           app.Controller
}

func NewController(
	auth authorization.Controller,
	orders orders.Contoller,
	app app.Controller,
) *controller {
	return &controller{
		authorization: auth,
		orders:        orders,
		app:           app,
	}
}

func (c *controller) addAuthMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return c.authorization.AuthMiddleware(next)
}

func (c *controller) Handlers(ctx context.Context) fasthttp.RequestHandler {
	r := router.New()

	apiV1 := r.Group("/api/v1")
	apiV1.GET("/app", c.app.GetFile)

	orders := apiV1.Group("/orders")
	orders.GET("/{status}", c.addAuthMiddleware(c.orders.Orders))
	orders.POST("/order/{orderId}", c.addAuthMiddleware(c.orders.UpdateOrder))
	orders.POST("/new-order", c.addAuthMiddleware(c.orders.NewOrder))

	auth := apiV1.Group("/auth")
	auth.GET("/access", c.authorization.Access)
	auth.POST("/refresh", c.authorization.Refresh)
	auth.POST("/login", c.authorization.Login)
	auth.POST("/registration", c.addAuthMiddleware(c.authorization.Register))

	return r.Handler
}
