package nyx

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/nyx/internal/user"
	"github.com/UrbiJr/nyx/internal/utils"
)

// NewProfilesView returns a view for the profiles management
func (app *Config) billingTab() *fyne.Container {

	// get current profile groups
	app.User.ProfileManager.Groups, _ = app.DB.AllProfileGroups()

	// define a list to display profile groups
	list := widget.NewList(
		func() int {
			return len(app.User.ProfileManager.Groups)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if label, ok := item.(*widget.Label); ok {
				label.SetText(
					fmt.Sprintf("%s\t (%d)", app.User.ProfileManager.Groups[id].Name, len(app.User.ProfileManager.FilterByGroupName(app.User.ProfileManager.Groups[id].Name))))
			}
		},
	)
	app.ProfileGroupsList = list

	// define a button to create a new group
	addProfileGroupButton := widget.NewButtonWithIcon("New Group", theme.ContentAddIcon(), func() {
		app.addProfileGroupDialog()
	})

	// define the left container
	leftContainer := container.NewVBox(canvas.NewText("PROFILE GROUPS", nil), addProfileGroupButton, list)

	// create main table to list profiles

	// display current group title on top of profiles table

	// define the center container

	// define the billingTab container
	profilesTabContainer := container.NewWithoutLayout(leftContainer)

	// resize and move billingTab elements
	leftContainer.Move(fyne.NewPos(10, 10))
	leftContainer.Resize(fyne.NewSize(280, 600))

	return profilesTabContainer
}

func (app *Config) addProfileGroupDialog() dialog.Dialog {
	nameEntry := widget.NewEntry()

	nameEntry.Validator = utils.IsStringEmpty

	// create a dialog
	addForm := dialog.NewForm(
		"Create Group",
		"Create",
		"Cancel",
		[]*widget.FormItem{
			{Text: "Group name", Widget: nameEntry},
		},
		func(valid bool) {
			if valid {

				_, err := app.DB.InsertProfileGroup(user.ProfileGroup{
					Name: nameEntry.Text,
				})

				if err != nil {
					app.Logger.Error(err)
				}
				app.refreshProfileGroupsList()
			}
		},
		app.MainWindow)

	// size and show the dialog
	addForm.Resize(fyne.Size{Width: 400})
	addForm.Show()

	return addForm
}

func (app *Config) currentProfileGroups() ([]user.ProfileGroup, error) {
	pfg, err := app.DB.AllProfileGroups()
	if err != nil {
		app.Logger.Error(err)
		return nil, err
	}

	return pfg, nil
}

func (app *Config) currentProfiles() ([]user.Profile, error) {
	p, err := app.DB.AllProfiles()
	if err != nil {
		app.Logger.Error(err)
		return nil, err
	}

	return p, nil
}

func (app *Config) refreshProfileGroupsList() {
	app.User.ProfileManager.Groups, _ = app.DB.AllProfileGroups()
	app.ProfileGroupsList.Refresh()
}
