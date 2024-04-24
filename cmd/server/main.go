package main

import (
	"log"
	"net/http"

	"github.com/Digivate-Labs-Pvt-Ltd/dvlutil"
	"github.com/dvl-mukesh/go-workshop/internal/comment"
	"github.com/dvl-mukesh/go-workshop/internal/config"
	"github.com/dvl-mukesh/go-workshop/internal/database"
	"github.com/dvl-mukesh/go-workshop/internal/middleware"
	transportHTTP "github.com/dvl-mukesh/go-workshop/internal/transport/http"
)

type App struct {
}

func (app *App) Run() error {
	log.Println("Settting up our APP")

	var envVars config.Environment

	if err := dvlutil.ReadEnvVars(&envVars); err != nil {
		return err
	}

	db, err := database.NewDatabase(&envVars)
	if err != nil {
		return err
	}

	err = database.MigrateDB(db)
	if err != nil {
		return err
	}

	commentService := comment.NewService(db)
	handler := transportHTTP.NewHandler(commentService)

	handler.SetupRoutes()
	log.Printf("Starting API server on PORT %s\n", envVars.Port)
	log.Println("Server Started...")

	stack := middleware.CreateStack(
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":" + envVars.Port,
		Handler: stack(handler.Router),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Println("Failed to setup server")
		return err
	}

	return nil
}

func main() {
	log.Println("GO REST API Course")
	app := App{}
	if err := app.Run(); err != nil {
		log.Println("Error starting up REST API")
		log.Fatal(err)
	}
}
