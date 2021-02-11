package pkg

import (
	"database/sql"
	"fmt"
	"github.com/Rahmadax/GOMuxServer/Api/conf"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/guests"
	"github.com/Rahmadax/GOMuxServer/Api/pkg/guests_repo"
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

	guestsRepo := guests_repo.NewGuestsRepository(newDBClient(config.Database))

	guestsService := guests.NewGuestsService(config, guestsRepo)
	guests.AddGuestsRoutes(config.Routes, guestsService, validator, app.Router)

	guestListService := guest_list.NewGuestListService(config, guestsRepo)
	guest_list.AddGuestListRoutes(config.Routes, guestListService, validator, app.Router)

	emptySeatsService := empty_seats.NewEmptySeatsService(config, guestsRepo)
	empty_seats.AddGuestListRoutes(config.Routes, emptySeatsService, app.Router)

	return app, nil
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Router.ServeHTTP(w, r)
}

func newDBClient(dbConfig conf.DatabaseConfig) *sql.DB {
	fmt.Println(fmt.Sprintf("Connecting to DB %s", dbConfig.Database))

	db, err := sql.Open(dbConfig.Driver, formatConnectionString(dbConfig))
	if err != nil {
		panic(err)
	}

	db.SetConnMaxIdleTime(dbConfig.MaxConnLifeTimeMinute)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)

	return db
}

func formatConnectionString(dbConfig conf.DatabaseConfig) string {
	return fmt.Sprintf(
		"%s:%s@/%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Database,
	)
}
