package api

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/saifutdinov/go-invoices-api/pkg/db"
	"github.com/saifutdinov/go-invoices-api/pkg/toml"

	// required implement to use postgres
	_ "github.com/lib/pq"
)

func RunServer(configPath string) {
	config, err := toml.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// init echo server
	echoServer := echo.New()

	// init db connection, migrate and return back prepared connection
	dbo := initDB(echoServer, config.DBDriver, config.DBConnectionString, "./db/migrations")

	initHandlers(echoServer, dbo, config)

	log.Fatal(echoServer.Start("0.0.0.0:5000"))

}

func initDB(echoServer *echo.Echo, dbDriver, dbString, migrationPath string) *sql.DB {
	dbs, err := sql.Open(dbDriver, dbString)
	if err != nil {
		echoServer.Logger.Fatal(err)
	}
	if err = dbs.Ping(); err != nil {
		echoServer.Logger.Fatal(err)
	}
	echoServer.Logger.Infof("DB opened")
	db.MigrateDB(dbs, migrationPath)
	echoServer.Logger.Infof("DB migration complete")
	return dbs
}
