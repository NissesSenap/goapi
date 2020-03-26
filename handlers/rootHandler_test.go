package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	rootMessage = "Running echo API v1"
)

func TestGetRoot(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertion
	if assert.NoError(t, Root(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, rootMessage, rec.Body.String())
	}
	/*
		err := Root(c)
		if err != nil {
			t.Fatalf("Test failed, got error: %v", err)
		}
		t.Log(rec.Code)

		if rec.Code != http.StatusCreated {
			t.Errorf("Got %v, expected: %v", rec.Code, http.StatusOK)
		}

		if rec.Body.String() != rootMessage {
			t.Errorf("Got %v, expected: %v", rec.Body.String(), rootMessage)
		}
	*/
}
