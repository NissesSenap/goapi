package handlers

import (
	"net/http"
	"strings"

	"github.com/NissesSenap/GoAPI/cache"
	"github.com/NissesSenap/GoAPI/user"
	"github.com/asdine/storm"
	"github.com/labstack/echo/v4"
	"gopkg.in/mgo.v2/bson"
)

// GetOneUser creates a type that is used in the struct
// To be able to be overwritten when mocking the API.
type GetOneUser func(id bson.ObjectId) (*user.User, error)

// UserFunction the struct that contains different types of function
// Depending on what we need to do
type UserFunction struct {
	GetOne GetOneUser
}

// NewUserFunction the wrapper function for GetOne
func NewUserFunction(uf GetOneUser) *UserFunction {
	return &UserFunction{GetOne: uf}
}

// UsersOptions give all the avliable API methods avaliabl in /users
func UsersOptions(c echo.Context) error {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodOptions}
	c.Response().Header().Set("Allow", strings.Join(methods, ","))
	return c.NoContent(http.StatusOK)
}

// UserOptions give all the avliable API methods avaliabl in /users/:id
func UserOptions(c echo.Context) error {
	// TODO automate the user & usersOptions to get it from the echo framework.
	methods := []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodHead, http.MethodDelete, http.MethodOptions}
	c.Response().Header().Set("Allow", strings.Join(methods, ","))
	return c.NoContent(http.StatusOK)
}

// UsersGetAll get all the users avliable in the db
func UsersGetAll(c echo.Context) error {
	users, err := user.All()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if c.Request().Method == http.MethodHead {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, jsonResponse{"users": users})
}

// UsersGetOne gets a single user /users/:id
func (uf *UserFunction) UsersGetOne(c echo.Context) error {

	if !bson.IsObjectIdHex(c.Param("id")) {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	id := bson.ObjectIdHex(c.Param("id"))
	u, err := uf.GetOne(id)
	// u, err := user.One(id)
	if err != nil {
		if err == storm.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if c.Request().Method == http.MethodHead {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, jsonResponse{"user": u})
}

// UsersDeleteOne deletes a single user /users/:id
func UsersDeleteOne(c echo.Context) error {
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

// UsersPutOne puts a single user /users/:id
func UsersPutOne(c echo.Context) error {
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
	return c.JSON(http.StatusOK, jsonResponse{"user": u})
}

// UsersPatchOne patches a single user /users/:id
func UsersPatchOne(c echo.Context) error {
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
	return c.JSON(http.StatusOK, jsonResponse{"user": u})
}

// UsersPostOne creates a single user /users/:id
func UsersPostOne(c echo.Context) error {
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
