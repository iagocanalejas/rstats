package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (app *Application) searchView() *tview.Flex {
	app.searchInput = tview.NewInputField().
		SetLabel("Search: ").
		SetFieldWidth(30).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite)

	app.searchInput.SetChangedFunc(func(text string) {
		app.currentSearch = text
	})

	legend := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetTextAlign(tview.AlignLeft).
		SetText("Filters -> year | trophy[_id] | flag[_id] | league[_id] | participant[_id]")

	searchBox := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(app.searchInput, 1, 0, true).
		AddItem(legend, 1, 0, false)
	searchBox.Box.SetBorder(true)

	return searchBox
}
