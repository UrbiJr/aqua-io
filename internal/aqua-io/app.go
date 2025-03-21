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
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/aqua-io/internal/client"
	"github.com/UrbiJr/aqua-io/internal/repository"
	"github.com/UrbiJr/aqua-io/internal/resources"
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
	SplashWindow    fyne.Window
	MainWindow      fyne.Window
	TopRightToolbar *widget.Toolbar
	AccountMenu     *fyne.Container
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
	if app.LoginWindow != nil {
		app.LoginWindow.Close()
	}
	if app.MainWindow != nil {
		app.MainWindow.Close()
	}
	os.Exit(0)
}

// Logout exits the current user and shows the login window
func (app *Config) Logout() {
	app.Logger.Debug("logging out")
	err := app.DB.DeleteAllUsers()
	if err != nil {
		app.Logger.Error(err)
		app.Quit()
	}
	app.User = nil
	app.MakeTray()
	app.LoginWindow.Show()
	app.MainWindow.Hide()
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
		return err
	}
	defer file.Close()

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

// copyTrader creates an order for each trader's position
func (app *Config) copyTrader(trader user.Trader, profile *user.Profile) error {
	positions, err := app.fetchTraderPositions(trader.EncryptedUid)
	if err != nil {
		return err
	}

	if len(positions) > 0 {
		var forms []*fyne.Container
		for _, p := range positions {
			if p.Amount < 0 {
				// get the open position form for each position
				forms = append(forms, app.openPositionForm(profile, utils.SHORT_POSITION, p.Symbol, p.MarkPrice))
			} else {
				forms = append(forms, app.openPositionForm(profile, utils.LONG_POSITION, p.Symbol, p.MarkPrice))
			}

		}
		index := 0
		// get the window which will show the form
		createOrderWindow := app.App.NewWindow(fmt.Sprintf("Copying %s's positions (%d/%d)", trader.NickName, index+1, len(forms)))
		content := container.NewVBox()
		// add the first form to window content
		content.Add(forms[index])
		var btn *widget.Button
		if len(forms) == 1 {
			btn = widget.NewButtonWithIcon("Close", theme.CancelIcon(), func() {
				createOrderWindow.Close()
			})
		} else {
			btn = widget.NewButtonWithIcon("Next", theme.NavigateNextIcon(), nil)
			btn.OnTapped = func() {
				if index+1 < len(forms) {
					// second-last element
					if index+1 == len(forms)-1 {
						btn.SetIcon(theme.CancelIcon())
						btn.SetText("Close")
					}
					index++
					// update window content based on current form index
					content.Objects[0] = forms[index]
					createOrderWindow.SetTitle(fmt.Sprintf("Copying %s's positions (%d/%d)", trader.NickName, index+1, len(forms)))
					content.Refresh()
				} else {
					// all forms have been showed, we can close the window
					createOrderWindow.Close()
				}
			}
		}

		createOrderWindow.SetContent(container.NewVBox(
			content,
			layout.NewSpacer(),
			btn,
		))
		content.Refresh()
		createOrderWindow.Resize(fyne.NewSize(520, 460))
		createOrderWindow.SetFixedSize(true)
		createOrderWindow.SetIcon(resources.ResourceIconPng)
		createOrderWindow.Show()
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
	app.refreshCopiedTradersTab()
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
	app.refreshCopiedTradersList()
	app.RefreshLeaderboardWithoutFetch()
	app.refreshCopiedTradersTab()
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
