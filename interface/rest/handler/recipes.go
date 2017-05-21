package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gobonoid/svc-recipes/model"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const (
	Limit = "limit"
	Page  = "page"
)

type RecipesHandler struct {
	recipesAggregator model.RecipesAggregator
}

func NewRecipesHandler(recipesAggregator model.RecipesAggregator) RecipesHandler {
	return RecipesHandler{
		recipesAggregator: recipesAggregator,
	}
}

func (h RecipesHandler) CreateRecipe(c echo.Context) error {
	recipe := &model.Recipe{}
	if err := c.Bind(recipe); err != nil {
		return err
	}
	recipe.Id = time.Now().Nanosecond()
	if err := h.recipesAggregator.CreateRecipe(recipe); err == model.DuplicateError {
		return echo.NewHTTPError(http.StatusConflict, "Recipe already exists")
	}
	return c.NoContent(http.StatusCreated)
}

func (h RecipesHandler) GetRecipesList(c echo.Context) error {
	recipes := h.recipesAggregator.FetchRecipes(recipesListLimiter(c))
	return c.JSON(http.StatusOK, recipes)
}

func (h RecipesHandler) GetRecipe(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("recipeID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect recipeID given")
	}
	recipe, err := h.recipesAggregator.FetchOneByID(id)
	if err == model.NotFoundError {
		return echo.NewHTTPError(http.StatusNotFound, "Recipe not found")
	}
	return c.JSON(http.StatusOK, recipe)
}

func (h RecipesHandler) UpdateRecipe(c echo.Context) error {
	recipe := &model.Recipe{}
	if err := c.Bind(recipe); err != nil {
		return err
	}
	id, err := strconv.Atoi(c.Param("recipeID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect recipeID given")
	}

	if err := h.recipesAggregator.UpdateRecipe(id, recipe); err == model.NotFoundError {
		return echo.NewHTTPError(http.StatusNotFound, "Recipe not found")
	}
	return c.NoContent(http.StatusOK)
}

func (h RecipesHandler) RateRecipe(c echo.Context) error {
	recipeRate := &model.RecipeRate{}
	if err := c.Bind(recipeRate); err != nil {
		return err
	}
	id, err := strconv.Atoi(c.Param("recipeID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Incorrect recipeID given")
	}
	if h.recipesAggregator.RateRecipe(id, recipeRate) == model.NotFoundError {
		return echo.NewHTTPError(http.StatusNotFound, "Recipe not found")
	}
	return c.NoContent(http.StatusCreated)
}

func recipesListLimiter(c echo.Context) *model.Limiter {
	var limit int
	var page int
	if l := c.QueryParam(Limit); l != "" {
		tmp, err := strconv.Atoi(l)
		if err != nil {
			c.Logger().Errorf("%#v", errors.Wrap(err, "Failed to parse given limit"))
		}
		limit = tmp
	}
	if p := c.QueryParam(Page); p != "" {
		tmp, err := strconv.Atoi(p)
		if err != nil {
			c.Logger().Errorf("%#v", errors.Wrap(err, "Failed to parse given page"))
		}
		page = tmp
	}
	if page == 0 {
		page = 1
	}
	return &model.Limiter{
		Limit: limit,
		Page:  page,
	}
}
