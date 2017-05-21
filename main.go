package main

import (
	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
	"github.com/gobonoid/svc-recipes/interface/rest/handler"
	"github.com/gobonoid/svc-recipes/interface/rest/server"
	"github.com/gobonoid/svc-recipes/model"
	"github.com/pkg/errors"
)

const (
	applicationPort = 8080 // maybe a flag...
	csvPath         = "recipe-data.csv" //possibly I could use flag here
)

func main() {
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}

	csv, err := os.Open(csvPath)
	if err != nil {
		logger.Fatalf("%#v", errors.Wrapf(err, "can't load csv: %s", csvPath))
	}
	recipesModel := model.NewRecipesModel()
	recipesModel.LoadFromCSV(csv)
	httpServer := server.NewRecipesServer(applicationPort, logger, handler.NewRecipesHandler(recipesModel))
	httpServer.Start()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	httpServer.Stop()
}
