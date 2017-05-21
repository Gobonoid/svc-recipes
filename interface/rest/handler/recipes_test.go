package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gobonoid/svc-recipes/interface/rest/handler"
	"github.com/gobonoid/svc-recipes/model"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestRecipesHandler_CreateRecipe(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(`{"created_at": "30/06/2015 17:58:00"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	//Normaly I would mock it but this is anyway in the memory
	model := model.NewRecipesModel()
	h := handler.NewRecipesHandler(model)

	c := e.NewContext(req, rec)
	if assert.NoError(t, h.CreateRecipe(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestRecipesHandler_CreateRecipe_InvalidRecipe(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(`/06/2015 17:58:00"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	//Normaly I would mock it but this is anyway in the memory
	model := model.NewRecipesModel()
	h := handler.NewRecipesHandler(model)

	c := e.NewContext(req, rec)
	assert.Error(t, h.CreateRecipe(c))
	assert.Equal(t, http.StatusBadRequest, h.CreateRecipe(c).(*echo.HTTPError).Code)
}
