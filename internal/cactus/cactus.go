package cactus

import (
	"fmt"
	"math/rand"

	"github.com/cactus-aio/go-cactus/internal/user"
	"github.com/gdamore/tcell/v2"
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
	list := tview.NewList().
		AddItem("Sitelist", "display Cactus-AIO sitelist", '1', nil).
		AddItem("Profiles", "manage your profiles", '2', nil).
		AddItem("Settings", "edit Cactus-AIO settings", '3', nil).
		AddItem("Quit", "close Cactus-AIO", '4', cactus.Quit)

	var pages = tview.NewPages() // Allows us to easily switch between views
	var flex = tview.NewFlex()   // Flexbox layout allows us to organize multiple widgets inside a view

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText(msg), 0, 1, false).
		AddItem(list, 0, 4, true).SetBorder(true)

	pages.AddPage("Main Menu", flex, true, true)

	// Enable mouse detection
	// The SetRoot function tells the tview app which widget to display when the application starts
	if err := cactus.UI.tui.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		return err
	}

	return nil
}
