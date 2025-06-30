package aqua_io

import (
	"fmt"
	"strconv"

	user "github.com/UrbiJr/aqua-io/backend/internal/user"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	bot "github.com/UrbiJr/aqua-io/backend/pkg/bot"
)

type TasksTab struct {
	*container.TabItem
	TasksTop     *fyne.Container
	TasksBottom  *fyne.Container
	TasksSlice   [][]any
	TasksTable   *widget.Table
	selectedTask *user.Task
}

func (app *Config) tasksTab() *fyne.Container {
	// define tasksTab
	app.TasksTab.TasksTop = container.NewStack()
	app.TasksTab.TasksBottom = container.NewHBox()

	app.TasksSlice = app.getTasksSlice()

	// update content
	app.refreshTasksTopContent()
	app.refreshTasksBottomContent()

	app.TasksTable = app.getTasksTable()

	// define the center container
	vScroll := container.NewScroll(app.TasksTable)

	// define the tasksTab container
	tasksTabContainer := container.NewWithoutLayout(app.TasksTop, vScroll, app.TasksBottom)

	app.TasksTop.Move(fyne.NewPos(10, 10))
	app.TasksTop.Resize(fyne.NewSize(1280, 40))

	vScroll.Move(fyne.NewPos(10, 60))
	vScroll.Resize(fyne.NewSize(1280, 500))

	app.TasksBottom.Move(fyne.NewPos(300, 580))
	app.TasksBottom.Resize(fyne.NewSize(900, 50))

	return tasksTabContainer
}

func (app *Config) getTasksTable() *widget.Table {
	data := app.getTasksSlice()

	table := widget.NewTable(
		func() (rows int, cols int) {
			return len(data), len(data[0])
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

			// if last cell and not first row
			if i.Col == (len(data[0])-1) && i.Row != 0 {
				lbl.Hide()
				toolbar.Hidden = false
				var taskStatusIcon fyne.Resource

				task, _ := app.DB.GetTaskByID(app.TasksSlice[i.Row][len(data[0])-1].(int64))
				if task != nil {
					if task.Running {
						taskStatusIcon = theme.CancelIcon()
					} else {
						taskStatusIcon = theme.ContentRedoIcon()
					}
				}

				if len(toolbar.Items) == 0 {
					toolbar.Append(
						widget.NewToolbarAction(
							theme.ContentCopyIcon(), func() {

							},
						),
					)

					toolbar.Append(
						widget.NewToolbarAction(
							taskStatusIcon, func() {
								// TODO
								task.Running = !task.Running
								app.refreshTasksTab()
							},
						),
					)

					toolbar.Append(
						widget.NewToolbarAction(theme.DeleteIcon(), func() {
							dialog.ShowConfirm("Delete?", "", func(deleted bool) {
								if deleted {
									err := app.DB.DeleteTask(task.ID)
									if err != nil {
										app.Logger.Error(err)
									}
								}
								app.refreshTasksTab()
							}, app.MainWindow)
						},
						),
					)
				}
			} else {
				toolbar.Hide()
				lbl.Hidden = false
				lbl.SetText(
					app.TasksSlice[i.Row][i.Col].(string))
			}
		})

	colWidths := []float32{40, 270, 180, 180, 180, 150, 40}
	for i, w := range colWidths {
		table.SetColumnWidth(i, w)
	}

	return table
}

func getTaskStatus(taskID int64) string {
	// TODO
	return ""
}

