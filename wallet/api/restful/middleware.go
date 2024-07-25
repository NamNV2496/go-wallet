package restful

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-wallet/internal/token"
)

var (
	authorizationheader = "Authorization"
	bearerToken         = "Bearer"
	authorPayloadKey    = "authorization_key"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationheader)
		if len(authHeader) == 0 {
			err := fmt.Errorf("invalid Authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			err := fmt.Errorf("invalid Authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		if fields[0] != bearerToken {
			err := fmt.Errorf("invalid type of token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		tokenString := fields[1]
		payload, err := tokenMaker.VerifyToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Set(authorPayloadKey, payload)
		ctx.Next()
	}
}
