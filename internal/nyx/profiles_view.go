package nyx

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

// NewTaskCreationView returns a view for the profiles management
func (nyx *Nyx) NewProfilesView() *ProfileView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view
	table := tview.NewTable().
		SetBorders(true).
		SetBordersColor(tcell.ColorGrey)

	footerForm := tview.NewForm()
	footerForm.AddButton("Create New", func() {
		nyx.pages.SwitchToPage(nyx.NewProfileView.Title)
	})
	footerForm.AddButton("Go Back", func() {
		nyx.UI.OnGoBackSelected()
	})

	flex.SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorGreen).SetText("👥 Profiles"), 0, 1, false).
		AddItem(table, 0, 4, false).
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(footerForm, 0, 4, true), 0, 1, true).
		SetBorder(true)

	return &ProfileView{Title: "Profiles", View: flex, ProfilesTable: table}
}

// RefreshProfileView refresh the list of profiles
func (nyx *Nyx) RefreshProfileView() {
	nyx.ProfilesView.ProfilesTable.Clear()

	for i, profile := range nyx.User.Profiles {
		// table cell containing profile name
		nyx.ProfilesView.ProfilesTable.SetCell(i, 0, tview.NewTableCell(profile.Title).
			SetExpansion(2).
			SetAlign(tview.AlignCenter))

		// table cell containing "Rename button"
		nyx.ProfilesView.ProfilesTable.SetCell(i, 1, tview.NewTableCell("✏️ Rename").
			SetExpansion(2).
			SetClickedFunc(
				func() bool {
					nyx.RenameProfileForm(profile)
					nyx.pages.SwitchToPage(nyx.EditProfileView.Title)
					return true
				},
			).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter))

		// table cell containing "Edit button"
		nyx.ProfilesView.ProfilesTable.SetCell(i, 2, tview.NewTableCell("📝 Edit").
			SetExpansion(2).
			SetClickedFunc(
				func() bool {
					nyx.EditProfileForm(profile)
					nyx.pages.SwitchToPage(nyx.EditProfileView.Title)
					return true
				},
			).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignCenter))

		// table cell containing "Delete button"
		nyx.ProfilesView.ProfilesTable.SetCell(i, 3, tview.NewTableCell("❌ Delete").
			SetExpansion(2).
			SetClickedFunc(
				func() bool {
					nyx.ShowConfirm(
						func() {
							nyx.DeleteProfile(profile.Title)
							nyx.RefreshProfileView()
							nyx.pages.SwitchToPage(nyx.ProfilesView.Title)
						},
						func() {
							nyx.pages.SwitchToPage(nyx.ProfilesView.Title)
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
