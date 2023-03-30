package nyx

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	App             fyne.App
	DB              repository.Repository
	MainWindow      fyne.Window
	TopRightToolbar *widget.Toolbar
	GlobalContent   *fyne.Container
	*HomeTab
	*CopiedTradersTab
	*LeaderboardTab
	*AnalyticsTab
	*ProfilesTab
	SiteList []*sites.SupportedSite
	Logger   *utils.AppLogger
	User     *user.User
	*http.Client
	*client.TLSClient
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
	dbPath := strings.TrimSpace(os.Getenv("DB_PATH"))
	if dbPath != "" && filepath.Ext(dbPath) == ".db" {
		path = dbPath
	} else {
		path = filepath.Join(app.App.Storage().RootURI().Path(), "sql.db")
	}

	if utils.DebugEnabled {
		app.Logger.Debug("DB path: " + path)
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

func (app *Config) downloadFile(URL, fileName string) error {

	if URL == "" {
		return errors.New("empty URL")
	}

	// get the response bytes from calling a url
	response, err := app.Client.Get(URL)
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New("received wrong response code when downloading image")
	}

	//open a file for writing
	file, err := os.Create(fmt.Sprintf("downloads/%s.jpg", fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (app *Config) copyTrader(trader user.Trader) error {

	app.User.CopiedTradersManager.Traders = append(app.User.CopiedTradersManager.Traders, trader)

	//TODO: call the actual ByBit APIs
	positions, err := app.fetchTraderPositions(trader.EncryptedUid)
	if err != nil {
		return err
	}

	for _, p := range positions {
		if !utils.Contains(app.SelectedProfile.BlacklistCoins, p.Symbol) {
			app.User.CopiedTradersManager.Positions = append(app.User.CopiedTradersManager.Positions)
		}
	}

	return nil
}
