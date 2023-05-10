package copy_io

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	x_widget "fyne.io/x/fyne/widget"
	"github.com/UrbiJr/aqua-io/internal/resources"
)

// UI contains fyne elements
type UI struct {
}

// MakeMobileUI instantiates all the needed pages and makes the UI layout, but does not display it yet.
func (app *Config) MakeMobileUI() {

	// add application tabs (home, tasks, proxies, profiles, settings)
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("          ", theme.HomeIcon(), canvas.NewText("Home content goes here", nil)),
		container.NewTabItemWithIcon("          ", theme.ListIcon(), canvas.NewText("Tasks content goes here", nil)),
		container.NewTabItemWithIcon("          ", app.App.Settings().Theme().Icon(resources.IconNameWifi), canvas.NewText("Proxies content goes here", nil)),
		container.NewTabItemWithIcon("          ", app.App.Settings().Theme().Icon(resources.IconNameCreditCard), canvas.NewText("Profiles content goes here", nil)),
	)
	// show tabs at the bottom of the window
	tabs.SetTabLocation(container.TabLocationBottom)

	// populate app tabs

	// populate window
	finalContent := container.NewVBox(tabs)

	app.MainWindow.SetContent(finalContent)

}

// MakeMobileUI instantiates all the needed pages and makes the UI layout, but does not display it yet.
func (app *Config) MakeDesktopUI() {

	greetings := []string{
		"how can Aqua.io assist you today? :-)",
		"how are you going to use Aqua.io today? :-)",
		"ready to have some fun with Aqua.io?",
		"Aqua.io is at your service :salute:",
		"it's been a while.",
		"time to make the checkout feed go brrr.",
		"everyday is a perfect day to run Aqua.io!",
	}
	msg := fmt.Sprintf("Hello %s, %s", app.User.Username, greetings[rand.Intn(len(greetings))])

	// get app tabs content
	app.HomeTab = &HomeTab{}
	app.CopiedTradersTab = &CopiedTradersTab{}
	app.LeaderboardTab = &LeaderboardTab{}
	app.AnalyticsTab = &AnalyticsTab{}
	app.ProfilesTab = &ProfilesTab{}

	// load profiles
	app.getProfiles()
	profilesTabContent := app.profilesTab()
	homeTabContent := app.homeTab(msg)
	copiedTradersTabContent := app.copiedTradersTab()
	leaderboardContent := app.leaderboardTab()

	app.HomeTab.TabItem = container.NewTabItemWithIcon("Home", theme.HomeIcon(), homeTabContent)
	app.CopiedTradersTab.TabItem = container.NewTabItemWithIcon("Copied Traders", theme.GridIcon(), copiedTradersTabContent)
	app.LeaderboardTab.TabItem = container.NewTabItemWithIcon("Leaderboard", theme.ListIcon(), leaderboardContent)
	app.AnalyticsTab.TabItem = container.NewTabItemWithIcon("Analytics", theme.ComputerIcon(), canvas.NewText("Analytics content goes here", nil))
	app.ProfilesTab.TabItem = container.NewTabItemWithIcon("Profiles", app.App.Settings().Theme().Icon(resources.IconNameCreditCard), profilesTabContent)

	// add application tabs (home, tasks, proxies, profiles, settings)
	tabs := container.NewAppTabs(
		app.HomeTab.TabItem,
		app.CopiedTradersTab.TabItem,
		app.LeaderboardTab.TabItem,
		app.AnalyticsTab.TabItem,
		app.ProfilesTab.TabItem,
	)
	tabs.SetTabLocation(container.TabLocationTop)

	// populate window
	app.TopRightToolbar = app.getToolbar()
	app.GlobalContent = container.NewWithoutLayout(tabs, app.TopRightToolbar)

	// resize and position widgets
	tabs.Resize(fyne.NewSize(1280, 720))
	app.TopRightToolbar.Resize(fyne.NewSize(100, 30))
	tabs.Move(fyne.NewPos(0, 0))
	app.TopRightToolbar.Move(fyne.NewPos(1180, 0))

	app.MainWindow.SetContent(app.GlobalContent)

}

