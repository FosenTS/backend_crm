package authorization

import (
	"backend_crm/internal/controller/http/fasthttp/authorization/dto"
	"backend_crm/internal/model"
	"backend_crm/internal/usecase/users"
	"encoding/json"
	"errors"
	"strings"

	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

type Controller struct {
	users  users.Usecase
	logger zerolog.Logger
}

func NewController(users users.Usecase, logger zerolog.Logger) *Controller {
	return &Controller{
		users:  users,
		logger: logger,
	}
}

func (c *Controller) Register(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.Error("Only POST method allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	body := ctx.PostBody()
	if len(body) == 0 {
		ctx.Error("Empty request body", fasthttp.StatusBadRequest)
		return
	}

	var register *dto.Register
	if err := json.Unmarshal(body, &register); err != nil {
		ctx.Error("Invalid JSON format", fasthttp.StatusBadRequest)
		return
	}

	if err := c.users.Register(ctx, &model.Register{
		RoleId:   model.Role(register.RoleId),
		Username: register.Username,
		Password: register.Password,
	}); err != nil {
		ctx.Error("Error on the server", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)
}

func (c *Controller) Login(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.Error("Only POST method allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	body := ctx.PostBody()
	if len(body) == 0 {
		ctx.Error("Empty request body", fasthttp.StatusBadRequest)
		return
	}

	var login *dto.Login
	if err := json.Unmarshal(body, &login); err != nil {
		ctx.Error("Invalid JSON format", fasthttp.StatusBadRequest)
		return
	}

	tokens, err := c.users.Login(ctx, &model.Login{
		Username: login.Username,
		Password: login.Password,
	})
	if err != nil {
		if errors.Is(err, users.ErrNotFoundUser) {
			ctx.Error("user not found", fasthttp.StatusUnauthorized)
			return
		} else if errors.Is(err, users.ErrIncorrectPassword) {
			ctx.Error("incorrect password", fasthttp.StatusUnauthorized)
			return
		}
		c.logger.Error().Err(err).Msg("Error on the server")
		ctx.Error("Error on the server", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	if err := json.NewEncoder(ctx).Encode(&dto.Token{
		Access:  tokens.AccessToken,
		Refresh: tokens.RefreshToken,
	}); err != nil {
		c.logger.Error().Err(err).Msg("Error creating response")
		ctx.Error("Error creating response", fasthttp.StatusInternalServerError)
	}
}

func (c *Controller) Access(ctx *fasthttp.RequestCtx) {
	if !ctx.IsGet() {
		ctx.Error("Only GET method allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	authHeader := string(ctx.Request.Header.Peek("Authorization"))
	if authHeader == "" {
		ctx.Error("Empty Authorization Header", fasthttp.StatusUnauthorized)
		return
	}

	// Extract token from "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		ctx.Error("Invalid Authorization Header format", fasthttp.StatusUnauthorized)
		return
	}

	if _, _, err := c.users.CheckAccess(ctx, parts[1]); err != nil {
		if errors.Is(err, users.ErrExpiredAccessToken) {
			ctx.Error("Expired access token", fasthttp.StatusUnauthorized)
			return
		}

		c.logger.Error().Err(err).Msg("Error on the server")
		ctx.Error("Error on the server", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (c *Controller) Refresh(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.Error("Only POST method allowed", fasthttp.StatusMethodNotAllowed)
		return
	}

	body := ctx.PostBody()
	if len(body) == 0 {
		ctx.Error("Empty request body", fasthttp.StatusBadRequest)
		return
	}

	var refresh *dto.Refresh
	if err := json.Unmarshal(body, &refresh); err != nil {
		ctx.Error("Invalid JSON format", fasthttp.StatusBadRequest)
		return
	}

	tokens, err := c.users.RefreshTokens(ctx, refresh.Refresh)
	if err != nil {
		ctx.Error("Error on the server", fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	if err := json.NewEncoder(ctx).Encode(&dto.Token{
		Access:  tokens.AccessToken,
		Refresh: tokens.RefreshToken,
	}); err != nil {
		ctx.Error("Error creating response", fasthttp.StatusInternalServerError)
	}
}

func (c *Controller) AuthMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		authHeader := string(ctx.Request.Header.Peek("Authorization"))
		if authHeader == "" {
			ctx.Error("Empty Authorization Header", fasthttp.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.Error("Invalid Authorization Header format", fasthttp.StatusUnauthorized)
			return
		}

		userId, userRole, err := c.users.CheckAccess(ctx, parts[1])
		if err != nil {
			if errors.Is(err, users.ErrExpiredAccessToken) {
				ctx.Error("Expired access token", fasthttp.StatusUnauthorized)
				return
			}

			ctx.Error("Error on the server", fasthttp.StatusInternalServerError)
			return
		}

		ctx.SetUserValue("user_id", userId)
		ctx.SetUserValue("user_role", userRole)

		next(ctx)
	}
}