func (app *Config) getTasksSlice() [][]any {
	var slice [][]any

	slice = append(slice, []any{"ID", "Status", "Module", "Profile", "Proxy List", "Payment Mode", "Actions"})
	allTasks, _ := app.DB.AllTasks()

	for _, x := range allTasks {
		var currentRow []any

		currentRow = append(currentRow, x.ID)
		currentRow = append(currentRow, getTaskStatus(x.ID))
		currentRow = append(currentRow, x.Module)
		p, err := app.DB.GetProfileByID(x.ProfileID)
		if err != nil {
			continue
		}
		if len(p.Title) > 30 {
			currentRow = append(currentRow, p.Title[:29]+"...")
		} else {
			currentRow = append(currentRow, p.Title)
		}

		proxyList, err := app.DB.GetProxyList(x.ProxyListID)
		if err != nil {
			continue
		}
		if len(proxyList.Title) > 30 {
			currentRow = append(currentRow, proxyList.Title[:29]+"...")
		} else {
			currentRow = append(currentRow, proxyList.Title)
		}

		currentRow = append(currentRow, x.PaymentMode)

		currentRow = append(currentRow, x.ID) // for toolbar actions

		slice = append(slice, currentRow)
	}

	return slice
}

func (app *Config) refreshTasksTable() {
	app.TasksSlice = app.getTasksSlice()
	app.TasksTable.Refresh()

	colWidths := []float32{40, 270, 180, 180, 180, 150, 40}
	for i, w := range colWidths {
		app.TasksTable.SetColumnWidth(i, w)
	}
}

func (app *Config) refreshTasksTab() {
	app.refreshTasksTable()
}

func (app *Config) addTaskDialog() dialog.Dialog {
	var profilesIDS []string
	profiles, err := app.DB.AllProfiles()
	if err != nil {
		app.Logger.Error(err)
	} else {
		for _, a := range profiles {
			profilesIDS = append(profilesIDS, a.Title)
		}
	}

	var proxyListsIDS []string
	proxyLists, err := app.DB.AllProxyLists()
	if err != nil {
		app.Logger.Error(err)
	} else {
		for _, a := range proxyLists {
			proxyListsIDS = append(proxyListsIDS, a.Title)
		}
	}

	var modulesIDS []string
	modules := bot.AllModules()
	for _, a := range modules {
		modulesIDS = append(modulesIDS, a)
	}

	profilesSelect := widget.NewSelect(profilesIDS, nil)
	profilesSelect.ClearSelected()

	proxyListsSelect := widget.NewSelect(proxyListsIDS, nil)
	proxyListsSelect.ClearSelected()

	modulesSelect := widget.NewSelect(modulesIDS, nil)
	modulesSelect.ClearSelected()

	paymentModesSelect := widget.NewSelect([]string{
		"paypal",
		"card",
	}, nil)
	paymentModesSelect.ClearSelected()

	vBox := container.NewVBox(
		widget.NewLabel("Profile"),
		profilesSelect,
		widget.NewLabel("Proxy List"),
		proxyListsSelect,
		widget.NewLabel("Module"),
		modulesSelect,
		widget.NewLabel("Payment Mode"),
		paymentModesSelect,
	)
	scrollContent := container.NewVScroll(vBox)

	addForm := dialog.NewCustomConfirm(
		"Add New Task",
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

			profile, err := app.DB.GetProfileByTitle(profilesSelect.Selected)
			if err != nil {
				return
			}
			proxyList, err := app.DB.GetProxyListByTitle(proxyListsSelect.Selected)
			if err != nil {
				return
			}

			if valid {
				t := user.Task{
					ProfileID:   profile.ID,
					ProxyListID: proxyList.ID,
					Module:      modulesSelect.Selected,
					PaymentMode: paymentModesSelect.Selected,
					Running:     false,
				}

				_, err := app.DB.InsertTask(t)

				if err != nil {
					app.Logger.Error(err)
				}
				app.refreshTasksTab()
			}
		},
		app.MainWindow,
	)

	// size and show the dialog
	addForm.Resize(fyne.NewSize(500, 600))
	addForm.Show()

	return addForm
}

