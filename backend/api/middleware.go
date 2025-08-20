package api

import (
	"amulyam/token"
	"amulyam/utils"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationType       = "header"
	authorizationPayloadKey = "authorization_payload"
)

func (server *Server) authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !server.config.RequireAuth {
			ctx.Next()
		}
		fmt.Println("Middleware")
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			fmt.Println("error")
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}

		authorizationtype := strings.ToLower(fields[0])

		if authorizationtype != authorizationTypeBearer {
			err := errors.New("authorization type is not supported")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}
		payload, err := server.tokenMaker.VerifyToken(fields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(err))
			return
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func (server *Server) isAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("Authorisation middleware")
		authorizationPayload, exists := ctx.Get(authorizationPayloadKey)
		if !exists {
			ctx.JSON(http.StatusForbidden, utils.ErrorResponse(errors.New("unable to find the the user information")))
			return
		}
		payload, ok := authorizationPayload.(*token.Payload)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(errors.New("unable to get the payload")))
			return
		}
		fmt.Println(!payload.IsAdmin)
		if !payload.IsAdmin {
			ctx.AbortWithStatusJSON(http.StatusForbidden, utils.ErrorResponse(errors.New("you should have admin privilleges to access this resource")))
			return
		}
		ctx.Next()
	}
}
