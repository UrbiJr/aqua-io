package aqua_io

import (
	user "github.com/UrbiJr/aqua-io/backend/internal/user"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TasksTab struct {
	*container.TabItem
	TasksSlice   [][]any
	TasksTable   *widget.Table
	selectedTask *user.Task
}

func (app *Config) tasksTab() *fyne.Container {
	grid := app.getTasksTab()

	max := container.NewStack(grid)

	return max
}

func (app *Config) getTasksTab() *widget.Table {
	data := app.getTasksSlice()

	table := widget.NewTable(
		func() (rows int, cols int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			ctr := container.NewVBox(widget.NewLabel(""))
			return ctr
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

				task, _ := app.DB.GetTaskByID(app.ProfilesSlice[i.Row][len(data[0])-1].(int64))
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
					app.ProfilesSlice[i.Row][i.Col].(string))
			}
		})

	return table
}

func (app *Config) getTasksSlice() [][]any {
	var slice [][]any

	slice = append(slice, []any{"ID", "Module", "Profile", "Proxy List", "Payment Mode", "Actions"})
	allTasks, _ := app.DB.AllTasks()

	for _, x := range allTasks {
		var currentRow []any

		currentRow = append(currentRow, x.ID)
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
	
	/*
		colWidths := []float32{220, 270, 200, 200, 100, 40}
		for i, w := range colWidths {
			app.ProfilesTable.SetColumnWidth(i, w)
		}
	*/
}

func (app *Config) refreshTasksTab() {
	app.refreshTasksTable()
}
