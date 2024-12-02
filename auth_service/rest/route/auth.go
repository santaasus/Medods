package route

import (
	repository "Medods/auth_service/inner_layer/repository/user"
	"Medods/auth_service/rest/adapter"

	"github.com/gin-gonic/gin"
)

const (
	AUTH_GROUP          = "/auth"
	TOKENS_PATH         = "/tokens"
	REFRESH_TOKENS_PATH = "/refresh_token"
)

func AuthRoutes(router *gin.Engine) {
	baseAdapter := adapter.BaseAdapter{
		Repository: repository.Repository{},
	}

	controller := baseAdapter.AuthAdapter()

	group := router.Group(AUTH_GROUP)
	{
		group.GET(TOKENS_PATH, controller.GetTokens)
		group.POST(REFRESH_TOKENS_PATH, controller.RefreshToken)
	}
}
