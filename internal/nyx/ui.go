package nyx

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI contains information about the user interface
type UI struct {
	tui                *tview.Application
	pages              *tview.Pages
	MainMenuView       *MainMenuView
	TaskCreationView   *TaskCreationView
	ProfilesView       *ProfileView
	StateSelectionView *StateSelectionView
	NewProfileView     *NewProfileView
	EditProfileView    *EditProfileView
	SettingsView       *SettingsView
	ErrorView          *ErrorView
	ConfirmView        *ConfirmView
}

// OnGoBackSelected should be called when a user choose to go to the previous view
func (ui *UI) OnGoBackSelected() {

	// switch selection based on current page
	switch currentPageTitle, _ := ui.pages.GetFrontPage(); currentPageTitle {
	case "Task Creation", "Profiles":
		// if current page is Task Creation, go back to main menu
		ui.pages.SwitchToPage(ui.MainMenuView.Title)
	}

}

// OnProfilesSelected should be called when a user choose Profiles entry on main menu
func (ui *UI) OnProfilesSelected() {

	// switch current view to Task Creation
	ui.pages.SwitchToPage(ui.ProfilesView.Title)
}

// OnSettingsSelected should be called when a user choose Settings entry on main menu
func (nyx *Nyx) OnSettingsSelected() {

	// switch current view to Task Creation
	nyx.EditSettingsForm()
	nyx.UI.pages.SwitchToPage(nyx.SettingsView.Title)
}

// OnTaskCreationSelected should be called when a user choose Task Creation entry on main menu
func (ui *UI) OnTaskCreationSelected() {

	// switch current view to Task Creation
	ui.pages.SwitchToPage(ui.TaskCreationView.Title)
}

// NewUI initializes a tview and assigns it to the returned UI object
func NewUI() *UI {
	ui := new(UI)
	ui.tui = tview.NewApplication()

	return ui
}

// DrawUI instantiates all the needed pages and makes the UI layout, but does not display it. For that, call Run()
func (nyx *Nyx) DrawUI() {

	nyx.UI.pages = tview.NewPages() // Allows us to easily switch between views

	// Listeners must be initialized before anything else
	nyx.SetListeners()

	entries := []MenuEntry{
		{name: "Task Creation", label: '1', description: "create new tasks", selected: nyx.OnTaskCreationSelected},
		{name: "Profiles", label: '2', description: "manage your profiles", selected: nyx.OnProfilesSelected},
		{name: "Settings", label: '3', description: "edit Nyx settings", selected: nyx.OnSettingsSelected},
		{name: "Quit", label: 'q', description: "close Nyx", selected: nyx.Quit},
	}

	nyx.ErrorView = nyx.NewErrorView()
	nyx.ConfirmView = nyx.NewConfirmView()
	nyx.MainMenuView = nyx.NewMainMenuView(entries) // create main view Menu
	nyx.TaskCreationView = nyx.NewTaskCreationView()
	nyx.ProfilesView = nyx.NewProfilesView()
	nyx.NewProfileView = nyx.NewNewProfileView()
	nyx.EditProfileView = nyx.NewEditProfileView()
	nyx.SettingsView = nyx.NewSettingsView()
	nyx.StateSelectionView = nyx.NewStateSelectionView()
	nyx.AddProfileForm()

	nyx.UI.pages.AddPage(nyx.ErrorView.Title, nyx.ErrorView.View, true, false)
	nyx.UI.pages.AddPage(nyx.ConfirmView.Title, nyx.ConfirmView.View, true, false)
	nyx.UI.pages.AddPage(nyx.MainMenuView.Title, nyx.MainMenuView.View, true, true)
	nyx.UI.pages.AddPage(nyx.TaskCreationView.Title, nyx.TaskCreationView.View, true, false)
	nyx.UI.pages.AddPage(nyx.ProfilesView.Title, nyx.ProfilesView.View, true, false)
	nyx.UI.pages.AddPage(nyx.NewProfileView.Title, nyx.NewProfileView.View, true, false)
	nyx.UI.pages.AddPage(nyx.StateSelectionView.Title, nyx.StateSelectionView.View, true, false)
	nyx.UI.pages.AddPage(nyx.EditProfileView.Title, nyx.EditProfileView.View, true, false)
	nyx.UI.pages.AddPage(nyx.SettingsView.Title, nyx.SettingsView.View, true, false)
}

/*
Run displays the UI. This method should be called as last, as it's blocking.
In particular, the Run function of tview.Application starts the application and takes care of handling user interface events, such as rendering and handling user input. Thus, once the application is started with Run, subsequent code will not execute until the application is stopped, such as by calling the Stop function or exiting the application window.

If you want to run other code after starting the application, you'll need to do it asynchronously, such as using a goroutine
*/
func (ui *UI) Run() error {

	// Enable mouse detection
	// The SetRoot function tells the tview app which widget to display when the application starts
	if err := ui.tui.SetRoot(ui.pages, true).EnableMouse(true).Run(); err != nil {
		return err
	}
	return nil
}

func (nyx *Nyx) SetListeners() {
	nyx.tui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		// user pressed CTRL-C
		if event.Key() == tcell.KeyCtrlC {
			nyx.Quit()
		}

		// switch selection based on current page
		switch currentPageTitle, _ := nyx.pages.GetFrontPage(); currentPageTitle {
		case "Profiles":
			switch pressedKey := event.Rune(); pressedKey {
			case 97: // user presses 'a' key
			case 100: // user presses 'd' key
			case 114: // user presses 'r' key
			}
		}
		return event
	})

}
