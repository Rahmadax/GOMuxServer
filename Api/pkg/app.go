package pkg

import (
	"database/sql"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/gorilla/mux"
	"net/http"
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

	app.Router = newRouter(config)

	app.addEndpoints(config)


	//dbClient, err := sql.Open("mySql", "")
	//if err != nil {
	//	return nil, err
	//}
	//
	//s.dbClient = dbClient

	return app, nil
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Router.Mux.ServeHTTP(w, r)
}

func newRouter(config conf.Configuration) *Router {
	router := &Router{
		Mux: mux.NewRouter(),
	}

	return router
}

func (app *App) addEndpoints(config conf.Configuration) {
	app.addGuestListsRoutes(config.Routes)
	app.addGuestRoutes(config.Routes)
	app.addInvitationRoutes(config.Routes)

}
