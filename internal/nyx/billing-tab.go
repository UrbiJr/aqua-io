package nyx

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// NewProfilesView returns a view for the profiles management
func (app *Config) billingTab() *fyne.Container {

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

	// define a button to create a new group
	addProfileGroupButton := widget.NewButtonWithIcon("New Group", theme.ContentAddIcon(), func() {
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

// RefreshProfilesView refresh the list of profiles
func (app *Config) refreshProfilesTab() {
	/*
		for i, profile := range app.User.Profiles {
			// table cell containing profile name

		}
	*/
}
