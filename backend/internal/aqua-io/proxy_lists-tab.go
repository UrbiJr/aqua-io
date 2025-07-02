package aqua_io

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/UrbiJr/aqua-io/backend/internal/user"
	"github.com/UrbiJr/aqua-io/backend/internal/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type ProxyListsTab struct {
	*container.TabItem
	ProxyListsTop    *fyne.Container
	ProxyListsTable  *widget.Table
	ProxyListsSlice  [][]any
	ProxyListsBottom *fyne.Container
}

// NewProfilesView returns a view for the profiles management
func (app *Config) proxyListsTab() *fyne.Container {
	// define profilesTab
	app.ProxyListsTab.ProxyListsTop = container.NewStack()
	app.ProxyListsTab.ProxyListsBottom = container.NewHBox()

	app.ProxyListsSlice = app.getProxyListsSlice()

	// update content
	app.refreshProxyListsTopContent()
	app.refreshProxyListsBottomContent()

	// get the profiles table
	app.ProxyListsTable = app.getProxyListsTable()

	// define the center container
	vScroll := container.NewScroll(app.ProxyListsTable)

	// define the profilesTab container
	ctr := container.NewWithoutLayout(app.ProxyListsTop, vScroll, app.ProxyListsBottom)

	app.ProxyListsTop.Move(fyne.NewPos(10, 10))
	app.ProxyListsTop.Resize(fyne.NewSize(1280, 40))

	vScroll.Move(fyne.NewPos(10, 60))
	vScroll.Resize(fyne.NewSize(1280, 500))

	app.ProxyListsBottom.Move(fyne.NewPos(300, 580))
	app.ProxyListsBottom.Resize(fyne.NewSize(900, 50))

	return ctr
}

func (app *Config) addProxyListDialog() dialog.Dialog {

	titleLabel := widget.NewLabel("Title")
	titleLabel.Move(fyne.NewPos(10, 10))

	title := widget.NewEntry()
	title.SetPlaceHolder("My Proxy list")
	title.Resize(fyne.NewSize(400, 40))
	title.Move(fyne.NewPos(10, 50))
	title.Validator = func(s string) error {
		s = strings.TrimSpace(s)
		p, _ := app.DB.GetProxyListByTitle(s)
		if p != nil {
			return fmt.Errorf("A proxy list with title \"%s\" already exists", s)
		} else if len(s) <= 0 {
			return errors.New("please insert a title")
		} else if strings.Contains(s, "|") || strings.Contains(s, "\n") {
			return errors.New("please remove invalid characters")
		} else {
			return nil
		}
	}

	proxiesLabel := widget.NewLabel("Proxies (ip:port:user:pass)")
	proxiesLabel.Move(fyne.NewPos(10, 100))

	proxies := widget.NewMultiLineEntry()
	proxies.Resize(fyne.NewSize(400, 300))
	proxies.Move(fyne.NewPos(10, 140))
	proxies.SetPlaceHolder("IP:PORT:USERNAME:PASSWORD")
	proxies.Validator = utils.IsStringEmpty

	vBox := container.NewWithoutLayout(
		titleLabel,
		title,
		proxiesLabel,
		proxies,
	)
	scrollContent := container.NewVScroll(vBox)

	addForm := dialog.NewCustomConfirm(
		"Add New Proxy List",
		"Create",
		"Cancel",
		scrollContent,
		func(valid bool) {

			for _, o := range vBox.Objects {
				switch o := o.(type) {
				case *widget.Entry:
					err := o.Validate()
					if err != nil {
						valid = false
						break
					}
				}
			}

			if valid {
				// split proxies by new line character
				proxiesSlice := strings.Split(strings.ReplaceAll(proxies.Text, "\r\n", "\n"), "\n")

				p := user.ProxyList{
					Title:   title.Text,
					Proxies: proxiesSlice,
				}

				_, err := app.DB.InsertProxyList(p)

				if err != nil {
					app.Logger.Error(err)
				}
				app.refreshProxyListsTab()
				app.refreshProfilesTab()
			}
		},
		app.MainWindow,
	)

	// size and show the dialog
	addForm.Resize(fyne.NewSize(500, 600))
	addForm.Show()

	return addForm
}

func (app *Config) editProxyListDialog(pf *user.ProxyList) dialog.Dialog {
	title := widget.NewEntry()
	title.SetText(pf.Title)
	title.Validator = func(s string) error {
		s = strings.TrimSpace(s)
		p, _ := app.DB.GetProxyListByTitle(s)
		if s != pf.Title && p != nil {
			return fmt.Errorf("A proxy list with title \"%s\" already exists", s)
		} else if len(s) <= 0 {
			return errors.New("please insert a title")
		} else if strings.Contains(s, "|") || strings.Contains(s, "\n") {
			return errors.New("please remove invalid characters")
		} else {
			return nil
		}
	}

	proxies := widget.NewMultiLineEntry()
	proxies.SetText(strings.Join(pf.Proxies, "\n"))
	proxies.Validator = utils.IsStringEmpty

	vBox := container.NewVBox(
		widget.NewLabel("Title"),
		title,
		widget.NewLabel("Proxies (ip:port:user:pass)"),
		proxies,
	)
	scrollContent := container.NewVScroll(vBox)

	editForm := dialog.NewCustomConfirm(
		"Edit Proxy List",
		"Update",
		"Cancel",
		scrollContent,
		func(valid bool) {

			for _, o := range vBox.Objects {
				switch o := o.(type) {
				case *widget.Entry:
					err := o.Validate()
					if err != nil {
						valid = false
						break
					}
				}
			}

			if valid {
				// split proxies by new line character
				proxiesSlice := strings.Split(strings.ReplaceAll(proxies.Text, "\r\n", "\n"), "\n")

				app.Logger.Debug(proxiesSlice)

				p := user.ProxyList{
					Title:   title.Text,
					Proxies: proxiesSlice,
				}

				err := app.DB.UpdateProxyList(pf.ID, p)

				if err != nil {
					app.Logger.Error(err)
				}
				app.refreshProxyListsTab()
				app.refreshProfilesTab()
			}
		},
		app.MainWindow,
	)

	// size and show the dialog
	editForm.Resize(fyne.NewSize(500, 600))
	editForm.Show()

	return editForm

}

