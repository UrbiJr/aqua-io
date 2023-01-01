package cactus

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ProfileView contains information about the "Profile" view
type ProfileView struct {
	Title         string
	ProfilesTable *tview.Table
	View          *tview.Flex
}

// NewSitelistView returns a view for the profiles management
func (cactus *Cactus) NewProfilesView() *ProfileView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	table := tview.NewTable().
		SetBorders(true).
		SetBordersColor(tcell.ColorGrey)

	footerForm := tview.NewForm()
	footerForm.AddButton("Create New", func() {
		cactus.pages.SwitchToPage(cactus.NewProfileView.Title)
	})
	footerForm.AddButton("Go Back", func() {
		cactus.UI.OnGoBackSelected()
	})

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("Profiles"), 0, 1, false).
		AddItem(table, 0, 4, false).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(footerForm, 0, 4, true), 0, 1, true).
		SetBorder(true)

	return &ProfileView{Title: "Profiles", View: flex, ProfilesTable: table}
}

// RefreshProfileView refresh the list of profiles
func (cactus *Cactus) RefreshProfileView() {
	cactus.ProfilesView.ProfilesTable.Clear()

	for i, profile := range cactus.User.Profiles {
		// table cell containing profile name
		cactus.ProfilesView.ProfilesTable.SetCell(i, 0, tview.NewTableCell(profile.Title).
			SetAlign(tview.AlignCenter))

		// table cell containing "Rename button"
		cactus.ProfilesView.ProfilesTable.SetCell(i, 1, tview.NewTableCell("Rename").
			SetClickedFunc(
				func() bool {
					cactus.RenameProfileForm(profile)
					cactus.pages.SwitchToPage(cactus.EditProfileView.Title)
					return true
				},
			).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter))

		// table cell containing "Edit button"
		cactus.ProfilesView.ProfilesTable.SetCell(i, 2, tview.NewTableCell("Edit").
			SetClickedFunc(
				func() bool {
					cactus.EditProfileForm(profile)
					cactus.pages.SwitchToPage(cactus.EditProfileView.Title)
					return true
				},
			).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter))

		// table cell containing "Delete button"
		cactus.ProfilesView.ProfilesTable.SetCell(i, 3, tview.NewTableCell("Delete").
			SetClickedFunc(
				func() bool {
					cactus.ShowConfirm(
						func() {
							cactus.DeleteProfile(profile.Title)
							cactus.RefreshProfileView()
							cactus.pages.SwitchToPage(cactus.ProfilesView.Title)
						},
						func() {
							cactus.pages.SwitchToPage(cactus.ProfilesView.Title)
						},
						fmt.Sprintf("Are you sure you want to delete profile \"%s\"?", profile.Title),
					)
					return true
				},
			).
			SetTextColor(tcell.ColorRed).
			SetAlign(tview.AlignCenter))
	}

}
