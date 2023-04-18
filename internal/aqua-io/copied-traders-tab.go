package copy_io

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/UrbiJr/aqua-io/internal/user"
	"github.com/UrbiJr/aqua-io/internal/utils"
)

type CopiedTradersTab struct {
	*container.TabItem
	OrdersSlice         [][]any
	OrdersTable         *widget.Table
	ProfileSelector     *widget.Select
	GroupSelector       *widget.Select
	SelectedProfile     *user.Profile
	ShowFromAllProfiles *widget.Check
	CopiedTradersList   *widget.List
	orders              []user.Order
	symbols             []string
	positionsLabel      *widget.Label
	SymbolList          *widget.List
}

func (app *Config) copiedTradersTab() *fyne.Container {
	app.getTraders()
	split := app.getCopiedTraders()

	max := container.NewMax(split)

	return max
}

func (app *Config) getCopiedTraders() *container.Split {

	/* RIGTH SPLIT */

	// get the top content

	var profileGroups []string
	for _, pfg := range app.User.ProfileManager.Groups {
		profileGroups = append(profileGroups, pfg.Name)
	}

	// get the profile selector
	app.CopiedTradersTab.ProfileSelector = widget.NewSelect([]string{}, func(s string) {
		group := app.User.ProfileManager.GetGroupByName(app.CopiedTradersTab.GroupSelector.Selected)
		if group != nil {
			app.CopiedTradersTab.SelectedProfile = app.User.ProfileManager.GetProfileByTitle(s, group.ID)
			app.refreshCopiedTradersTab()
		}
	})
	app.CopiedTradersTab.ProfileSelector.Disable()

	// get the profile group selector
	app.CopiedTradersTab.GroupSelector = widget.NewSelect(profileGroups, func(s string) {
		app.CopiedTradersTab.ProfileSelector.Options = []string{}
		profiles := app.User.ProfileManager.FilterByGroupName(s)
		for _, p := range profiles {
			app.CopiedTradersTab.ProfileSelector.Options = append(app.CopiedTradersTab.ProfileSelector.Options, p.Title)
		}
		if len(app.CopiedTradersTab.ProfileSelector.Options) > 0 {
			app.CopiedTradersTab.ProfileSelector.Enable()
		}
		app.CopiedTradersTab.ProfileSelector.ClearSelected()
		app.CopiedTradersTab.ProfileSelector.Refresh()
	})

	// get the show from all profiles checkbox
	app.CopiedTradersTab.ShowFromAllProfiles = widget.NewCheck("Show From All Profiles", func(b bool) {
	})
	// set it to true as default
	app.CopiedTradersTab.ShowFromAllProfiles.SetChecked(true)

	// wrap the top content in a container
	top := container.NewVBox(
		app.CopiedTradersTab.ShowFromAllProfiles,
		widget.NewSeparator(),
		widget.NewLabel("...Or Select Profile"),
		app.CopiedTradersTab.GroupSelector,
		app.CopiedTradersTab.ProfileSelector,
	)

	// get the main content

	// get the order history table
	var slice [][]any
	slice = append(slice, []any{"Profile", "Symbol", "Order ID", "Status", "Qty", "Price", "Side", "Leverage", "Created Time"})
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

	// set default columns widths
	colWidths := []float32{100, 200, 100, 200, 100, 100, 100, 100, 200}
	for i, w := range colWidths {
		app.OrdersTable.SetColumnWidth(i, w)
	}

	// get the orders tab
	ordersTab := container.NewVScroll(app.OrdersTable)

	// get the positions tab
	app.CopiedTradersTab.positionsLabel = widget.NewLabel("")
	// get the position info text
	positionInfoMarkdown := widget.NewRichTextFromMarkdown("Select a symbol to view its positions")
	// get the positions list
	app.CopiedTradersTab.SymbolList = app.getSymbolList()
	app.CopiedTradersTab.SymbolList.UnselectAll()
	app.CopiedTradersTab.SymbolList.OnSelected = func(id widget.ListItemID) {
		symbol := app.CopiedTradersTab.symbols[id]
		positionInfoMarkdown.ParseMarkdown(fmt.Sprintf("Loading %s positions...", symbol))
		go func() {
			elements := strings.Split(symbol, "|")
			group := app.User.ProfileManager.GetGroupByName(elements[0])
			profile := app.User.ProfileManager.GetProfileByTitle(elements[1], group.ID)
			positionInfoArr := app.getPositionInfo("linear", elements[2], *profile)
			markdownText := ""
			for i, p := range positionInfoArr {
				//Unix Timestamp to time.Time
				timeT := time.Unix(0, p.UpdatedTime*int64(time.Millisecond))
				layout := "02-01-2006 15:04:05"
				readable := timeT.Format(layout)
				var mode, side string
				if p.PositionIdx == 0 {
					mode = "one-way mode position"
				} else if p.PositionIdx == 1 {
					mode = "Buy side of hedge-mode position"
				} else if p.PositionIdx == 2 {
					mode = "Sell side of hedge-mode position"
				}
				if p.Side == "None" {
					side = "Empty position"
				}

				markdownText += fmt.Sprintf("%d. Position", i) + `
` + "    ```" + `
    Symbol:                     ` + p.Symbol + `
	Mode:                       ` + mode + `
    Leverage:                   ` + fmt.Sprintf("%d", p.Leverage) + `
    Average Entry Price:        ` + fmt.Sprintf("%.2f", p.AvgPrice) + `
    Position Liquidation Price: ` + fmt.Sprintf("%.2f", p.LiqPrice) + `
    Take Profit:                ` + p.TakeProfit.(string) + `
    Stop Loss:                  ` + p.StopLoss.(string) + `
    Position Value:             ` + fmt.Sprintf("%.2f", p.PositionValue) + `
    Unrealised Pnl:             ` + fmt.Sprintf("%.2f", p.UnrealisedPnl) + `
    Cumulative Realised Pnl:    ` + fmt.Sprintf("%.2f", p.CumRealisedPnl) + `
    Market Price:               ` + fmt.Sprintf("%.2f", p.MarkPrice) + `
    Last Update Time:           ` + readable + `
    Side (buy/sell/empty):      ` + side + `
    Position Status:            ` + p.PositionStatus + `
` + "    ```\n"
			}
			positionInfoMarkdown.ParseMarkdown(markdownText)
		}()
	}
	app.CopiedTradersTab.SymbolList.OnUnselected = func(id widget.ListItemID) {
		positionInfoMarkdown.ParseMarkdown("Select a symbol to view its positions")
	}

	positionsTab := container.NewBorder(
		app.CopiedTradersTab.positionsLabel,
		nil,
		nil,
		nil,
		container.NewHSplit(app.CopiedTradersTab.SymbolList, container.NewVScroll(positionInfoMarkdown)))

	// get the tabs container
	tabs := container.NewAppTabs(
		container.NewTabItem("Order History", ordersTab),
		container.NewTabItem("Positions", positionsTab))
	tabs.SetTabLocation(container.TabLocationTop)
	tabs.OnSelected = func(ti *container.TabItem) {
		// when selecting the order history tab...
		if ti.Text == "Order History" {
			// ...clear selection from symbol list on the positions tab
			app.CopiedTradersTab.SymbolList.UnselectAll()
		}
	}

	/* LEFT SPLIT */

	// get the copied traders list
	app.CopiedTradersList = widget.NewList(
		func() int {
			return len(app.User.CopiedTradersManager.Traders)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.AccountIcon()), widget.NewLabel("Template Object"), widget.NewToolbar())
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			// define the icon + trader name + toolbar actions for each trader
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

	// set default height for each list item
	for i := range app.User.CopiedTradersManager.Traders {
		app.CopiedTradersList.SetItemHeight(i, 50)
	}

	// refresh the content
	app.refreshCopiedTradersTab()

	app.CopiedTradersTab.ShowFromAllProfiles.OnChanged = func(b bool) {
		positionInfoMarkdown.ParseMarkdown("Select a symbol to view its positions")
		app.refreshCopiedTradersTab()
	}

	// get the horizontal split
	return container.NewHSplit(container.NewBorder(
		widget.NewLabelWithStyle("Copied Traders", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		app.CopiedTradersList,
	), container.NewBorder(
		top,
		nil,
		nil,
		nil,
		tabs))
}

func (app *Config) getTraders() {
	traders, err := app.DB.AllTraders()
	if err != nil {
		app.Logger.Error(err)
		app.Quit()
	}
	app.User.CopiedTradersManager.Traders = traders
}

func (app *Config) getOrderSlice() [][]any {
	var slice [][]any

	slice = append(slice, []any{"Profile", "Symbol", "Order ID", "Status", "Qty", "Price", "Side", "Leverage", "Created Time"})

	for _, x := range app.CopiedTradersTab.orders {
		var currentRow []any

		if x.OrderStatus == "Filled" {
			profile := app.User.ProfileManager.GetProfileByID(x.ProfileID, x.ProfileGroupID)

			currentRow = append(currentRow, profile.Title)

			currentRow = append(currentRow, fmt.Sprintf("%s Perpetual", x.Symbol))

			if len(x.OrderID) > 16 {
				currentRow = append(currentRow, x.OrderID[:10]+"...")
			} else {
				currentRow = append(currentRow, x.OrderID)
			}

			currentRow = append(currentRow, x.OrderStatus)

			currentRow = append(currentRow, fmt.Sprintf("%.2f", x.Qty))

			currentRow = append(currentRow, fmt.Sprintf("%.2f", x.Price))

			currentRow = append(currentRow, x.Side)

			currentRow = append(currentRow, fmt.Sprintf("%d", x.IsLeverage))

			//Unix Timestamp to time.Time
			timeT := time.Unix(0, x.CreatedTime*int64(time.Millisecond))
			layout := "02-01-2006 15:04:05"
			readable := timeT.Format(layout)

			currentRow = append(currentRow, readable)

			slice = append(slice, currentRow)
		}
	}

	return slice
}

func (app *Config) getOrders(allProfiles bool, profile *user.Profile) []user.Order {
	var allOrders []user.Order
	var allSymbols []string

	if allProfiles {
		for _, p := range app.User.Profiles {
			//TODO: better check if bybit or binance api key, add binance api key code
			if p.BybitApiKey != "" {
				byBitOrders := app.getOrderHistory("spot", p)
				if byBitOrders != nil {
					allOrders = append(allOrders, byBitOrders...)
				}
			}
		}
	} else {
		if profile != nil {
			if profile.BybitApiKey != "" {
				byBitOrders := app.getOrderHistory("spot", *profile)
				if byBitOrders != nil {
					allOrders = append(allOrders, byBitOrders...)
				}
			}
		}
	}

	app.CopiedTradersTab.orders = allOrders

	for _, o := range allOrders {
		if o.OrderStatus == "Filled" {
			groupName := app.User.ProfileManager.GetGroupByID(o.ProfileGroupID).Name
			profileName := app.User.ProfileManager.GetProfileByID(o.ProfileID, o.ProfileGroupID).Title
			if !utils.Contains(allSymbols, fmt.Sprintf("%s|%s|%s", groupName, profileName, o.Symbol)) {
				allSymbols = append(allSymbols, fmt.Sprintf("%s|%s|%s", groupName, profileName, o.Symbol))
			}

		}
	}

	app.CopiedTradersTab.symbols = allSymbols

	return allOrders
}

func (app *Config) getSymbolList() *widget.List {

	list := widget.NewList(
		func() int {
			return len(app.CopiedTradersTab.symbols)
		},
		func() fyne.CanvasObject {
			return container.NewMax(widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(app.CopiedTradersTab.symbols[id])
		},
	)

	return list
}

func (app *Config) refreshCopiedTradersTab() {
	go func() {
		// refresh orders table
		app.getOrders(app.ShowFromAllProfiles.Checked, app.CopiedTradersTab.SelectedProfile)
		app.OrdersSlice = app.getOrderSlice()
		app.OrdersTable.Refresh()

		// refresh opened positions title
		if app.ShowFromAllProfiles.Checked {
			app.CopiedTradersTab.positionsLabel.SetText("Showing opened positions from all profiles")
		} else if app.CopiedTradersTab.SelectedProfile != nil {
			app.CopiedTradersTab.positionsLabel.SetText(fmt.Sprintf("Showing opened positions from profile %s", app.CopiedTradersTab.SelectedProfile.Title))
		}

		// refresh opened positions list
		app.CopiedTradersTab.SymbolList.Refresh()
	}()
}
