package echoserver

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/rianekacahya/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type response struct {
	Message interface{} `json:"message"`
}

func Handler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)

	if GetServer().Debug {
		msg = err.Error()
		switch err.(type) {
		case *echo.HTTPError:
			code = err.(*echo.HTTPError).Code
		}
	} else {
		switch e := err.(type) {
		case *echo.HTTPError:
			code = e.Code
			msg = e.Message
			if e.Internal != nil {
				msg = fmt.Sprintf("%v, %v", err, e.Internal)
			}
		default:
			msg = http.StatusText(code)
		}

		if _, ok := msg.(string); ok {
			msg = response{Message: msg}
		}
	}

	if !c.Response().Committed {
		if c.Request().Method == "HEAD" {
			c.NoContent(code)
		} else {
			c.JSON(code, msg)
		}

		logger.Error(
			"http",
			zap.Int("status", code),
			zap.String("time", time.Now().Format(time.RFC1123Z)),
			zap.String("hostname", c.Request().Host),
			zap.String("user_agent", c.Request().UserAgent()),
			zap.String("method", c.Request().Method),
			zap.String("path", c.Path()),
			zap.String("query", c.QueryString()),
			zap.Any("req", c.Request().Body),
			zap.Any("res", msg),
		)
	}
}