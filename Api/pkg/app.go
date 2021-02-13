package pkg

import (
	"database/sql"
	"fmt"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/guests"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/guests_repository"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/system_validator"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/Rahmadax/GOMuxServer/Api/pkg/empty_seats"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/guest_list"
	_ "github.com/go-sql-driver/mysql"
)

type App struct {
	Config    conf.Configuration
	Router    *mux.Router
}

func NewApp(config conf.Configuration) (*App, error) {
	app := &App{
		Config: config,
		Router: mux.NewRouter(),
	}

	validator := system_validator.NewSystemValidator(config)

	newDb, err := newDBClient(config.Database)
	if err != nil {
		return nil, err
	}

	guestsRepo := guests_repository.NewGuestsRepository(newDb)

	guestsService := guests.NewGuestsService(config, guestsRepo)
	guests.AddGuestsRoutes(config.Routes, guestsService, validator, app.Router)

	guestListService := guest_list.NewGuestListService(config, guestsRepo)
	guest_list.AddGuestListRoutes(config.Routes, guestListService, validator, app.Router)

	emptySeatsService := empty_seats.NewEmptySeatsService(config, guestsRepo)
	empty_seats.AddGuestListRoutes(config, emptySeatsService, app.Router)

	return app, nil
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Router.ServeHTTP(w, r)
}

func newDBClient(dbConfig conf.DatabaseConfig) (*sql.DB, error) {
	fmt.Println(fmt.Sprintf("Connecting to DB %s", dbConfig.Database))

	db, err := sql.Open(dbConfig.Driver, formatConnectionString(dbConfig))
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(dbConfig.MaxConnLifeTimeMinute)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)

	return db, nil
}

func formatConnectionString(dbConfig conf.DatabaseConfig) string {
	return fmt.Sprintf(
		"%s:%s@/%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Database,
	)
}