func (app *Config) getProxyListsSlice() [][]any {
	var slice [][]any
	slice = append(slice, []any{"Title", "Proxies Loaded", "Actions"})
	all, _ := app.DB.AllProxyLists()

	for _, x := range all {
		var currentRow []any

		if len(x.Title) > 20 {
			currentRow = append(currentRow, x.Title[:19]+"...")
		} else {
			currentRow = append(currentRow, x.Title)
		}

		currentRow = append(currentRow, strconv.Itoa(len(x.Proxies)))

		currentRow = append(currentRow, x.ID) // for toolbar actions

		slice = append(slice, currentRow)
	}

	return slice
}

func (app *Config) getProxyListsTable() *widget.Table {

	t := widget.NewTable(
		func() (int, int) {
			return len(app.ProxyListsSlice), len(app.ProxyListsSlice[0])
		},
		func() fyne.CanvasObject {
			lbl := widget.NewLabel("")
			toolbar := widget.NewToolbar()
			toolbar.Hide()
			return container.NewStack(lbl, toolbar)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			container := o.(*fyne.Container)
			lbl := container.Objects[0].(*widget.Label)
			toolbar := container.Objects[1].(*widget.Toolbar)

			if i.Col == (len(app.ProxyListsSlice[0])-1) && i.Row != 0 {
				lbl.Hide()
				toolbar.Hidden = false

				if len(toolbar.Items) == 0 {
					toolbar.Append(widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
						pf, _ := app.DB.GetProxyList(app.ProxyListsSlice[i.Row][len(app.ProxyListsSlice[0])-1].(int64))
						if pf != nil {
							pf.Title = pf.Title + " - Copy"
							_, err := app.DB.InsertProxyList(*pf)
							if err != nil {
								app.Logger.Error(err)
							}
							app.refreshProxyListsTab()
							app.refreshProfilesTab()
						}
					}))
					toolbar.Append(widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
						pf, _ := app.DB.GetProxyList(app.ProxyListsSlice[i.Row][len(app.ProxyListsSlice[0])-1].(int64))
						if pf != nil {
							app.editProxyListDialog(pf)
						}
					}))
					toolbar.Append(widget.NewToolbarAction(theme.DeleteIcon(), func() {
						pf, _ := app.DB.GetProxyList(app.ProxyListsSlice[i.Row][len(app.ProxyListsSlice[0])-1].(int64))
						if pf == nil {
							return
						}

						dialog.ShowConfirm("Delete?", "", func(deleted bool) {
							if deleted {
								err := app.DB.DeleteProxyList(pf.ID)
								if err != nil {
									app.Logger.Error(err)
								}
							}
							app.refreshProxyListsTab()
							app.refreshProfilesTab()
						}, app.MainWindow)

					}))
				}
			} else {
				toolbar.Hide()
				lbl.Hidden = false
				// we're just putting in textual information
				lbl.SetText(
					app.ProxyListsSlice[i.Row][i.Col].(string))
			}
		})

	colWidths := []float32{200, 200, 40}
	for i, w := range colWidths {
		t.SetColumnWidth(i, w)
	}

	return t
}

func (app *Config) refreshProxyListsTable() {
	app.ProxyListsSlice = app.getProxyListsSlice()
	app.ProxyListsTable.Refresh()

	colWidths := []float32{200, 200, 40}
	for i, w := range colWidths {
		app.ProxyListsTable.SetColumnWidth(i, w)
	}
}

func (app *Config) refreshProxyListsBottomContent() {

	all, _ := app.DB.AllProxyLists()

	btnAdd := widget.NewButtonWithIcon("Add Proxy List", theme.ContentAddIcon(), func() {
		app.addProxyListDialog()
	})
	btnClear := widget.NewButtonWithIcon("Clear Proxy Lists", theme.ContentRemoveIcon(), func() {
		dialog.ShowConfirm(
			"Delete all proxy lists?",
			fmt.Sprintf("Do you really want to delete %d proxy lists?", len(all)),
			func(deleted bool) {
				if deleted {
					for _, p := range all {
						err := app.DB.DeleteProxyList(p.ID)
						if err != nil {
							app.Logger.Error(err)
						}
					}
					app.refreshProxyListsTab()
					app.refreshProfilesTab()
				}
			}, app.MainWindow)
	})
	btnClear.Importance = widget.DangerImportance

	app.ProxyListsBottom.Objects = []fyne.CanvasObject{
		layout.NewSpacer(),
		btnAdd,
		btnClear,
	}

	app.ProxyListsBottom.Refresh()
}

func (app *Config) refreshProxyListsTopContent() {

	all, _ := app.DB.AllProxyLists()

	txt := widget.NewRichTextFromMarkdown(`## ` + strconv.Itoa(len(all)) + ` Proxy Lists Loaded`)
	app.ProxyListsTop.Objects = []fyne.CanvasObject{
		txt,
	}

	app.ProxyListsTop.Refresh()
}

func (app *Config) refreshProxyListsTab() {
	app.refreshProxyListsTopContent()
	app.refreshProxyListsTable()
}
