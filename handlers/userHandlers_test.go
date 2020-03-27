package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/NissesSenap/GoAPI/user"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
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
		assert.Equal(t, strings.Join(usersMethods, ",("), rec.Header().Get("Allow"))
	}
}

func mockGetOne(id bson.ObjectId) (*user.User, error) {
	u := new(user.User)
	u.ID = id
	u.Name = "Mark"
	u.Role = "lead developer"

	return u, nil
}
func TestGetUserID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/", http.NoBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Generate ID
	id := bson.NewObjectId()
	t.Logf("The generated ID is: %v string: %v", id, id.String())
	c.SetParamNames("id")
	c.SetParamValues(id.String())

	uf := NewUserFunction(mockGetOne)

	if assert.NoError(t, uf.UsersGetOne(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// {"user":{"id":"5e7cbed74044734d097c98e3","name":"Mar1k","role":"lead developer"}}

	// This might be due to: https://github.com/labstack/echo/issues/1492
}
