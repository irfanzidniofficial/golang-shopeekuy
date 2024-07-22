package main

import (
	"database/sql"
	"golang-shopeekuy/src/util/config"
	"golang-shopeekuy/src/util/routes"

	"github.com/thedevsaddam/renderer"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	sqlDb, err := config.ConnectToDatabase(config.Connection{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		return
	}
	defer sqlDb.Close()

	render := renderer.New()
	routes := setupRoutes(render, sqlDb)
	routes.Run(cfg.AppPort)
}

func setupRoutes(render *renderer.Render, myDb *sql.DB) *routes.Routes {

	return &routes.Routes{}
}
