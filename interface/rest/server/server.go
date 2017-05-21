package server

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gobonoid/svc-recipes/interface/rest/handler"
	"github.com/labstack/echo"
	echoMiddleware "github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"github.com/sandalwing/echo-logrusmiddleware"
	"golang.org/x/net/context"
)

const (
	recipesPath = "/recipes"
)

type RecipesServer struct {
	echo *echo.Echo
	port int
}

func NewRecipesServer(port int, log *logrus.Logger, handler handler.RecipesHandler) *RecipesServer {
	e := echo.New()
	e.Logger = logrusmiddleware.Logger{Logger: log}
	e.HideBanner = true
	e.Use(logrusmiddleware.Hook())
	e.Use(echoMiddleware.Recover())
	s := &RecipesServer{port: port}

	recipes := e.Group(recipesPath)

	//Production like project  would use binding and validating middleware, I really value my time here
	recipes.POST("", handler.CreateRecipe)
	recipes.GET("", handler.GetRecipesList)
	recipes.PUT("/:recipeID", handler.UpdateRecipe)
	recipes.GET("/:recipeID", handler.GetRecipe)
	recipes.POST("/:recipeID/rates", handler.RateRecipe)
	s.echo = e
	return s
}

func (s *RecipesServer) Start() {
	go func() {
		if err := s.echo.Start(fmt.Sprintf(":%d", s.port)); err != nil {
			s.echo.Logger.Fatal(errors.Wrap(err, "Failed to start Recipes Server"))
		}
		s.echo.Logger.Info("Recipes Server started")
	}()
}

func (s *RecipesServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.echo.Shutdown(ctx); err != nil {
		s.echo.Logger.Fatal(errors.Wrap(err, "failed to shutdown Recipes Server"))
	}
}
