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
	Router   *Router
	dbClient *sql.DB
}

type Router struct {
	Mux *mux.Router
}

func NewApp(config conf.Configuration) (*App, error) {
	app := &App{}

	app.newRouter()
	app.addEndpoints(config)

	app.newDBClient(config.Database)

	return app, nil
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Router.Mux.ServeHTTP(w, r)
}

func (app *App) newRouter(){
	app.Router = &Router{
		Mux: mux.NewRouter(),
	}
}

func (app *App) addEndpoints(config conf.Configuration) {
	app.addGuestListsRoutes(config.Routes)
	app.addGuestRoutes(config.Routes)
	app.addInvitationRoutes(config.Routes)
	app.addSeatsRoutes(config.Routes)
}

func (app *App) newDBClient(dbConfig conf.DatabaseConfig) {
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