func (app *Config) editTaskDialog(pf *user.Task) dialog.Dialog {
	var profilesIDS []string
	profiles, err := app.DB.AllProfiles()
	if err != nil {
		app.Logger.Error(err)
	} else {
		for _, a := range profiles {
			profilesIDS = append(profilesIDS, a.Title)
		}
	}

	var proxyListsIDS []string
	proxyLists, err := app.DB.AllProxyLists()
	if err != nil {
		app.Logger.Error(err)
	} else {
		for _, a := range proxyLists {
			proxyListsIDS = append(proxyListsIDS, a.Title)
		}
	}

	var modulesIDS []string
	modules := bot.AllModules()
	for _, a := range modules {
		modulesIDS = append(modulesIDS, a)
	}

	profilesSelect := widget.NewSelect(profilesIDS, nil)
	profileSelected, _ := app.DB.GetProfileByID(pf.ProfileID)
	if profileSelected != nil {
		profilesSelect.SetSelected(profileSelected.Title)
	}

	proxyListsSelect := widget.NewSelect(proxyListsIDS, nil)
	proxyListSelected, _ := app.DB.GetProxyList(pf.ProxyListID)
	if proxyListSelected != nil {
		proxyListsSelect.SetSelected(proxyListSelected.Title)
	}

	modulesSelect := widget.NewSelect(modulesIDS, nil)
	modulesSelect.SetSelected(pf.Module)

	paymentModesSelect := widget.NewSelect([]string{
		"paypal",
		"card",
	}, nil)
	paymentModesSelect.SetSelected(pf.PaymentMode)

	vBox := container.NewVBox(
		widget.NewLabel("Profile"),
		profilesSelect,
		widget.NewLabel("Proxy List"),
		proxyListsSelect,
		widget.NewLabel("Module"),
		modulesSelect,
		widget.NewLabel("Payment Mode"),
		paymentModesSelect,
	)
	scrollContent := container.NewVScroll(vBox)

	editForm := dialog.NewCustomConfirm(
		"Edit Profile",
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

			profile, err := app.DB.GetProfileByTitle(profilesSelect.Selected)
			if err != nil {
				return
			}
			proxyList, err := app.DB.GetProxyListByTitle(proxyListsSelect.Selected)
			if err != nil {
				return
			}

			if valid {
				t := user.Task{
					ProfileID:   profile.ID,
					ProxyListID: proxyList.ID,
					Module:      modulesSelect.Selected,
					PaymentMode: paymentModesSelect.Selected,
					Running:     false,
				}

				err := app.DB.UpdateTask(pf.ID, t)

				if err != nil {
					app.Logger.Error(err)
				}
				app.refreshTasksTab()
			}
		},
		app.MainWindow,
	)

	// size and show the dialog
	editForm.Resize(fyne.NewSize(500, 600))
	editForm.Show()

	return editForm

}

func (app *Config) refreshTasksBottomContent() {

	allTasks, _ := app.DB.AllTasks()

	btnAdd := widget.NewButtonWithIcon("Add Task", theme.ContentAddIcon(), func() {
		app.addTaskDialog()
	})
	btnClear := widget.NewButtonWithIcon("Clear Tasks", theme.ContentRemoveIcon(), func() {
		dialog.ShowConfirm(
			"Delete all tasks?",
			fmt.Sprintf("Do you really want to stop and delete %d tasks?", len(allTasks)),
			func(deleted bool) {
				if deleted {
					for _, t := range allTasks {
						err := app.DB.DeleteTask(t.ID)
						if err != nil {
							app.Logger.Error(err)
						}
					}
					app.refreshTasksTab()
				}
			}, app.MainWindow)
	})
	btnClear.Importance = widget.DangerImportance

	app.TasksBottom.Objects = []fyne.CanvasObject{
		layout.NewSpacer(),
		btnAdd,
		btnClear,
	}

	app.TasksBottom.Refresh()
}

func (app *Config) refreshTasksTopContent() {

	txt := widget.NewRichTextFromMarkdown(`## ` + strconv.Itoa(len(app.TasksSlice)-1) + ` Tasks Loaded`)
	app.TasksTop.Objects = []fyne.CanvasObject{
		txt,
	}

	app.TasksTop.Refresh()
}
