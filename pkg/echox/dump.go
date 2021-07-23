package echox

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/xinpianchang/xservice/pkg/log"
)

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Dump(filename ...string) echo.MiddlewareFunc {
	var logger log.Logger
	if len(filename) > 0 {
		logger, _ = log.NewLogger(filename[0])
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var l log.Logger
			if logger == nil {
				l = log.For(c.Request().Context())
			} else {
				l = logger.For(c.Request().Context())
			}

			req, _ := httputil.DumpRequest(c.Request(), true)
			start := time.Now()
			l.Debug(
				c.Response().Header().Get(echo.HeaderXRequestID),
				zap.Any("url", c.Request().URL),
				zap.String("req", string(req)),
			)

			resBody := new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer

			err := next(c)

			if err != nil {
				return err
			}

			l.Debug(
				c.Response().Header().Get(echo.HeaderXRequestID),
				zap.Int64("cost", time.Since(start).Milliseconds()),
				zap.Int("status", c.Response().Status),
				zap.Any("header", c.Response().Header()),
				zap.String("rsp", resBody.String()),
			)

			return nil
		}
	}
}
