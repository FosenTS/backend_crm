package orders

import (
	"backend_crm/internal/controller/http/fasthttp/orders/dto"
	"backend_crm/internal/model"
	"backend_crm/internal/repository/orders"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

type Contoller struct {
	orders orders.Repository
	logger zerolog.Logger
}

func NewController(orders orders.Repository, logger zerolog.Logger) *Contoller {
	return &Contoller{
		orders: orders,
		logger: logger,
	}
}

func (c *Contoller) UpdateOrder(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.Error("Only PATCH method allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	orderId, ok := ctx.UserValue("orderId").(string)
	if !ok {
		ctx.Error("Invalid request", fasthttp.StatusBadRequest)
		return
	}

	body := ctx.PostBody()
	if len(body) == 0 {
		ctx.Error("Empty request body", fasthttp.StatusBadRequest)
		return
	}

	var st *dto.Status
	if err := json.Unmarshal(body, &st); err != nil {
		ctx.Error("Invalid JSON format", fasthttp.StatusBadRequest)
		return
	}

	if err := c.orders.UpdateOrderStatus(ctx, orderId, model.OrderStatus(st.Status)); err != nil {
		ctx.Error("Error on the server", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (c *Contoller) Orders(ctx *fasthttp.RequestCtx) {
	if !ctx.IsGet() {
		ctx.Error("Only GET method allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	status, ok := ctx.UserValue("status").(model.OrderStatus)
	if !ok {
		ctx.Error("Invalid request", fasthttp.StatusBadRequest)
		return
	}

	queryArgs := ctx.QueryArgs()
	phoneFilter := false
	var phone string
	p := queryArgs.Peek("phone")
	if len(p) != 0 {
		phone = string(p)
		phoneFilter = true
	}

	emailFilter := false
	var email string
	e := queryArgs.Peek("email")
	if len(p) != 0 {
		email = string(e)
		emailFilter = true
	}

	userRole := ctx.UserValue("user_role").(model.Role)
	if !ok {
		ctx.Error("Invalid request", fasthttp.StatusBadRequest)
		return
	}

	if userRole == model.Director {
		var orders []*model.Order
		var err error
		if phoneFilter && emailFilter {
			orders, err = c.orders.GetByStatusAndPhoneAndEmail(ctx, status, phone, email)
			if err != nil {
				ctx.Error("Error on the server", fasthttp.StatusInternalServerError)
				return
			}
		} else if phoneFilter {
			orders, err = c.orders.GetByStatusAndPhone(ctx, status, phone)
			if err != nil {
				ctx.Error("Error on the server", fasthttp.StatusInternalServerError)
				return
			}
		} else if emailFilter {
			orders, err = c.orders.GetByStatusAndEmail(ctx, status, email)
			if err != nil {
				ctx.Error("Error on the server", fasthttp.StatusInternalServerError)
				return
			}
		} else {
			orders, err = c.orders.GetByStatus(ctx, status)
			if err != nil {
				ctx.Error("Error on the server", fasthttp.StatusInternalServerError)
				return
			}
		}

		respOrders := make([]*dto.Order, 0, len(orders))
		for _, order := range orders {
			respOrders = append(respOrders, &dto.Order{
				Phone:       order.Phone,
				Email:       order.Email,
				Description: order.Description,
				Product: dto.Product{
					Name:        order.Product.Name,
					Weigth:      fmt.Sprintf("%f kg", order.Product.Weigth),
					Description: order.Product.Description,
				},
				Status: int(order.Status),
			})
		}

		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusOK)
		if err := json.NewEncoder(ctx).Encode(respOrders); err != nil {
			ctx.Error("Error creating response", fasthttp.StatusInternalServerError)
		}
		return
	}

	userId := ctx.UserValue("user_id").(string)

	var orders []*model.Order
	var err error

	if phoneFilter && emailFilter {
		orders, err = c.orders.GetByUserIdAndStatusAndPhoneAndEmail(ctx, userId, status, phone, email)
		if err != nil {
			ctx.Error("Error on the server", fasthttp.StatusInternalServerError)
			return
		}
	} else if phoneFilter {
		orders, err = c.orders.GetByUserIdAndStatusAndPhone(ctx, userId, status, phone)
		if err != nil {
			ctx.Error("Error on the server", fasthttp.StatusInternalServerError)
			return
		}
	} else if emailFilter {
		orders, err = c.orders.GetByUserIdAndStatusAndEmail(ctx, userId, status, email)
		if err != nil {
			ctx.Error("Error on the server", fasthttp.StatusInternalServerError)
			return
		}
	} else {
		orders, err = c.orders.GetByUserIdAndStatus(ctx, userId, status)
		if err != nil {
			ctx.Error("Error on the server", fasthttp.StatusInternalServerError)
			return
		}
	}

	respOrders := make([]*dto.Order, 0, len(orders))
	for _, order := range orders {
		respOrders = append(respOrders, &dto.Order{
			Phone:       order.Phone,
			Email:       order.Email,
			Description: order.Description,
			Product: dto.Product{
				Name:        order.Product.Name,
				Weigth:      fmt.Sprintf("%f kg", order.Product.Weigth),
				Description: order.Product.Description,
			},
			Status: int(order.Status),
		})
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	if err := json.NewEncoder(ctx).Encode(respOrders); err != nil {
		ctx.Error("Error creating response", fasthttp.StatusInternalServerError)
	}
}

func (c *Contoller) NewOrder(ctx *fasthttp.RequestCtx) {
	if ctx.IsPost() {
		ctx.Error("Only POST method allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	body := ctx.PostBody()
	if len(body) == 0 {
		ctx.Error("Empty request body", fasthttp.StatusBadRequest)
		return
	}

	var newOrder *dto.NewOrder
	if err := json.Unmarshal(body, &newOrder); err != nil {
		ctx.Error("Invalid JSON format", fasthttp.StatusBadRequest)
		return
	}

	if err := c.orders.Save(ctx, &model.NewOrder{
		Phone:       newOrder.Phone,
		Email:       newOrder.Email,
		Description: newOrder.Description,
		ProductId:   newOrder.ProductId,
		Status:      model.Consideration,
	}); err != nil {
		ctx.Error("Error on the server", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
}
