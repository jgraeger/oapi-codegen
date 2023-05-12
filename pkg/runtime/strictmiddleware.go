package runtime

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bunrouter"
)

type StrictEchoHandlerFunc func(ctx echo.Context, request interface{}) (response interface{}, err error)

type StrictEchoMiddlewareFunc func(f StrictEchoHandlerFunc, operationID string) StrictEchoHandlerFunc

type StrictHttpHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error)

type StrictHttpMiddlewareFunc func(f StrictHttpHandlerFunc, operationID string) StrictHttpHandlerFunc

type StrictGinHandlerFunc func(ctx *gin.Context, request interface{}) (response interface{}, err error)

type StrictGinMiddlewareFunc func(f StrictGinHandlerFunc, operationID string) StrictGinHandlerFunc

type StrictBunrouterHandlerFunc func(ctx context.Context, w http.ResponseWriter, r bunrouter.Request, request interface{}) (response interface{}, err error)

type StrictBunrouterMiddlewareFunc func(f StrictBunrouterHandlerFunc, operationID string) StrictBunrouterHandlerFunc
