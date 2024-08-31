package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo/v5"
)

func (a *Api) V1(c echo.Context) error {

	targetUrl, _ := url.Parse("http://localhost:8080")

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)
	c.Request().Header.Add("x-api-key", "secret")

	proxy.ServeHTTP(c.Response(), c.Request())
	return nil

}

func (a *Api) V2(c echo.Context) error {

	targetUrl, _ := url.Parse("http://localhost:8081")

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	proxy.ServeHTTP(c.Response(), c.Request())
	return nil

}
func (a *Api) Hello2(c echo.Context) error {

	fmt.Println(c)

	obj := map[string]interface{}{"message": "qqqqqqqqqqqqqqqqq"}
	fmt.Println("sssssssssssssssss")
	return c.JSON(http.StatusOK, obj)
}
