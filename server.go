package echoserver

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/rianekacahya/config"
	custom_middleware "github.com/rianekacahya/echoserver/middleware"
)

var (
	server *echo.Echo
	mutex  sync.Once
)

func GetServer() *echo.Echo {
	mutex.Do(func() {
		server = newServer()
	})
	return server
}

func newServer() *echo.Echo {
	return echo.New()
}

func InitServer() {

	// Hide banner
	GetServer().HideBanner = true

	// Set debug status parameter
	GetServer().Debug = config.GetEchoServerDebug()

	// init default middleware
	GetServer().Use(
		middleware.Recover(),
		custom_middleware.CORS(),
		custom_middleware.Headers(),
	)

	// healthCheck endpoint
	GetServer().GET("/infrastructure/healthcheck", func (c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	// custom error handler
	GetServer().HTTPErrorHandler = Handler
}

func StartServer(ctx context.Context) {
	select {
	case <-ctx.Done():
		if err := GetServer().Shutdown(ctx); err != nil {}
	default:
		if err := GetServer().StartServer(&http.Server{
			Addr:         config.GetEchoServerPort(),
			ReadTimeout:  time.Duration(config.GetEchoServerReadTimeout()) * time.Second,
			WriteTimeout: time.Duration(config.GetEchoServerWriteTimeout()) * time.Second,
		}); err != nil {
			panic(err)
		}
	}
}