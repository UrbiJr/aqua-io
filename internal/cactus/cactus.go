package cactus

import (
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/cactus-aio/go-cactus/internal/user"
	"github.com/rivo/tview"
)

// Cactus is the container of the main app, it contains the main attributes
type Cactus struct {
	*UI
	SiteList []*SupportedSite
	*log.Logger
	User *user.User
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

// SaveProfile appends profile to user profiles list and writes the updated list to file
func (cactus *Cactus) SaveProfile(profile user.Profile) error {
	for _, p := range cactus.User.Profiles {
		if p.Title == profile.Title {
			return errors.New("a profile with this title is already existent")
		}
	}
	cactus.User.Profiles = append(cactus.User.Profiles, profile)
	user.WriteProfiles(cactus.User.Profiles)
	return nil
}

// DeleteProfile removes profile from user profiles list and writes the updated list to file
func (cactus *Cactus) DeleteProfile(profileTitle string) error {
	for i, p := range cactus.User.Profiles {
		if p.Title == profileTitle {
			// remove from slice and preserve order
			cactus.User.Profiles = append(cactus.User.Profiles[:i], cactus.User.Profiles[i+1:]...)
			user.WriteProfiles(cactus.User.Profiles)
			return nil
		}
	}
	return errors.New("cannot find a profile with this title")
}

// Run initializes the app along with its layout
func (cactus *Cactus) Run() error {

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
