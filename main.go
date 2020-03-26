package main

import (
	"crypto/subtle"

	"github.com/NissesSenap/GoAPI/cache"
	"github.com/NissesSenap/GoAPI/handlers"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jsonResponse map[string]interface{}

/* serveCache returns the cached value as an echo middleware
The ninja return func is a way for echo middleware to be able to "stack"
multiple middlewares after eachother.
*/
func serveCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if cache.Serve(c.Response(), c.Request()) {
			return nil
		}
		return next(c)
	}
}

// cacheResponse saves the return value to the cache as an echo middleware
func cacheResponse(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Writer = cache.NewWriter(c.Response().Writer, c.Request())
		return next(c)
	}
}

func auth(username, password string, c echo.Context) (bool, error) {
	// Be careful to use constant time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
		return true, nil
	}
	return false, nil
}

func main() {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	// The Recover middleware will handle panics and dump a stack trace
	e.Use(middleware.Recover())
	// The Secure middleware will help to secure the page for example protecting against XSS attack
	e.Use(middleware.Secure())

	// Enable logging
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"method": "${method}", "uri": "${uri}", "status": "${status}", "latency_human": "${latency_human}"}` + "\n",
	}))
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	e.GET("/", handlers.Root)

	u := e.Group("/users")

	u.OPTIONS("", handlers.UsersOptions)
	u.HEAD("", handlers.UsersGetAll, serveCache)
	u.GET("", handlers.UsersGetAll, serveCache, cacheResponse)
	u.POST("", handlers.UsersPostOne, middleware.BasicAuth(auth))

	uid := u.Group("/:id")

	uid.OPTIONS("", handlers.UserOptions)
	uid.HEAD("", handlers.UsersGetOne, serveCache)
	uid.GET("", handlers.UsersGetOne, serveCache, cacheResponse)
	uid.PUT("", handlers.UsersPutOne, middleware.BasicAuth(auth), cacheResponse)
	uid.PATCH("", handlers.UsersPatchOne, middleware.BasicAuth(auth), cacheResponse)
	uid.DELETE("", handlers.UsersDeleteOne, middleware.BasicAuth(auth))

	e.Logger.Fatal(e.Start(":12345"))
}