// getAccountIcon returns either a static image or an animated gif for the current user image
func (app *Config) getAccountIcon() (fyne.Resource, *x_widget.AnimatedGif) {

	// if profile image is not stored locally
	if app.User.ProfilePicturePath == "" {
		// see if we have information to download it
		if app.User != nil && app.User.ProfilePictureURL != "" && app.User.DiscordID != "" {
			ext := filepath.Ext(app.User.ProfilePictureURL)
			// if profilePicture url does not contain file extension
			if ext == "" {
				// download the file now
				response, err := app.Client.Get(app.User.ProfilePictureURL)
				if err != nil {
					app.Logger.Error(err)
					return theme.AccountIcon(), nil
				}

				if response.StatusCode != 200 {
					return theme.AccountIcon(), nil
				}

				// store downloaded image locally
				file, err := os.Create(fmt.Sprintf("downloads/%s", app.User.DiscordID))
				if err != nil {
					app.Logger.Error(err)
					return theme.AccountIcon(), nil
				}
				defer file.Close()

				_, err = io.Copy(file, response.Body)
				if err != nil {
					app.Logger.Error(err)
					return theme.AccountIcon(), nil
				}

				// get the file extension
				ext = response.Header.Get("Content-Type")
				switch ext {
				case "image/gif":
					ext = "gif"
				case "image/jpeg", "image/jpg":
					ext = "jpg"
				case "image/png":
					ext = "png"
				default:
					app.Logger.Error(errors.New("no content-type header found when getting profile image"))
					return theme.AccountIcon(), nil
				}

				// finally, add extension to the stored image
				src := fmt.Sprintf("downloads/%s", app.User.DiscordID)
				dst := fmt.Sprintf("downloads/%s.%s", app.User.DiscordID, ext)

				// rename file
				os.Rename(src, dst)
				app.User.ProfilePicturePath = dst
				err = app.DB.UpdateUser(app.User.ID, *app.User)
				if err != nil {
					app.Logger.Error(err)
				}
			} else {
				// download and save to file
				err := app.downloadFile(fmt.Sprintf("downloads/%s%s", app.User.DiscordID, ext), app.User.DiscordID, ext)
				if err != nil {
					app.Logger.Error(err)
					// return default account icon
					return theme.AccountIcon(), nil
				}
				app.User.ProfilePicturePath = fmt.Sprintf("downloads/%s.%s", app.User.DiscordID, ext)
				err = app.DB.UpdateUser(app.User.ID, *app.User)
				if err != nil {
					app.Logger.Error(err)
				}
			}
		} else {
			// if not, return a default account icon
			return theme.AccountIcon(), nil
		}
	}

	ext := filepath.Ext(app.User.ProfilePicturePath)
	switch ext {
	case ".jpg", ".png":
		resource, err := fyne.LoadResourceFromPath(app.User.ProfilePicturePath)
		if err != nil {
			app.Logger.Error(err)
			// return default account icon
			return theme.AccountIcon(), nil
		}
		return resource, nil
	case ".gif":
		gif, err := x_widget.NewAnimatedGif(storage.NewFileURI(app.User.ProfilePicturePath))
		if err != nil {
			app.Logger.Error(err)
			// return default account icon
			return theme.AccountIcon(), nil
		}
		return nil, gif
	default:
		return theme.AccountIcon(), nil
	}
}

func (app *Config) getToolbar() *widget.Toolbar {
	var accIcon fyne.Resource
	accIcon, _ = app.getAccountIcon()

	if accIcon != nil {
		return widget.NewToolbar(
			widget.NewToolbarSpacer(),
			widget.NewToolbarAction(accIcon, func() {
			}),
			widget.NewToolbarSpacer())
	} else {
		return widget.NewToolbar(
			widget.NewToolbarSpacer(),
			// AnimatedGif does not implement fyne.Resource so we cannot use it as toolabr icon :(
			widget.NewToolbarAction(theme.AccountIcon(), func() {
			}),
			widget.NewToolbarSpacer())
	}
}

