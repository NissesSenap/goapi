package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

/* TODO, create a generl function
that takes in a struct that contains
method http.Method
body string
ParamNames string
ParamValues string

Return a c.context

Need a way to get a "gobal echo.New() to get that to work.
Is it possible with a
func TestMain(m *testing.M) {
	m.Run()
	os.Remove(dbPath)
}
?
*/

var (
	userMethods  = []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodHead, http.MethodDelete, http.MethodOptions}
	usersMethods = []string{http.MethodGet, http.MethodPost, http.MethodHead, http.MethodOptions}
)

func TestGetUserOptions(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodOptions, "/users", http.NoBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, UserOptions(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, strings.Join(userMethods, ","), rec.Header().Get("Allow"))
	}
}

func TestGetUsersOptions(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodOptions, "/users", http.NoBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, UsersOptions(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, strings.Join(usersMethods, ","), rec.Header().Get("Allow"))
	}
}

/*
func TestGetUserID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/", http.NoBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Generate ID
	id := bson.NewObjectId()
	t.Logf("The generated ID is: %v", id)
	c.SetParamNames("id")
	c.SetParamValues(string(id))

}
*/
