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
	OrdersSlice          [][]any
	OrdersTable          *widget.Table
	CopiedTradersList    *widget.List
	orders               []user.Order
	positionsLabel       *widget.Label
	positionsContainer   *fyne.Container
	selectedCopiedTrader *user.Profile
	tradersList          []string
}

func (app *Config) copiedTradersTab() *fyne.Container {
	split := app.getCopiedTraders()

	max := container.NewMax(split)

	return max
}

func (app *Config) getCopiedTraders() *container.Split {

	/* LEFT SPLIT */

	// get the copied traders list (profiles with a trader)
	app.CopiedTradersList = widget.NewList(
		func() int {
			return len(app.CopiedTradersTab.tradersList)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.AccountIcon()), widget.NewLabel("Template Object"), widget.NewToolbar())
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			// get the current trader...
			trader := app.CopiedTradersTab.tradersList[id]
			traderID := strings.Split(trader, "\n")[1]

			// display profile title and trader nickname
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(trader)

			toolbar := item.(*fyne.Container).Objects[2].(*widget.Toolbar)
			if len(toolbar.Items) == 0 {
				toolbar.Append(widget.NewToolbarAction(theme.VisibilityIcon(), func() {
					app.traderDialog(*t)
				}))
				toolbar.Append(widget.NewToolbarAction(theme.DeleteIcon(), func() {
					dialog.ShowConfirm("Stop copying?", "Doing that will NOT close positions copied from this trader.", func(deleted bool) {
						err := app.stopCopyingTrader(*t)
						if err != nil {
							app.Logger.Error(err)
						}
					}, app.MainWindow)
				}))
			}
		},
	)

	// set default height for each list item
	for i := range app.CopiedTradersTab.tradersList {
		app.CopiedTradersList.SetItemHeight(i, 50)
	}

	/* RIGTH SPLIT */

	// get the top content

	// get the refresh toolbar button
	topRightToolbar := widget.NewToolbar(widget.NewToolbarSpacer(), widget.NewToolbarAction(
		theme.ViewRefreshIcon(), func() {
			app.refreshCopiedTradersTab()
		}))

	// wrap the top content in a container
	top := container.NewHBox(
		topRightToolbar,
	)

	// get the main content

	// get the order history table
	var slice [][]any
	slice = append(slice, []any{"Symbol", "Order ID", "Status", "Qty", "Price", "Side", "Leverage", "Created Time"})
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
	colWidths := []float32{200, 100, 200, 100, 100, 100, 100, 200}
	for i, w := range colWidths {
		app.OrdersTable.SetColumnWidth(i, w)
	}

	// get the orders tab
	ordersTab := container.NewVScroll(app.OrdersTable)

	// get the positions tab
	app.CopiedTradersTab.positionsLabel = widget.NewLabel("Select a trader")

	// get the positions container
	app.CopiedTradersTab.positionsContainer = container.NewVBox()

	positionsTab := container.NewBorder(
		app.CopiedTradersTab.positionsLabel,
		nil,
		nil,
		nil,
		container.NewVScroll(app.CopiedTradersTab.positionsContainer))

	// get the tabs container
	tabs := container.NewAppTabs(
		container.NewTabItem("Order History", ordersTab),
		container.NewTabItem("Positions", positionsTab))
	tabs.SetTabLocation(container.TabLocationTop)

	// when a trader is selected...
	app.CopiedTradersList.OnSelected = func(id widget.ListItemID) {
		// ... get the selected item (profile)
		trader := app.CopiedTradersTab.tradersList[id]
		traderID := strings.Split(trader, "\n")[1]
		profile := app.User.ProfileManager.GetProfileByTraderID(traderID)
		app.CopiedTradersTab.selectedCopiedTrader = profile

		// update order history
		app.updateOrderHistoryContent(profile)

		// update positions
		app.CopiedTradersTab.positionsLabel.SetText("Fetching positions...")
		app.updatePositionsContent(profile)
	}

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

func (app *Config) getOrderSlice() [][]any {
	var slice [][]any

	slice = append(slice, []any{"Symbol", "Order ID", "Status", "Qty", "Price", "Side", "Leverage", "Created Time"})

	for _, x := range app.CopiedTradersTab.orders {
		var currentRow []any

		if x.OrderStatus == "Filled" {

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

func (app *Config) getOrders(profile *user.Profile) []user.Order {
	var allOrders []user.Order

	if profile != nil {
		if profile.BybitApiKey != "" {
			byBitOrders := app.getOrderHistory("spot", *profile)
			if byBitOrders != nil {
				allOrders = append(allOrders, byBitOrders...)
			}
		}
	}

	app.CopiedTradersTab.orders = allOrders

	return allOrders
}

func (app *Config) updateOrderHistoryContent(profile *user.Profile) {
	go func() {
		app.getOrders(profile)
		app.OrdersSlice = app.getOrderSlice()
		app.OrdersTable.Refresh()
	}()
}

func (app *Config) updatePositionsContent(profile *user.Profile) {
	go func() {
		positionInfoArr := app.getPositionInfo("linear", *profile)
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

			app.CopiedTradersTab.positionsContainer.RemoveAll()
			app.CopiedTradersTab.positionsContainer.Add(widget.NewRichTextFromMarkdown(markdownText))
			app.CopiedTradersTab.positionsContainer.Add(container.NewGridWithColumns(3,
				widget.NewButtonWithIcon("Set TP", theme.DocumentCreateIcon(), func() {

				}),
				widget.NewButtonWithIcon("Set SL", theme.DocumentCreateIcon(), func() {

				}),
				widget.NewButtonWithIcon("Close", theme.DeleteIcon(), func() {

				}),
			))
			app.CopiedTradersTab.positionsContainer.Add(widget.NewSeparator())
		}
	}()
}

func (app *Config) refreshCopiedTradersTab() {
	go func() {
		if app.CopiedTradersTab.selectedCopiedTrader == nil {
			// reset orders table
			var slice [][]any
			slice = append(slice, []any{"Symbol", "Order ID", "Status", "Qty", "Price", "Side", "Leverage", "Created Time"})
			app.OrdersSlice = slice
			app.OrdersTable.Refresh()

			// reset positions
			app.CopiedTradersTab.positionsLabel.SetText("Select a profile")
			app.CopiedTradersTab.positionsContainer.RemoveAll()
		} else {
			// update content based on the current selection
			app.updateOrderHistoryContent(app.CopiedTradersTab.selectedCopiedTrader)
			app.updatePositionsContent(app.CopiedTradersTab.selectedCopiedTrader)
		}
	}()
}
