package cactus

import (
	"fmt"
	"math/rand"

	"github.com/cactus-aio/go-cactus/internal/user"
	"github.com/rivo/tview"
)

// Cactus is the container of the main app, it contains the main attributes
type Cactus struct {
	*UI
	SiteList []*SupportedSite
}

// NewCactus returns a new instance of the app
func NewCactus() *Cactus {
	sitelist := []*SupportedSite{
		{Name: Nike, Category: SneakerSite, CSVFields: []string{"PID", "MIN SIZE", "MAX SIZE", "PROFILE", "MODE", "REGION"}},
		{Name: Schuh, Category: SneakerSite, CSVFields: []string{"PID", "MIN SIZE", "MAX SIZE", "PROFILE", "MODE", "REGION"}},
	}

	cactus := &Cactus{
		UI:       NewUI(),
		SiteList: sitelist,
	}

	return cactus
}

// Quit exits the app gracefully
func (cactus *Cactus) Quit() {
	cactus.UI.tui.Stop()
}

// Run initializes the app along with its layout
func (cactus *Cactus) Run(user *user.User) error {

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

	msg := fmt.Sprintf("Hello %s, %s", user.Username, greetings[rand.Intn(len(greetings))])
	entries := []MenuEntry{
		{name: "Sitelist", label: '1', description: "display Cactus-AIO sitelist", selected: nil},
		{name: "Profiles", label: '2', description: "manage your profiles", selected: nil},
		{name: "Settings", label: '3', description: "edit Cactus-AIO settings", selected: nil},
		{name: "Quit", label: '4', description: "close Cactus-AIO", selected: cactus.Quit},
	}
	menu := cactus.UI.NewMainMenu(msg, entries) // create main view Menu
	var pages = tview.NewPages()                // Allows us to easily switch between views
	pages.AddPage("Main Menu", menu, true, true)

	// Enable mouse detection
	// The SetRoot function tells the tview app which widget to display when the application starts
	if err := cactus.UI.tui.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		return err
	}

	return nil
}
