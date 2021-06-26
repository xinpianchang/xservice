package echox

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/xinpianchang/xservice/pkg/log"
	"github.com/xinpianchang/xservice/pkg/responsex"
	"github.com/xinpianchang/xservice/pkg/tracingx"
)

func HTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	traceId := tracingx.GetTraceID(c.Request().Context())

	if he, ok := err.(*echo.HTTPError); ok {
		switch err {
		case middleware.ErrJWTMissing:
			_ = c.JSON(http.StatusUnauthorized, echo.Map{"message": "login required"})
		default:
			status := he.Code
			message := fmt.Sprintf("%v", he.Message)
			if he.Internal != nil {
				message = fmt.Sprintf("%v, cause: %v", he.Message, he.Internal)
			}
			_ = responsex.R(c, responsex.New(status, message, nil).SetHttpStatus(status))
		}
		return
	} else if ve, ok := err.(*validator.ValidationErrors); ok {
		_ = responsex.R(c, responsex.New(http.StatusBadRequest, ve.Error(), nil).SetHttpStatus(http.StatusOK))
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		_ = responsex.R(c, responsex.New(http.StatusNotFound, err.Error(), nil).SetHttpStatus(http.StatusNotFound))
		return
	} else if ve, ok := err.(*responsex.Error); ok {
		httpStatus := http.StatusOK
		if ve.HttpStatus != http.StatusOK && ve.HttpStatus != 0 {
			httpStatus = ve.HttpStatus
		}
		if httpStatus >= 500 {
			sentryecho.GetHubFromContext(c).CaptureException(err)
		}
		v := responsex.New(ve.Status, ve.Message, nil)
		if ve.Internal != nil {
			v = v.SetData(map[string]interface{}{
				"internalErr": fmt.Sprint(ve.Internal),
				"requestId":   c.Get(echo.HeaderXRequestID),
				"traceId":     traceId,
			})
		}
		v = v.SetHttpStatus(httpStatus)
		_ = responsex.R(c, v)
		return
	} else {
		_ = c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":     err.Error(),
			"requestId": c.Get(echo.HeaderXRequestID),
			"traceId":   traceId,
		})
	}

	log.For(c.Request().Context()).Warn(
		c.Response().Header().Get(echo.HeaderXRequestID),
		zap.Any("method", c.Request().Method),
		zap.Int("status", c.Response().Status),
		zap.Any("url", c.Request().URL),
		zap.Any("type", reflect.TypeOf(err)),
		zap.String("error", err.Error()),
	)

	sentryecho.GetHubFromContext(c).CaptureException(err)
}
