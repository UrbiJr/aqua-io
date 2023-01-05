package nyx

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ProxiesView contains information about the "Profile" view
type ProxiesView struct {
	Title string
	View  *tview.Flex
}

// NewTaskCreationView returns a view for the profiles management
func (nyx *Nyx) NewProxiesView() *ProxiesView {

	var flex = tview.NewFlex() // Flexbox layout allows us to organize multiple widgets inside a view

	return &ProxiesView{
		Title: "Proxies",
		View:  flex,
	}
}

// RefreshProxiesView refreshes this view with updated proxy profiles and the current selection, if any
func (nyx *Nyx) RefreshProxiesView(proxyProfileIdx int) {

	leftContainer := tview.NewFlex()
	proxyProfilesTable := tview.NewTable()
	headerContainer := tview.NewFlex()
	proxiesTable := tview.NewTable()

	proxyProfilesTable.SetCell(0, 0, tview.NewTableCell("➕ Create New").
		//SetExpansion(2).
		SetClickedFunc(func() bool {
			nyx.pages.SwitchToPage(nyx.NewProxyProfileView.Title)
			return true
		}).
		SetAlign(tview.AlignLeft))

	for i, proxyProfile := range nyx.User.ProxyProfiles {
		// table cell containing proxyProfile name
		proxyProfilesTable.SetCell(i+1, 0, tview.NewTableCell(fmt.Sprintf("%s\t(%d)", proxyProfile.Name, len(proxyProfile.Proxies))).
			//SetExpansion(2).
			SetClickedFunc(func() bool {
				nyx.RefreshProxiesView(i)
				return true
			}).
			SetAlign(tview.AlignLeft))
	}

	proxyProfilesTable.SetTitle("📁 Proxy Profiles")
	navigationForm := tview.NewForm()
	navigationForm.AddButton("Go Back", func() {
		nyx.UI.OnGoBackSelected()
	})
	leftContainer.SetDirection(tview.FlexRow).
		AddItem(proxyProfilesTable, 0, 4, false).
		AddItem(navigationForm, 0, 1, false).
		SetBorder(true).
		SetBorderPadding(1, 1, 0, 0)

	if proxyProfileIdx != -1 {
		for i, proxy := range nyx.User.ProxyProfiles[proxyProfileIdx].Proxies {
			// table cell containing proxyProfile name
			proxiesTable.SetCell(i, 0, tview.NewTableCell(proxy).
				SetExpansion(2).
				SetAlign(tview.AlignLeft))

			// table cell containing "Delete button"
			proxiesTable.SetCell(i, 1, tview.NewTableCell("🗑️").
				SetClickedFunc(
					func() bool {
						return true
					},
				).
				SetTextColor(tcell.ColorRed).
				SetAlign(tview.AlignLeft))
		}

		headerContainer.
			AddItem(tview.NewTextView().SetText(fmt.Sprintf("%s 🌐", nyx.User.ProxyProfiles[proxyProfileIdx].Name)), 1, 1, false).
			AddItem(
				tview.NewForm().
					AddButton("Add Proxies", func() {
					}).
					AddButton("Clear All", func() {
					}).
					AddButton("Delete", func() {
					}), 0, 2, false)

		nyx.ProxiesView.View.
			AddItem(leftContainer, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(headerContainer, 1, 1, false).
				AddItem(proxiesTable, 0, 4, false), 0, 3, false)
	} else {
		nyx.ProxiesView.View.
			AddItem(leftContainer, 0, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexRow).SetBorder(true), 0, 3, false)
	}

}
