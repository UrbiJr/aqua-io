package copy_io

import (
	"fmt"
	"math/rand"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/copy-io/internal/resources"
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
		"how can Copy.io assist you today? :-)",
		"how are you going to use Copy.io today? :-)",
		"ready to have some fun with Copy.io?",
		"Copy.io is at your service :salute:",
		"it's been a while.",
		"time to make the checkout feed go brrr.",
		"everyday is a perfect day to run Copy.io!",
	}
	msg := fmt.Sprintf("Hello %s, %s", app.User.Username, greetings[rand.Intn(len(greetings))])

	// get app tabs content
	app.HomeTab = &HomeTab{}
	app.CopiedTradersTab = &CopiedTradersTab{}
	app.LeaderboardTab = &LeaderboardTab{}
	app.AnalyticsTab = &AnalyticsTab{}
	app.ProfilesTab = &ProfilesTab{}

	profilesTabContent := app.profilesTab() // important: load profiles tab content before anyone else!
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
	tabs.Resize(fyne.NewSize(1390, 848))
	app.TopRightToolbar.Resize(fyne.NewSize(100, 30))
	tabs.Move(fyne.NewPos(0, 0))
	app.TopRightToolbar.Move(fyne.NewPos(1290, 0))

	app.MainWindow.SetContent(app.GlobalContent)

}

func (app *Config) getToolbar() *widget.Toolbar {
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.AccountIcon(), func() {
		}),
		widget.NewToolbarAction(app.App.Settings().Theme().Icon(resources.IconNameBell), func() {
		}),
		widget.NewToolbarSpacer(),
	)
	return toolbar
}

func (app *Config) refreshGlobalContent() {
	// refresh custom icons for app tabs
	app.ProfilesTab.Icon = app.App.Settings().Theme().Icon(resources.IconNameCreditCard)

	app.TopRightToolbar = app.getToolbar()
	app.GlobalContent.Objects[1] = app.TopRightToolbar

	// resize and position widgets
	app.TopRightToolbar.Resize(fyne.NewSize(100, 30))
	app.TopRightToolbar.Move(fyne.NewPos(1180, 0))

	app.GlobalContent.Refresh()
}

func (app *Config) MakeMenu() *fyne.MainMenu {

	setDarkThemeItem := fyne.NewMenuItem("Dark", func() {
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
			u, _ := url.Parse("https://docs.copy-trading.io")
			_ = app.App.OpenURL(u)
		}))

	main := fyne.NewMainMenu(
		fileMenu,
		helpMenu,
	)

	return main
}