func (app *Config) refreshGlobalContent() {
	go func() {
		// refresh custom icons for app tabs
		app.ProfilesTab.Icon = app.App.Settings().Theme().Icon(resources.IconNameCreditCard)

		app.TopRightToolbar = app.getToolbar()
		app.GlobalContent.Objects[1] = app.TopRightToolbar

		// resize and position widgets
		app.TopRightToolbar.Resize(fyne.NewSize(100, 30))
		app.TopRightToolbar.Move(fyne.NewPos(1180, 0))

		app.GlobalContent.Refresh()
	}()
}

func (app *Config) MakeMenu() *fyne.MainMenu {

	setDarkThemeItem := fyne.NewMenuItem("Dark", func() {
		// set dark theme and refresh widgets
		app.User.Theme = "dark"
		app.DB.UpdateUser(app.User.ID, *app.User)
		app.App.Settings().SetTheme(&resources.DarkTheme{})
		app.HomeTab.Content.Refresh()
		app.CopiedTradersTab.Content.Refresh()
		app.RefreshLeaderboardWithoutFetch()
		app.LeaderboardTab.Content.Refresh()
		app.AnalyticsTab.Content.Refresh()
		app.ProfilesTab.Content.Refresh()
		app.refreshGlobalContent()
	})
	setLightThemeItem := fyne.NewMenuItem("Light", func() {
		// set light theme and refresh widgets
		app.User.Theme = "light"
		app.DB.UpdateUser(app.User.ID, *app.User)
		app.App.Settings().SetTheme(&resources.LightTheme{})
		app.HomeTab.Content.Refresh()
		app.CopiedTradersTab.Content.Refresh()
		app.RefreshLeaderboardWithoutFetch()
		app.LeaderboardTab.Content.Refresh()
		app.AnalyticsTab.Content.Refresh()
		app.ProfilesTab.Content.Refresh()
		app.refreshGlobalContent()
	})
	themeItem := fyne.NewMenuItem("Theme", nil)
	themeItem.ChildMenu = fyne.NewMenu("",
		setDarkThemeItem,
		setLightThemeItem)

	fileMenu := fyne.NewMenu("Settings", themeItem)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://docs.aqua-io.trading")
			_ = app.App.OpenURL(u)
		}),
		fyne.NewMenuItem("Contacts", func() {
			app.Contacts()
		}),
		fyne.NewMenuItem("Legal Notes", func() {
			u, _ := url.Parse("https://aqua-io.trading/legal")
			_ = app.App.OpenURL(u)
		}),
	)

	main := fyne.NewMainMenu(
		fileMenu,
		helpMenu,
	)

	return main
}

func (app *Config) MakeTray() {
	if desk, ok := app.App.(desktop.App); ok {
		var menu *fyne.Menu
		if app.User != nil {
			h := fyne.NewMenuItem("Logout", func() {})
			h.Icon = theme.LogoutIcon()
			menu = fyne.NewMenu("", h)
			h.Action = func() {
				app.Logout()
			}
		} else {
			menu = fyne.NewMenu("")
		}

		desk.SetSystemTrayIcon(resources.ResourceIconPng)
		desk.SetSystemTrayMenu(menu)
	}
}

// Contacts shows a splash screen for a few seconds
func (app *Config) Contacts() {
	drv := fyne.CurrentApp().Driver()
	if drv, ok := drv.(desktop.Driver); ok {
		w := drv.CreateSplashWindow()
		appLogo := canvas.NewImageFromResource(resources.ResourceIconPng)
		appLogo.SetMinSize(fyne.NewSize(25, 25))
		appLogo.FillMode = canvas.ImageFillContain
		aqua_io, _ := url.Parse("https://aqua-io.trading")
		discord, _ := url.Parse("https://discord.gg/JPw4uuRDQX")
		vBox := container.NewVBox(
			container.NewCenter(container.NewHBox(widget.NewRichTextFromMarkdown(`## Contact Us`), appLogo)),
			widget.NewHyperlinkWithStyle("Website", aqua_io, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewHyperlinkWithStyle("Discord", discord, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithStyle("info@aqua-io.trading", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		)
		w.SetContent(vBox)
		w.Resize(fyne.NewSize(300, 200))
		w.Show()
		go func() {
			time.Sleep(time.Second * 5)
			w.Close()
		}()
	}
}
