package middleware

import (
	"net/http"
	"time"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/rianekacahya/logger"
	"go.uber.org/zap"
)

func Logger (c echo.Context, reqBody, resBody []byte) {
	if c.Response().Status >= http.StatusOK && c.Response().Status < http.StatusMultipleChoices {
		logger.Info(
			"http",
			zap.Int("status", c.Response().Status),
			zap.String("time", time.Now().Format(time.RFC1123Z)),
			zap.String("hostname", c.Request().Host),
			zap.String("user_agent", c.Request().UserAgent()),
			zap.String("method", c.Request().Method),
			zap.String("path", c.Path()),
			zap.String("query", c.QueryString()),
			zap.Any("req", json.RawMessage(reqBody)),
			zap.Any("res", json.RawMessage(resBody)),
		)
	}
}
