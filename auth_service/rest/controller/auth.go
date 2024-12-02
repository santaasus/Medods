package controller

import (
	service "Medods/auth_service/inner_layer/service/auth"
	"errors"
	"net/http"

	domainErrors "github.com/santaasus/errors-handler"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service *service.Service
}

func (c *Controller) GetTokens(ctx *gin.Context) {
	params := ctx.Query("guid")
	if params == "" {
		err := errors.New("guid is missed")
		_ = ctx.Error(err)
		return
	}

	clientIP := ctx.ClientIP()
	tokensModel, err := c.Service.GetNewTokens(params, clientIP)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, tokensModel)
}

func (c *Controller) RefreshToken(ctx *gin.Context) {
	var accessToken AccessTokenRequest

	if err := ctx.BindJSON(&accessToken); err != nil {
		appError := domainErrors.ThrowAppErrorWith(domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	clientIP := ctx.ClientIP()
	tokenModel, err := c.Service.AccessTokenByRefreshToken(accessToken.RefreshToken, clientIP)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, tokenModel)
}
