package main

import (
	"net/http"
	"strings"

	"github.com/NissesSenap/GoAPI/cache"
	"github.com/NissesSenap/GoAPI/user"
	"github.com/asdine/storm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/mgo.v2/bson"
)

type jsonResponse map[string]interface{}

func usersOptions(c echo.Context) error {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodOptions}
	c.Response().Header().Set("Allow", strings.Join(methods, ","))
	return c.NoContent(http.StatusOK)
}

func userOptions(c echo.Context) error {
	// TODO automate the user & usersOptions to get it from the echo framework.
	methods := []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodHead, http.MethodDelete, http.MethodOptions}
	c.Response().Header().Set("Allow", strings.Join(methods, ","))
	return c.NoContent(http.StatusOK)
}

func usersGetAll(c echo.Context) error {
	if cache.Serve(c.Response(), c.Request()) {
		return nil
	}
	users, err := user.All()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if c.Request().Method == http.MethodHead {
		return c.NoContent(http.StatusOK)
	}
	c.Response().Writer = cache.NewWriter(c.Response().Writer, c.Request())
	return c.JSON(http.StatusOK, jsonResponse{"users": users})
}

func usersGetOne(c echo.Context) error {
	if cache.Serve(c.Response(), c.Request()) {
		return nil
	}
	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	id := bson.ObjectIdHex(c.Param("id"))
	u, err := user.One(id)
	if err != nil {
		if err == storm.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if c.Request().Method == http.MethodHead {
		return c.NoContent(http.StatusOK)
	}
	c.Response().Writer = cache.NewWriter(c.Response().Writer, c.Request())
	return c.JSON(http.StatusOK, jsonResponse{"user": u})
}

func usersDeleteOne(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	id := bson.ObjectIdHex(c.Param("id"))
	err := user.Delete(id)
	if err != nil {
		if err == storm.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)

		}
		return echo.NewHTTPError(http.StatusInternalServerError)

	}
	cache.Drop("/users")
	cache.Drop(cache.MakeResource(c.Request()))
	return c.NoContent(http.StatusOK)
}

func usersPutOne(c echo.Context) error {
	// TODO catch if id exist or not
	u := new(user.User)
	err := c.Bind(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)

	}
	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	id := bson.ObjectIdHex(c.Param("id"))
	u.ID = id
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)

	}
	cache.Drop("/users")
	c.Response().Writer = cache.NewWriter(c.Response().Writer, c.Request())
	return c.JSON(http.StatusOK, jsonResponse{"user": u})
}

func usersPatchOne(c echo.Context) error {
	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	id := bson.ObjectIdHex(c.Param("id"))
	u, err := user.One(id)
	if err != nil {
		if err == storm.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)

		}
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	err = c.Bind(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	// u.ID = id
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	cache.Drop("/users")
	c.Response().Writer = cache.NewWriter(c.Response().Writer, c.Request())
	return c.JSON(http.StatusOK, jsonResponse{"user": u})
}

func usersPostOne(c echo.Context) error {
	u := new(user.User)
	err := c.Bind(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	u.ID = bson.NewObjectId()
	err = u.Save()
	if err != nil {
		if err == user.ErrRecordInvalid {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	cache.Drop("/users")
	c.Response().Header().Set("Location", "/users/"+u.ID.Hex())
	return c.NoContent(http.StatusCreated)
}

func root(c echo.Context) error {
	return c.String(http.StatusOK, "Running echo API v1")
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

	e.GET("/", root)

	u := e.Group("/users")

	u.OPTIONS("", usersOptions)
	u.HEAD("", usersGetAll)
	u.GET("", usersGetAll)
	u.POST("", usersPostOne)

	uid := u.Group("/:id")

	uid.OPTIONS("", userOptions)
	uid.HEAD("", usersGetOne)
	uid.GET("", usersGetOne)
	uid.PUT("", usersPutOne)
	uid.PATCH("", usersPatchOne)
	uid.DELETE("", usersDeleteOne)

	e.Logger.Fatal(e.Start(":12345"))
}
