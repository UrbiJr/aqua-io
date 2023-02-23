package nyx

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/nyx/internal/client"
	"github.com/UrbiJr/nyx/internal/repository"
	"github.com/UrbiJr/nyx/internal/sites"
	"github.com/UrbiJr/nyx/internal/user"
	"github.com/UrbiJr/nyx/internal/utils"
)

// Config is the container of the main app, it contains the main attributes
type Config struct {
	App                     fyne.App
	DB                      repository.Repository
	MainWindow              fyne.Window
	ProfileGroupsList       *widget.List
	CheckoutsChartContainer *fyne.Container
	SiteList                []*sites.SupportedSite
	Logger                  *utils.AppLogger
	User                    *user.User
	HTTPClient              *http.Client
	TLSClient               *client.Client
}

// NewNyx returns a new instance of the app
func NewNyx() *Config {
	sitelist := []*sites.SupportedSite{
		{Name: sites.Kickz, Category: sites.SneakerSite, CSVFields: []string{"PID", "MIN SIZE", "MAX SIZE", "PROFILE", "MODE", "REGION"}},
	}

	nyx := &Config{
		SiteList: sitelist,
	}

	return nyx
}

// Quit exits the app gracefully
func (app *Config) Quit() {
	app.Logger.QuitLogger()
}

func (app *Config) ConnectSQL() (*sql.DB, error) {
	path := ""

	if strings.HasSuffix(os.Getenv("DB_PATH"), ".db") {
		path = os.Getenv("DB_PATH")
	} else {
		path = app.App.Storage().RootURI().Path() + "/sql.db"
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *Config) SetupDB(sqldb *sql.DB) {
	app.DB = repository.NewSQLiteRepository(sqldb)

	err := app.DB.Migrate()
	if err != nil {
		app.Logger.Error(err)
		log.Panic()
	}
}
