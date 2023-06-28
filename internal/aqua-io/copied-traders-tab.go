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
)

type CopiedTradersTab struct {
	*container.TabItem
	limitMarketOrdersSlice  [][]any
	tpSlOrdersSlice         [][]any
	limitMarketOrdersTable  *widget.Table
	tpSlOrdersTable         *widget.Table
	CopiedTradersList       *widget.List
	limitMarketOrders       []user.Order
	tpSlOrders              []user.Order
	positionsContainer      *fyne.Container
	selectedCopiedTrader    *user.Profile
	profilesWithTraderNames []string
}

func (app *Config) copiedTradersTab() *fyne.Container {
	app.CopiedTradersTab.profilesWithTraderNames = app.formatCopiedTradersList(app.User.GetProfilesWithTrader())
	split := app.getCopiedTraders()

	max := container.NewMax(split)

	return max
}

func (app *Config) formatCopiedTradersList(profilesWithTrader []user.Profile) []string {
	var names []string
	for _, p := range profilesWithTrader {
		names = append(names, fmt.Sprintf("%s\n%s", p.Title, p.TraderID))
	}
	return names
}

func (app *Config) getCopiedTraders() *container.Split {

	/* LEFT SPLIT */

	// get the copied traders list (profiles with a trader)
	app.CopiedTradersList = widget.NewList(
		func() int {
			return len(app.CopiedTradersTab.profilesWithTraderNames)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.AccountIcon()), widget.NewLabel("Template Object"), widget.NewToolbar())
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			// get the current trader...
			trader := app.CopiedTradersTab.profilesWithTraderNames[id]
			traderID := strings.Split(trader, "\n")[1]

			// display profile title and trader nickname
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(trader)

			toolbar := item.(*fyne.Container).Objects[2].(*widget.Toolbar)
			if len(toolbar.Items) == 0 {
				toolbar.Append(widget.NewToolbarAction(theme.VisibilityIcon(), func() {
					app.traderDialog(user.Trader{}, traderID)
				}))
				toolbar.Append(widget.NewToolbarAction(theme.DeleteIcon(), func() {
					dialog.ShowConfirm("Stop copying?", "Doing that will NOT close positions copied from this trader.", func(deleted bool) {
						if deleted {
							err := app.stopCopyingTrader(user.Trader{}, traderID)
							if err != nil {
								app.Logger.Error(err)
							}
						}
					}, app.MainWindow)
				}))
			}
		},
	)

	/* RIGTH SPLIT */

	// get the top content

	// get the refresh toolbar button
	topRightToolbar := widget.NewToolbar(widget.NewToolbarSpacer(), widget.NewToolbarAction(
		theme.ViewRefreshIcon(), func() {
			var limitMarketOrdersSlice, tpSlOrdersSlice [][]any
			limitMarketOrdersSlice = append(limitMarketOrdersSlice, []any{"Symbol", "Order Type", "Side", "Avg. Filled Price", "Filled Qty", "Order Price", "Order Qty", "Order Status", "Order Time", "Order ID"})
			tpSlOrdersSlice = append(tpSlOrdersSlice, []any{"Symbol", "Order Type", "Side", "Order Value", "Order Qty", "Order Price", "Trigger Price", "Order Status", "Order Time", "Order ID"})
			app.limitMarketOrdersSlice = limitMarketOrdersSlice
			app.tpSlOrdersSlice = tpSlOrdersSlice
			app.CopiedTradersTab.positionsContainer.RemoveAll()
			app.CopiedTradersTab.positionsContainer.Add(widget.NewLabel("Getting positions..."))
			app.refreshCopiedTradersTab(true)
		}))

	// wrap the top content in a container
	top := container.NewHBox(
		topRightToolbar,
	)

	// get the main content

	// get the order history table
	var limitMarketOrdersSlice, tpSlOrdersSlice [][]any
	limitMarketOrdersSlice = append(limitMarketOrdersSlice, []any{"Symbol", "Order Type", "Side", "Avg. Filled Price", "Filled Qty", "Order Price", "Order Qty", "Order Status", "Order Time", "Order ID"})
	tpSlOrdersSlice = append(tpSlOrdersSlice, []any{"Symbol", "Order Type", "Side", "Order Value", "Order Qty", "Order Price", "Trigger Price", "Order Status", "Order Time", "Order ID"})
	app.limitMarketOrdersSlice = limitMarketOrdersSlice
	app.tpSlOrdersSlice = tpSlOrdersSlice
	app.limitMarketOrdersTable = widget.NewTable(
		func() (int, int) {
			return len(app.limitMarketOrdersSlice), len(app.limitMarketOrdersSlice[0])
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

			if i.Col == 2 && i.Row != 0 {
				lbl.Hide()
				canvasText.Hidden = false
				side := app.limitMarketOrdersSlice[i.Row][i.Col].(string)
				if side == "Buy" {
					canvasText.Color = color.RGBA{R: 14, G: 203, B: 129, A: 255}
				} else if side == "Sell" {
					canvasText.Color = color.RGBA{R: 246, G: 70, B: 93, A: 255}
				}
				canvasText.Text = side
			} else {
				canvasText.Hide()
				lbl.Hidden = false
				lbl.SetText(
					app.limitMarketOrdersSlice[i.Row][i.Col].(string))
			}
		})

	app.tpSlOrdersTable = widget.NewTable(
		func() (int, int) {
			return len(app.tpSlOrdersSlice), len(app.tpSlOrdersSlice[0])
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

			if i.Col == 2 && i.Row != 0 {
				lbl.Hide()
				canvasText.Hidden = false
				side := app.tpSlOrdersSlice[i.Row][i.Col].(string)
				if side == "Buy" {
					canvasText.Color = color.RGBA{R: 14, G: 203, B: 129, A: 255}
				} else if side == "Sell" {
					canvasText.Color = color.RGBA{R: 246, G: 70, B: 93, A: 255}
				}
				canvasText.Text = side
			} else {
				canvasText.Hide()
				lbl.Hidden = false
				lbl.SetText(
					app.tpSlOrdersSlice[i.Row][i.Col].(string))
			}
		})

	// set default columns widths
	colWidths := []float32{100, 100, 50, 100, 100, 100, 100, 200, 200, 200}
	for i, w := range colWidths {
		app.tpSlOrdersTable.SetColumnWidth(i, w)
	}

	colWidths = []float32{100, 100, 50, 100, 100, 100, 100, 200, 200, 200}
	for i, w := range colWidths {
		app.limitMarketOrdersTable.SetColumnWidth(i, w)
	}

	// get the orders tab
	limitMarketOrdersTab := container.NewScroll(app.limitMarketOrdersTable)

	// get the open orders tab
	tpSlOrdersTab := container.NewScroll(app.tpSlOrdersTable)

	// get the positions container
	app.CopiedTradersTab.positionsContainer = container.NewVBox(widget.NewLabel("Select a profile"))
	positionsTab := container.NewVScroll(app.CopiedTradersTab.positionsContainer)

	// get the tabs container
	tabs := container.NewAppTabs(
		container.NewTabItem("Positions", positionsTab),
		container.NewTabItem("Limit/Market Orders", limitMarketOrdersTab),
		container.NewTabItem("TP/SL Orders", tpSlOrdersTab))
	tabs.SetTabLocation(container.TabLocationTop)

	// when a trader is selected...
	app.CopiedTradersList.OnSelected = func(id widget.ListItemID) {
		// ... get the selected item (profile)
		trader := app.CopiedTradersTab.profilesWithTraderNames[id]
		traderID := strings.Split(trader, "\n")[1]
		profile := app.User.ProfileManager.GetProfileByTraderID(traderID)
		app.CopiedTradersTab.selectedCopiedTrader = profile

		// update order history
		app.updateOrderHistoryContent(profile)

		// update positions
		app.CopiedTradersTab.positionsContainer.RemoveAll()
		app.CopiedTradersTab.positionsContainer.Add(widget.NewLabel("Getting positions..."))
		app.updatePositionsContent(profile)
	}

	// refresh the content
	app.refreshCopiedTradersTab(false)

	// get the horizontal split
	return container.NewHSplit(container.NewBorder(
		widget.NewLabelWithStyle("Profiles", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
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

func (app *Config) getTpSlOrdersSlice() [][]any {
	var slice [][]any

	slice = append(slice, []any{"Symbol", "Order Type", "Side", "Order Value", "Order Qty", "Order Price", "Trigger Price", "Order Status", "Order Time", "Order ID"})

	for _, x := range app.tpSlOrders {
		var currentRow []any

		currentRow = append(currentRow, x.Symbol)

		currentRow = append(currentRow, x.OrderType)

		currentRow = append(currentRow, x.Side)

		orderValue := ""
		if x.Qty > 0 && x.Price > 0 {
			orderValue = fmt.Sprintf("%f", x.Qty*x.Price)
		}
		currentRow = append(currentRow, orderValue)

		currentRow = append(currentRow, fmt.Sprintf("%f", x.Qty))

		currentRow = append(currentRow, fmt.Sprintf("%f", x.Price))

		if x.TriggerDirection == 0 {
			currentRow = append(currentRow, fmt.Sprintf(">= %.2f", x.TriggerPrice))
		} else {
			currentRow = append(currentRow, fmt.Sprintf("<= %.2f", x.TriggerPrice))
		}

		currentRow = append(currentRow, x.OrderStatus)

		//Unix Timestamp to time.Time
		timeT := time.Unix(0, x.CreatedTime*int64(time.Millisecond))
		layout := "02-01-2006 15:04:05"
		readable := timeT.Format(layout)

		currentRow = append(currentRow, readable)

		if len(x.OrderID) > 26 {
			currentRow = append(currentRow, x.OrderID[:20]+"...")
		} else {
			currentRow = append(currentRow, x.OrderID)
		}

		slice = append(slice, currentRow)
	}

	return slice
}

func (app *Config) getLimitMarketOrdersSlice() [][]any {
	var slice [][]any

	slice = append(slice, []any{"Symbol", "Order Type", "Side", "Avg. Filled Price", "Filled Qty", "Order Price", "Order Qty", "Order Status", "Order Time", "Order ID"})

	for _, x := range app.limitMarketOrders {
		var currentRow []any

		currentRow = append(currentRow, x.Symbol)

		currentRow = append(currentRow, x.OrderType)

		currentRow = append(currentRow, x.Side)

		currentRow = append(currentRow, fmt.Sprintf("%f", x.AvgFilledPrice))

		currentRow = append(currentRow, fmt.Sprintf("%f", x.FilledQty))

		currentRow = append(currentRow, fmt.Sprintf("%.2f", x.Price))

		currentRow = append(currentRow, fmt.Sprintf("%f", x.Qty))

		currentRow = append(currentRow, x.OrderStatus)

		//Unix Timestamp to time.Time
		timeT := time.Unix(0, x.CreatedTime*int64(time.Millisecond))
		layout := "02-01-2006 15:04:05"
		readable := timeT.Format(layout)

		currentRow = append(currentRow, readable)

		if len(x.OrderID) > 26 {
			currentRow = append(currentRow, x.OrderID[:20]+"...")
		} else {
			currentRow = append(currentRow, x.OrderID)
		}

		slice = append(slice, currentRow)
	}

	return slice
}

func (app *Config) getOrders(profile *user.Profile) {
	var limitMarketOrders, tpSlOrders []user.Order

	if profile != nil {
		if profile.BybitApiKey != "" {
			byBitOrders, err := app.fetchOrderHistory("spot", "Order", *profile)
			if err != nil && strings.Contains(err.Error(), "Timestamp") {
				app.App.SendNotification(fyne.NewNotification(
					"⚠️ Error getting orders",
					err.Error(),
				))
			}
			if byBitOrders != nil {
				limitMarketOrders = append(limitMarketOrders, byBitOrders...)
			}
			byBitOrders, err = app.fetchOrderHistory("spot", "tpslOrder", *profile)
			if err != nil && strings.Contains(err.Error(), "Timestamp") {
				app.App.SendNotification(fyne.NewNotification(
					"⚠️ Error getting orders",
					err.Error(),
				))
			}
			if byBitOrders != nil {
				tpSlOrders = append(tpSlOrders, byBitOrders...)
			}
		}
	}

	app.CopiedTradersTab.limitMarketOrders = limitMarketOrders
	app.CopiedTradersTab.tpSlOrders = tpSlOrders
}

func (app *Config) updateOrderHistoryContent(profile *user.Profile) {
	go func() {
		app.getOrders(profile)
		app.limitMarketOrdersSlice = app.getLimitMarketOrdersSlice()
		app.tpSlOrdersSlice = app.getTpSlOrdersSlice()
		app.limitMarketOrdersTable.Refresh()
		app.tpSlOrdersTable.Refresh()
	}()
}

func (app *Config) updatePositionsContent(p *user.Profile) {
	go func() {
		var positionInfoArr []user.PositionInfo
		openedPositions := app.User.CopiedTradersManager.GetOpenedPositionsByProfileID(p.ID)
		for _, position := range openedPositions {
			positionInfoArr = append(positionInfoArr, app.getPositionInfo("linear", position.Symbol, position.OrderID, *p)...)
		}

		app.CopiedTradersTab.positionsContainer.RemoveAll()

		for _, pos := range positionInfoArr {
			//Unix Timestamp to time.Time
			timeT := time.Unix(0, pos.UpdatedTime*int64(time.Millisecond))
			layout := "02-01-2006 15:04:05"
			readable := timeT.Format(layout)
			var mode, side string
			if pos.PositionIdx == 0 {
				mode = "one-way mode position"
			} else if pos.PositionIdx == 1 {
				mode = "Buy side of hedge-mode position"
			} else if pos.PositionIdx == 2 {
				mode = "Sell side of hedge-mode position"
			}
			if pos.Side == "None" {
				side = "Empty position"
			}

			markdownText := "```" + `
    Order ID:                   ` + pos.OrderID + `
    Symbol:                     ` + pos.Symbol + `
    Mode:                       ` + mode + `
    Leverage:                   ` + fmt.Sprintf("%d", pos.Leverage) + `
    Average Entry Price:        ` + fmt.Sprintf("%.2f", pos.AvgPrice) + `
    Position Liquidation Price: ` + fmt.Sprintf("%.2f", pos.LiqPrice) + `
    Take Profit:                ` + pos.TakeProfit.(string) + `
    Stop Loss:                  ` + pos.StopLoss.(string) + `
    Position Value:             ` + fmt.Sprintf("%.2f", pos.PositionValue) + `
    Unrealised Pnl:             ` + fmt.Sprintf("%.2f", pos.UnrealisedPnl) + `
    Cumulative Realised Pnl:    ` + fmt.Sprintf("%.2f", pos.CumRealisedPnl) + `
    Market Price:               ` + fmt.Sprintf("%.2f", pos.MarkPrice) + `
    Last Update Time:           ` + readable + `
    Side (buy/sell/empty):      ` + side + `
    Position Status:            ` + pos.PositionStatus + `
` + "```"

			app.CopiedTradersTab.positionsContainer.Add(widget.NewRichTextFromMarkdown(markdownText))
			app.CopiedTradersTab.positionsContainer.Add(container.NewGridWithColumns(3,
				widget.NewButtonWithIcon("Set TP", theme.DocumentCreateIcon(), func() {

				}),
				widget.NewButtonWithIcon("Set SL", theme.DocumentCreateIcon(), func() {

				}),
				widget.NewButtonWithIcon("Close", theme.DeleteIcon(), func() {
					dialog.ShowConfirm("Close Position?", fmt.Sprintf("Confirming will cancel %s order with ID %s.", pos.Symbol, pos.OrderID), func(deleted bool) {
						if deleted {
							err := app.cancelOrder(p, "spot", pos.OrderID, pos.Symbol)
							if err != nil {
								app.App.SendNotification(fyne.NewNotification(
									"⚠️ Close Position Failed",
									fmt.Sprintf("Error: %s", err.Error()),
								))
							} else {
								err = app.DB.DeleteOpenedPosition(pos.OrderID)
								if err != nil {
									app.Logger.Error(err)
								} else {
									app.User.CopiedTradersManager.DeleteOpenedPosition(pos.OrderID)
									app.Logger.Debug(fmt.Sprintf("closed %s position with order ID %s", pos.Symbol, pos.OrderID))
									app.App.SendNotification(fyne.NewNotification(
										"🔴 Successfully Closed Position",
										fmt.Sprintf("Order %s cancelled", pos.OrderID),
									))
								}
								app.refreshCopiedTradersTab(true)
							}
						}
					}, app.MainWindow)
				}),
			))
			app.CopiedTradersTab.positionsContainer.Add(widget.NewSeparator())
		}
		app.CopiedTradersTab.positionsContainer.Refresh()
	}()
}

func (app *Config) refreshCopiedTradersTab(isCopiedTraderSelected bool) {
	if isCopiedTraderSelected && app.CopiedTradersTab.selectedCopiedTrader != nil {
		// update content based on the current selection
		go func() {
			app.updateOrderHistoryContent(app.CopiedTradersTab.selectedCopiedTrader)
			app.updatePositionsContent(app.CopiedTradersTab.selectedCopiedTrader)
		}()
	} else {
		// reset orders table
		var limitMarketOrdersSlice, tpSlOrdersSlice [][]any
		limitMarketOrdersSlice = append(limitMarketOrdersSlice, []any{"Symbol", "Order Type", "Side", "Avg. Filled Price", "Filled Qty", "Order Price", "Order Qty", "Order Status", "Order Time", "Order ID"})
		tpSlOrdersSlice = append(tpSlOrdersSlice, []any{"Symbol", "Order Type", "Side", "Order Value", "Order Qty", "Order Price", "Trigger Price", "Order Status", "Order Time", "Order ID"})
		app.limitMarketOrdersSlice = limitMarketOrdersSlice
		app.tpSlOrdersSlice = tpSlOrdersSlice
		app.limitMarketOrdersTable.Refresh()
		app.tpSlOrdersTable.Refresh()

		// reset positions
		app.CopiedTradersTab.positionsContainer.RemoveAll()
		app.CopiedTradersTab.positionsContainer.Add(widget.NewLabel("Select a profile"))
	}

	app.CopiedTradersTab.profilesWithTraderNames = app.formatCopiedTradersList(app.User.ProfileManager.GetProfilesWithTrader())
	// set default height for each list item
	for i := range app.CopiedTradersTab.profilesWithTraderNames {
		app.CopiedTradersList.SetItemHeight(i, 50)
	}
	app.CopiedTradersList.Refresh()
}

func (app *Config) getOpenedPositions() {
	openedPositions, err := app.DB.AllOpenedPositions()
	if err != nil {
		app.Logger.Error(err)
		app.Quit()
	}
	app.User.CopiedTradersManager.OpenedPositions = openedPositions
}
