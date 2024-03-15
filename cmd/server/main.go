package main

import (
	"fmt"
	"net/http"

	"github.com/dvl-mukesh/go-workshop/internal/comment"
	"github.com/dvl-mukesh/go-workshop/internal/database"
	transportHTTP "github.com/dvl-mukesh/go-workshop/internal/transport/http"
)

type App struct {
}

func (app *App) Run() error {
	fmt.Println("Settting up our APP")
	db, err := database.NewDatabase()
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

	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		fmt.Println("Failed to setup server")
		return err
	}

	return nil
}

func main() {
	fmt.Println("GO REST API Course")
	app := App{}
	if err := app.Run(); err != nil {
		fmt.Println("Error starting up REST API")
		fmt.Println(err)
	}
}
