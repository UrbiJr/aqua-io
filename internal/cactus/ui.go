package cactus

import (
	"fmt"
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI contains information about the user interface
type UI struct {
	tui             *tview.Application
	pages           *tview.Pages
	MainMenuView    *tview.Flex
	SitelistView    *tview.Flex
	ProfilesView    *ProfileView
	SelectStateView *tview.Flex
	NewProfileView  *NewProfileView
}

// OnGoBackSelected should be called when a user choose to go to the previous view
func (ui *UI) OnGoBackSelected() {

	// switch selection based on current page
	switch currentPageTitle, _ := ui.pages.GetFrontPage(); currentPageTitle {
	case "Sitelist", "Profiles":
		// if current page is Sitelist, go back to main menu
		ui.pages.SwitchToPage("Main Menu")
	}

}

// OnProfilesSelected should be called when a user choose Profiles entry on main menu
func (ui *UI) OnProfilesSelected() {

	// switch current view to Sitelist
	ui.pages.SwitchToPage("Profiles")
}

// OnSitelistSelected should be called when a user choose Sitelist entry on main menu
func (ui *UI) OnSitelistSelected() {

	// switch current view to Sitelist
	ui.pages.SwitchToPage("Sitelist")
}

// NewUI initializes a tview and assigns it to the returned UI object
func NewUI() *UI {
	ui := new(UI)
	ui.tui = tview.NewApplication()

	return ui
}

func (cactus *Cactus) InitUI() error {

	cactus.UI.pages = tview.NewPages() // Allows us to easily switch between views

	// Listeners must be initialized before anything else
	cactus.UI.SetListeners()

	greetings := []string{
		"how can Cactus-AIO assist you today? :-)",
		"how are you going to use Cactus-AIO today? :-)",
		"ready to have some fun with Cactus-AIO?",
		"Cactus-AIO is at your service :cactus-salute:",
		"it's been a while.",
		"time to make the success log go brrr.",
		"everyday is a perfect day to run Cactus-AIO!",
	}

	msg := fmt.Sprintf("Hello %s, %s", cactus.User.Username, greetings[rand.Intn(len(greetings))])
	entries := []MenuEntry{
		{name: "Sitelist", label: '1', description: "display Cactus-AIO sitelist", selected: cactus.OnSitelistSelected},
		{name: "Profiles", label: '2', description: "manage your profiles", selected: cactus.OnProfilesSelected},
		{name: "Settings", label: '3', description: "edit Cactus-AIO settings", selected: nil},
		{name: "Quit", label: 'q', description: "close Cactus-AIO", selected: cactus.Quit},
	}
	cactus.MainMenuView = cactus.NewMainMenuView(msg, entries) // create main view Menu
	cactus.SitelistView = cactus.NewSitelistView()
	cactus.ProfilesView = cactus.NewProfilesView()
	cactus.NewProfileView = cactus.NewNewProfileView()
	cactus.AddProfileForm()

	cactus.UI.pages.AddPage("Main Menu", cactus.MainMenuView, true, true)
	cactus.UI.pages.AddPage("Sitelist", cactus.SitelistView, true, false)
	cactus.UI.pages.AddPage("Profiles", cactus.ProfilesView.View, true, false)
	cactus.UI.pages.AddPage("New Profile", cactus.NewProfileView.View, true, false)

	// Enable mouse detection
	// The SetRoot function tells the tview app which widget to display when the application starts
	if err := cactus.UI.tui.SetRoot(cactus.UI.pages, true).EnableMouse(true).Run(); err != nil {
		return err
	}

	return nil
}

func (ui *UI) SetListeners() {
	ui.tui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// switch selection based on current page
		switch currentPageTitle, _ := ui.pages.GetFrontPage(); currentPageTitle {
		case "Profiles":
			switch pressedKey := event.Rune(); pressedKey {
			case 97: // user presses 'a' key
				// switch current view to Sitelist
				ui.pages.SwitchToPage("New Profile")
			case 100: // user presses 'd' key
			case 114: // user presses 'r' key
			}
		}
		return event
	})

}
