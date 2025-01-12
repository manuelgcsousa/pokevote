package main

import (
	"html/template"
	"log/slog"
	"os"

	"github.com/manuelgcsousa/pokevote/internal/data"
	"github.com/manuelgcsousa/pokevote/internal/server"
	"github.com/manuelgcsousa/pokevote/internal/settings"
)

func main() {
	// Application settings
	settings.Setup()

	// Connect to DB
	db, err := data.NewConnection()
	if err != nil {
		slog.Default().Error("Error while connecting to DB: " + err.Error())
		os.Exit(1)
	}
	defer db.CloseDatabase()

	// Start server
	server := &server.Server{
		Port:     2727,
		Database: db,
		Tmpl:     template.Must(template.ParseGlob("templates/*.html")),
	}
	server.Start()
}
