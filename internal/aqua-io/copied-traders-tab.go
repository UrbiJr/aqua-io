package copy_io

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/aqua-io/internal/user"
)

type CopiedTradersTab struct {
	*container.TabItem
	OrdersSlice       [][]any
	OrdersTable       *widget.Table
	ProfileSelector   *widget.Select
	GroupSelector     *widget.Select
	SelectedProfile   *user.Profile
	CopiedTradersList *widget.List
}

func (app *Config) copiedTradersTab() *fyne.Container {
	app.getTraders()
	copiedTradersList := app.getCopiedTraders()

	max := container.NewMax(copiedTradersList)

	return max
}

func (app *Config) getCopiedTraders() *container.Split {
	var slice [][]any

	slice = append(slice, []any{"Profile", "Symbol", "Order ID", "Status", "Qty", "Price", "Side", "Is Leverage", "Created Time"})
	app.OrdersSlice = slice
	app.OrdersTable = widget.NewTable(
		func() (int, int) {
			return len(app.OrdersSlice), len(app.OrdersSlice[0])
		},
		func() fyne.CanvasObject {
			lbl := widget.NewLabel("")
			canvasText := canvas.NewText("", nil)
			canvasText.Hide()
			return container.NewMax(lbl, canvasText)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			container := o.(*fyne.Container)
			lbl := container.Objects[0].(*widget.Label)
			canvasText := container.Objects[1].(*canvas.Text)

			if i.Col == 999 && i.Row != 0 {
				lbl.Hide()
				canvasText.Hidden = false
				funding := app.OrdersSlice[i.Row][i.Col].(float64)
				if funding > 0 {
					canvasText.Color = color.RGBA{R: 14, G: 203, B: 129, A: 255}
				} else {
					canvasText.Color = color.RGBA{R: 246, G: 70, B: 93, A: 255}
				}
				canvasText.Text = fmt.Sprintf("%f", funding)
			} else {
				canvasText.Hide()
				lbl.Hidden = false
				lbl.SetText(
					app.OrdersSlice[i.Row][i.Col].(string))
			}
		})
	colWidths := []float32{100, 100, 200, 200, 100, 100, 100, 100, 200}
	for i, w := range colWidths {
		app.OrdersTable.SetColumnWidth(i, w)
	}

	var profileGroups []string
	for _, pfg := range app.User.ProfileManager.Groups {
		profileGroups = append(profileGroups, pfg.Name)
	}

	app.CopiedTradersTab.ProfileSelector = widget.NewSelect([]string{}, func(s string) {
		group := app.User.ProfileManager.GetGroupByName(app.CopiedTradersTab.GroupSelector.Selected)
		if group != nil {
			app.CopiedTradersTab.SelectedProfile = app.User.ProfileManager.GetProfileByTitle(s, group.ID)
		}
	})
	app.CopiedTradersTab.ProfileSelector.Disable()
	app.CopiedTradersTab.GroupSelector = widget.NewSelect(profileGroups, func(s string) {
		app.CopiedTradersTab.ProfileSelector.Options = []string{}
		profiles := app.User.ProfileManager.FilterByGroupName(s)
		for _, p := range profiles {
			app.CopiedTradersTab.ProfileSelector.Options = append(app.CopiedTradersTab.ProfileSelector.Options, p.Title)
		}
		if len(app.CopiedTradersTab.ProfileSelector.Options) > 0 {
			app.CopiedTradersTab.ProfileSelector.Enable()
		}
		app.CopiedTradersTab.ProfileSelector.Refresh()
	})

	go func() {
		app.refreshOrdersTable()
	}()

	top := container.NewVBox(
		widget.NewCheck("Show From All Profiles", func(b bool) {}),
		widget.NewSeparator(),
		widget.NewLabel("...Or Select Profile"),
		app.CopiedTradersTab.GroupSelector,
		app.CopiedTradersTab.ProfileSelector,
	)

	border := container.NewBorder(
		top,
		nil,
		nil,
		nil,
		container.NewVScroll(app.OrdersTable),
	)

	app.CopiedTradersList = widget.NewList(
		func() int {
			return len(app.User.CopiedTradersManager.Traders)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.AccountIcon()), widget.NewLabel("Template Object"), widget.NewToolbar())
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			t := app.User.CopiedTradersManager.Traders[id]
			toolbar := item.(*fyne.Container).Objects[2].(*widget.Toolbar)
			if len(toolbar.Items) == 0 {
				toolbar.Append(widget.NewToolbarAction(theme.VisibilityIcon(), func() {
					app.traderDialog(t)
				}))
				toolbar.Append(widget.NewToolbarAction(theme.DeleteIcon(), func() {
					dialog.ShowConfirm("Stop copying?", "Confirming will delete all current positions copied from this trader.", func(deleted bool) {
						app.stopCopyingTrader(t)
					}, app.MainWindow)
				}))
			}
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(t.NickName + "\n" + t.EncryptedUid)
		},
	)
	app.CopiedTradersList.OnSelected = func(id widget.ListItemID) {
	}
	app.CopiedTradersList.OnUnselected = func(id widget.ListItemID) {
	}
	for i := range app.User.CopiedTradersManager.Traders {
		app.CopiedTradersList.SetItemHeight(i, 50)
	}

	return container.NewHSplit(container.NewBorder(
		widget.NewLabelWithStyle("Copied Traders", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		app.CopiedTradersList,
	), container.NewBorder(widget.NewLabelWithStyle("Orders History", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		border))
}

func (app *Config) getTraders() {
	traders, err := app.DB.AllTraders()
	if err != nil {
		app.Logger.Error(err)
		app.Quit()
	}
	app.User.CopiedTradersManager.Traders = traders
}

func (app *Config) getOrdersSlice() [][]any {
	var slice [][]any
	var allOrders []user.Order

	slice = append(slice, []any{"Profile", "Symbol", "Order ID", "Status", "Qty", "Price", "Side", "Is Leverage", "Created Time"})
	for _, p := range app.User.Profiles {
		//TODO: better check if bybit or binance api key, add binance api key code
		if p.BybitApiKey != "" {
			byBitOrders := app.getOrderHistory("spot", p)
			if byBitOrders != nil {
				allOrders = append(allOrders, byBitOrders...)
			}
		}
	}

	for _, x := range allOrders {
		var currentRow []any

		profile := app.User.ProfileManager.GetProfileByID(x.ProfileID, x.ProfileGroupID)

		currentRow = append(currentRow, profile.Title)

		currentRow = append(currentRow, fmt.Sprintf("%s Perpetual", x.Symbol))

		currentRow = append(currentRow, x.OrderID)

		currentRow = append(currentRow, x.OrderStatus)

		currentRow = append(currentRow, fmt.Sprintf("%.2f", x.Qty))

		currentRow = append(currentRow, fmt.Sprintf("%.2f", x.Price))

		currentRow = append(currentRow, x.Side)

		currentRow = append(currentRow, fmt.Sprintf("%.2f", x.IsLeverage))

		currentRow = append(currentRow, fmt.Sprintf("%d", x.CreatedTime))

		slice = append(slice, currentRow)
	}

	return slice
}

func (app *Config) refreshOrdersTable() {
	app.OrdersSlice = app.getOrdersSlice()
	app.OrdersTable.Refresh()
}
