package pkg

import (
	"database/sql"
	"fmt"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/gorilla/mux"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	Config   conf.Configuration
	Router   *mux.Router
	dbClient *sql.DB
}

func NewApp(config conf.Configuration) (*App, error) {
	app := &App{
		Config: config,
		Router: mux.NewRouter(),
	}
	app.addEndpoints()
	app.newDBClient()

	return app, nil
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Router.ServeHTTP(w, r)
}

func (app *App) addEndpoints() {
	app.addGuestListRoutes()
	app.addGuestsRoutes()
	app.addSeatsRoutes()
}

func (app *App) newDBClient() {
	dbConfig := app.Config.Database

	db, err := sql.Open(dbConfig.Driver, formatConnectionString(dbConfig))
	if err != nil {
		panic(err)
	}

	db.SetConnMaxIdleTime(dbConfig.MaxConnLifeTimeMinute)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)

	app.dbClient = db
}

func formatConnectionString(dbConfig conf.DatabaseConfig) string {
	return fmt.Sprintf(
		"%s:%s@/%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Database,
	)
}
