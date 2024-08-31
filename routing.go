package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func (a *Api) InitRouting(publicDirFlag string) {

	a.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {



		config := middleware.RateLimiterConfig{
			Skipper: middleware.DefaultSkipper,
			Store: middleware.NewRateLimiterMemoryStoreWithConfig(
				middleware.RateLimiterMemoryStoreConfig{Rate: float64(a.Config.RateLimiterMemoryStore.Rate), Burst: a.Config.RateLimiterMemoryStore.Burst, ExpiresIn: a.Config.RateLimiterMemoryStore.ExpiresIn},
			),
			IdentifierExtractor: func(ctx echo.Context) (string, error) {
				id := ctx.RealIP()
				for _, ip := range a.Config.RateLimiterMemoryStore.BlacklistIPs {
					if id == ip {
						return "", fmt.Errorf("IP %s is blacklisted", id)
					}
				}

				return id, nil
			},
			ErrorHandler: func(context echo.Context, err error) error {
				return context.JSON(http.StatusForbidden, nil)
			},
			DenyHandler: func(context echo.Context, identifier string, err error) error {


				a.App.Logger().Warn("DenyHandler for host " + identifier)
				return context.JSON(http.StatusTooManyRequests, nil)
			},
		}

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		e.Router.Use(middleware.RateLimiterWithConfig(config), middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogStatus:   true,
			LogURI:      true,
			LogError:    true,
			HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				if v.Error == nil {
					logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
						slog.String("uri", v.URI),
						slog.Int("status", v.Status),
					)
				} else {
					logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
						slog.String("uri", v.URI),
						slog.Int("status", v.Status),
						slog.String("err", v.Error.Error()),
					)
				}
				return nil
			},
		}))

		api := e.Router.Group("/api")

		all_methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}

		for _, v := range a.Config.Services {

			for _, r := range v.Routes {
				targetServices, err := url.Parse(v.Protocol + "://" + v.Host + ":" + v.Port)
				if err != nil {
					break
				}
				api.Match(all_methods, r.Paths, func(c echo.Context) error {

					proxy := httputil.NewSingleHostReverseProxy(targetServices)
					c.Request().Header.Add("x-api-key", "secret")

					proxy.ServeHTTP(c.Response(), c.Request())
					return nil
				}, v.GetPremessionType())

			}

		}

		// serves static files from the provided public dir (if exists)
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS(publicDirFlag), true))

		return nil
	})
}
