package copy_io

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
	"github.com/UrbiJr/copy-io/internal/client"
	"github.com/UrbiJr/copy-io/internal/repository"
	"github.com/UrbiJr/copy-io/internal/sites"
	"github.com/UrbiJr/copy-io/internal/user"
	"github.com/UrbiJr/copy-io/internal/utils"
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

// NewApp returns a new instance of the app
func NewApp() *Config {

	config := &Config{}

	return config
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

func (app *Config) copyTrader(trader user.Trader, profile *user.Profile) error {
	var createdOrders []string
	positions, err := app.fetchTraderPositions(trader.EncryptedUid)
	if err != nil {
		return err
	}

	// create order for each trader's position
	for _, p := range positions {
		app.Logger.Debug(fmt.Sprintf("creating order for symbol %s for profile %s", p.Symbol, profile.Title))
		orderId, err := app.createOrder(profile, p.Symbol, "Buy", "Market", "0.1", p.EntryPrice)
		if err == nil {
			createdOrders = append(createdOrders, orderId)
			// TODO: send notification about created order
		}
	}

	_, err = app.DB.InsertTrader(trader)
	if err != nil {
		app.Logger.Error(err)
	} else {
		app.User.CopiedTradersManager.Traders = append(app.User.CopiedTradersManager.Traders, trader)
	}

	// refresh affected widgets
	app.CopiedTradersList.Refresh()

	return nil
}

func (app *Config) stopCopyingTrader(trader user.Trader) {
	err := app.DB.DeleteTrader(trader.EncryptedUid)
	if err != nil {
		app.Logger.Error(err)
	}

	app.User.CopiedTradersManager.RemoveTraderByUid(trader.EncryptedUid)

	// refresh affected widgets
	app.CopiedTradersList.Refresh()
	app.RefreshLeaderboardWithoutFetch()
}
