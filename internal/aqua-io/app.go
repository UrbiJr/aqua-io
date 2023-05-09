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
	"github.com/UrbiJr/aqua-io/internal/client"
	"github.com/UrbiJr/aqua-io/internal/repository"
	"github.com/UrbiJr/aqua-io/internal/sites"
	"github.com/UrbiJr/aqua-io/internal/user"
	"github.com/UrbiJr/aqua-io/internal/utils"
	"github.com/UrbiJr/aqua-io/internal/whop"
)

// Config is the container of the main app, it contains the main attributes
type Config struct {
	App             fyne.App
	Whop            *whop.Whop
	DB              repository.Repository
	LoginWindow     fyne.Window
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

/*
Downloads a file and stores it locally.

URL is the download target.
filename is used to rename the file locally.
ext is the original file extension.
*/
func (app *Config) downloadFile(URL, fileName, ext string) error {

	if URL == "" {
		return errors.New("empty URL")
	}

	// check if file extension is supported
	if ext != ".jpg" && ext != ".png" && ext != ".gif" {
		return errors.New("unsupported file type")
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
	file, err := os.Create(fmt.Sprintf("downloads/%s%s", fileName, ext))
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
	positions, err := app.fetchTraderPositions(trader.EncryptedUid)
	if err != nil {
		return err
	}

	// create order for each trader's position
	success := 0
	for _, p := range positions {
		app.Logger.Debug(fmt.Sprintf("creating order for symbol %s for profile %s", p.Symbol, profile.Title))
		_, err := app.createOrder(profile, p.Symbol, "Market", p.Amount, p.MarkPrice)
		if err == nil {
			success += 1
		} else {
			app.App.SendNotification(&fyne.Notification{
				Title:   "Order create fail :-(",
				Content: err.Error(),
			})
			app.Logger.Error(err.Error())
		}
	}
	if success > 0 {
		app.App.SendNotification(&fyne.Notification{
			Title:   "Orders create success! :D",
			Content: fmt.Sprintf("Successfully created %d orders", success),
		})
	}

	profile.TraderID = trader.EncryptedUid
	err = app.DB.UpdateProfile(profile.ID, *profile)
	if err != nil {
		app.Logger.Error(err)
	} else {
		app.User.ProfileManager.UpdateProfile(profile.ID, *profile)
	}

	// refresh affected widgets
	app.CopiedTradersList.Refresh()
	app.RefreshLeaderboardWithoutFetch()
	app.refreshCopiedTradersTab(true)
	app.refreshProfilesTab()

	return nil
}

func (app *Config) stopCopyingTrader(trader user.Trader, traderID string) error {
	if trader.EncryptedUid == "" {
		if traderID != "" {
			// fetch trader
			traderInfo, err := app.fetchTraderByUid(traderID)
			if err != nil {
				return err
			}
			trader = *traderInfo
		} else {
			return errors.New("no trader provided")
		}
	}

	profile := app.User.ProfileManager.GetProfileByTraderID(trader.EncryptedUid)
	if profile == nil {
		return fmt.Errorf("no profiles found with trader id %s", trader.EncryptedUid)
	}

	profile.TraderID = ""
	err := app.DB.UpdateProfile(profile.ID, *profile)
	if err != nil {
		app.Logger.Error(err.Error())
	} else {
		app.User.ProfileManager.UpdateProfile(profile.ID, *profile)
	}

	// refresh affected widgets
	app.CopiedTradersList.Refresh()
	app.RefreshLeaderboardWithoutFetch()
	app.refreshCopiedTradersTab(false)
	app.refreshProfilesTab()

	return nil
}

func (app *Config) getProfiles() {
	profiles, err := app.DB.AllProfiles()
	if err != nil {
		app.Logger.Error(err)
		app.Quit()
	}
	app.User.ProfileManager.Profiles = profiles
}
