package copy_io

import (
	"fmt"
	"image/color"
	"strconv"
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
	limitMarketOrdersSlice [][]any
	tpSlOrdersSlice        [][]any
	limitMarketOrdersTable *widget.Table
	tpSlOrdersTable        *widget.Table
	CopiedTradersList      *widget.List
	limitMarketOrders      []user.Order
	tpSlOrders             []user.Order
	positionsTable         *widget.Table
	positionsSlice         [][]any
	selectedCopiedTrader   *user.Profile
	profilesWithTrader     []string
}

func (app *Config) copiedTradersTab() *fyne.Container {
	app.CopiedTradersTab.profilesWithTrader = app.formatCopiedTradersList(app.User.GetProfilesWithTrader())
	app.positionsSlice = app.getPositionsSlice()
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
			return len(app.CopiedTradersTab.profilesWithTrader)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.AccountIcon()), widget.NewLabel("Template Object"), widget.NewToolbar())
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			// get the current trader...
			profileTitleAndTraderID := app.CopiedTradersTab.profilesWithTrader[id]
			traderID := strings.Split(profileTitleAndTraderID, "\n")[1]

			// display profile title and trader nickname
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(profileTitleAndTraderID + "\t") // to align horizontally

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
			app.refreshCopiedTradersTab()
		}))

	// wrap the top content in a container
	top := container.NewHBox(
		topRightToolbar,
	)

	// get the main content

	// get the positions table
	app.positionsTable = app.getPositionsTable()

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

	// get the positions tab
	positionsTab := container.NewScroll(app.positionsTable)

	// get the tabs container
	tabs := container.NewAppTabs(
		container.NewTabItem("Positions", positionsTab),
		container.NewTabItem("Limit/Market Orders", limitMarketOrdersTab),
		container.NewTabItem("TP/SL Orders", tpSlOrdersTab))
	tabs.SetTabLocation(container.TabLocationTop)

	// refresh the content
	app.refreshCopiedTradersTab()

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

func (app *Config) updateOrderHistoryContent() {
	for _, profile := range app.User.Profiles {
		go func() {
			app.getOrders(&profile)
			app.limitMarketOrdersSlice = app.getLimitMarketOrdersSlice()
			app.tpSlOrdersSlice = app.getTpSlOrdersSlice()
			app.limitMarketOrdersTable.Refresh()
			app.tpSlOrdersTable.Refresh()
		}()
	}
}

func (app *Config) getPositionsSlice() [][]any {
	var slice [][]any
	slice = append(slice, []any{"Side", "Symbol", "Size", "Entry Price", "Market Price", "Liq. Price", "PNL", "Actions"})

	for _, profile := range app.User.Profiles {
		var currentRow []any
		var positionInfoArr []user.PositionInfo
		openedPositions := app.User.CopiedTradersManager.GetOpenedPositionsByProfileID(profile.ID)
		for _, position := range openedPositions {
			positionInfoArr = append(positionInfoArr, app.getPositionInfo("linear", position.Symbol, position.OrderID, profile)...)
		}

		for _, pos := range positionInfoArr {
			var side string

			if pos.Side == "None" {
				// TODO: handle "Empty position" aka None
				continue
			} else {
				side = pos.Side
			}

			symbol := pos.Symbol
			size := fmt.Sprintf("%.2f", pos.Size)
			entryPrice := fmt.Sprintf("%.2f", pos.AvgPrice)
			marketPrice := fmt.Sprintf("%.2f", pos.MarkPrice)
			liqPrice := fmt.Sprintf("%.2f", pos.LiqPrice)
			pnl := fmt.Sprintf("%.2f", pos.CumRealisedPnl)

			currentRow = append(currentRow, side)
			currentRow = append(currentRow, symbol)
			currentRow = append(currentRow, size)
			currentRow = append(currentRow, entryPrice)
			currentRow = append(currentRow, marketPrice)
			currentRow = append(currentRow, liqPrice)
			currentRow = append(currentRow, pnl)
			currentRow = append(currentRow, pos.OrderID) // for toolbar actions

			slice = append(slice, currentRow)
		}
	}

	return slice
}

