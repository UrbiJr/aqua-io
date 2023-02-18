package nyx

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/nyx/internal/resources"
)

// UI contains fyne elements
type UI struct {
}

// MakeMobileUI instantiates all the needed pages and makes the UI layout, but does not display it yet.
func (app *Config) MakeMobileUI() {

	// add application tabs (home, tasks, proxies, billing, settings)
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("          ", theme.HomeIcon(), canvas.NewText("Home content goes here", nil)),
		container.NewTabItemWithIcon("          ", theme.ListIcon(), canvas.NewText("Tasks content goes here", nil)),
		container.NewTabItemWithIcon("          ", app.App.Settings().Theme().Icon(resources.IconNameWifi), canvas.NewText("Proxies content goes here", nil)),
		container.NewTabItemWithIcon("          ", app.App.Settings().Theme().Icon(resources.IconNameCreditCard), canvas.NewText("Billing content goes here", nil)),
		container.NewTabItemWithIcon("          ", theme.SettingsIcon(), canvas.NewText("Settings content goes here", nil)),
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

	// get app tabs content
	homeTabContent := app.homeTab()

	// add application tabs (home, tasks, proxies, billing, settings)
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), homeTabContent),
		container.NewTabItemWithIcon("Tasks", theme.ListIcon(), canvas.NewText("Tasks content goes here", nil)),
		container.NewTabItemWithIcon("Proxies", app.App.Settings().Theme().Icon(resources.IconNameWifi), canvas.NewText("Proxies content goes here", nil)),
		container.NewTabItemWithIcon("Billing", app.App.Settings().Theme().Icon(resources.IconNameCreditCard), canvas.NewText("Billing content goes here", nil)),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), canvas.NewText("Settings content goes here", nil)),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	// populate window
	toolbar := app.getToolbar()
	finalContent := container.NewWithoutLayout(tabs, toolbar)

	// resize and position widgets
	tabs.Resize(fyne.NewSize(1280, 720))
	toolbar.Resize(fyne.NewSize(100, 30))
	tabs.Move(fyne.NewPos(0, 0))
	toolbar.Move(fyne.NewPos(1180, 0))

	app.MainWindow.SetContent(finalContent)

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
