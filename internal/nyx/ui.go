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
func (nyx *Config) MakeMobileUI() {

	// add application tabs (home, tasks, proxies, billing, settings)
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("          ", theme.HomeIcon(), canvas.NewText("Home content goes here", nil)),
		container.NewTabItemWithIcon("          ", theme.ListIcon(), canvas.NewText("Tasks content goes here", nil)),
		container.NewTabItemWithIcon("          ", nyx.App.Settings().Theme().Icon(resources.IconNameWifi), canvas.NewText("Proxies content goes here", nil)),
		container.NewTabItemWithIcon("          ", nyx.App.Settings().Theme().Icon(resources.IconNameCreditCard), canvas.NewText("Billing content goes here", nil)),
		container.NewTabItemWithIcon("          ", theme.SettingsIcon(), canvas.NewText("Settings content goes here", nil)),
	)
	// show tabs at the bottom of the window
	tabs.SetTabLocation(container.TabLocationBottom)

	// populate app tabs

	// populate window
	finalContent := container.NewVBox(tabs)

	nyx.MainWindow.SetContent(finalContent)

}

// MakeMobileUI instantiates all the needed pages and makes the UI layout, but does not display it yet.
func (nyx *Config) MakeDesktopUI() {

	// add application tabs (home, tasks, proxies, billing, settings)
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), canvas.NewText("Home content goes here", nil)),
		container.NewTabItemWithIcon("Tasks", theme.ListIcon(), canvas.NewText("Tasks content goes here", nil)),
		container.NewTabItemWithIcon("Proxies", nyx.App.Settings().Theme().Icon(resources.IconNameWifi), canvas.NewText("Proxies content goes here", nil)),
		container.NewTabItemWithIcon("Billing", nyx.App.Settings().Theme().Icon(resources.IconNameCreditCard), canvas.NewText("Billing content goes here", nil)),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), canvas.NewText("Settings content goes here", nil)),
	)
	tabs.SetTabLocation(container.TabLocationTop)
	tabs.MinSize()

	// populate app tabs

	// populate window
	toolbar := nyx.getToolbar()
	finalContent := container.NewWithoutLayout(tabs, toolbar)

	// resize and position widgets
	tabs.Resize(fyne.NewSize(1200, 600))
	toolbar.Resize(fyne.NewSize(100, 30))
	tabs.Move(fyne.NewPos(0, 0))
	toolbar.Move(fyne.NewPos(1100, 0))

	nyx.MainWindow.SetContent(finalContent)

}

func (nyx *Config) getToolbar() *widget.Toolbar {
	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.AccountIcon(), func() {
		}),
		widget.NewToolbarAction(nyx.App.Settings().Theme().Icon(resources.IconNameBell), func() {
		}),
		widget.NewToolbarSpacer(),
	)
	return toolbar
}