func (app *Config) getPositionsTable() *widget.Table {

	t := widget.NewTable(
		func() (int, int) {
			return len(app.positionsSlice), len(app.positionsSlice[0])
		},
		func() fyne.CanvasObject {
			lbl := widget.NewLabel("")
			toolbar := widget.NewToolbar()
			toolbar.Hide()
			return container.NewMax(lbl, toolbar)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			container := o.(*fyne.Container)
			lbl := container.Objects[0].(*widget.Label)
			toolbar := container.Objects[1].(*widget.Toolbar)

			if i.Row != 0 && i.Col == 7 {
				lbl.Hide()
				toolbar.Hidden = false

				if len(toolbar.Items) == 0 {
					toolbar.Append(widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
						// SET TP

					}))
					toolbar.Append(widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
						// SET SL

					}))
					toolbar.Append(widget.NewToolbarAction(theme.DeleteIcon(), func() {
						// CLOSE POSITION

						var err error
						var side, symbol, orderId string
						var amount float64

						// get and parse columns content
						side = app.positionsSlice[i.Row][0].(string)
						symbol = app.positionsSlice[i.Row][1].(string)
						orderId = app.positionsSlice[i.Row][i.Col].(string)

						// retrieve profile from position orderID
						openedPosition := app.User.CopiedTradersManager.GetOpenedPositionByOrderID(orderId)
						profile := app.User.ProfileManager.GetProfileByID(openedPosition.ProfileID)

						amount, err = strconv.ParseFloat(app.positionsSlice[i.Row][2].(string), 64)

						if err != nil {
							return
						}

						dialog.ShowConfirm("Close Position?", fmt.Sprintf("Confirming will close %s %s Position.", symbol, side), func(deleted bool) {
							if deleted {
								var orderData OrderData

								orderData.OrderType = utils.ORDER_MARKET
								orderData.Symbol = symbol
								orderData.Qty = amount

								if side == "Buy" {
									_, err = app.closeLongPosition(profile, orderData)
								} else {
									_, err = app.closeShortPosition(profile, orderData)
								}

								if err != nil {
									app.App.SendNotification(fyne.NewNotification(
										"⚠️ Close Position Failed",
										fmt.Sprintf("Error: %s", err.Error()),
									))
								} else {
									err = app.DB.DeleteOpenedPosition(orderId)
									if err != nil {
										app.Logger.Error(err)
									} else {
										app.User.CopiedTradersManager.DeleteOpenedPosition(orderId)
										app.Logger.Debug(fmt.Sprintf("closed %s position with order ID %s", symbol, orderId))
										app.App.SendNotification(fyne.NewNotification(
											"🔴 Successfully Closed Position",
											fmt.Sprintf("closed %s %s position", symbol, side),
										))
									}
									app.refreshCopiedTradersTab()
								}
							}
						}, app.MainWindow)
					}))
				}
			} else {
				toolbar.Hide()
				lbl.Hidden = false
				// we're just putting in textual information
				lbl.SetText(
					app.positionsSlice[i.Row][i.Col].(string))
			}
		})

	colWidths := []float32{100, 100, 100, 100, 100, 100, 100, 150}
	for i, w := range colWidths {
		t.SetColumnWidth(i, w)
	}

	return t
}

func (app *Config) refreshCopiedTradersList() {
	app.CopiedTradersTab.profilesWithTrader = app.formatCopiedTradersList(app.User.GetProfilesWithTrader())
	app.CopiedTradersList.Refresh()
}

func (app *Config) refreshPositionsTable() {
	app.positionsSlice = app.getPositionsSlice()
	app.positionsTable.Refresh()

	colWidths := []float32{100, 100, 100, 100, 100, 100, 100, 150}
	for i, w := range colWidths {
		app.positionsTable.SetColumnWidth(i, w)
	}
}

func (app *Config) refreshCopiedTradersTab() {

	var limitMarketOrdersSlice, tpSlOrdersSlice [][]any
	limitMarketOrdersSlice = append(limitMarketOrdersSlice, []any{"Symbol", "Order Type", "Side", "Avg. Filled Price", "Filled Qty", "Order Price", "Order Qty", "Order Status", "Order Time", "Order ID"})
	tpSlOrdersSlice = append(tpSlOrdersSlice, []any{"Symbol", "Order Type", "Side", "Order Value", "Order Qty", "Order Price", "Trigger Price", "Order Status", "Order Time", "Order ID"})
	app.limitMarketOrdersSlice = limitMarketOrdersSlice
	app.tpSlOrdersSlice = tpSlOrdersSlice
	app.limitMarketOrdersTable.Refresh()
	app.tpSlOrdersTable.Refresh()

	// update content based on the current selection
	app.updateOrderHistoryContent()
	app.refreshPositionsTable()

	app.CopiedTradersTab.profilesWithTrader = app.formatCopiedTradersList(app.User.ProfileManager.GetProfilesWithTrader())
	// set default height for each list item
	for i := range app.CopiedTradersTab.profilesWithTrader {
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
