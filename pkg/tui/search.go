package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (app *Application) searchInput() *tview.InputField {
	searchInput := tview.NewInputField().
		SetLabel("Search: ").
		SetFieldWidth(30).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite)
	searchInput.Box.SetBorder(true)

	searchInput.SetChangedFunc(func(text string) {
		app.currentSearch = text

		if len(text) == 0 || (len(text) > 3 && len(text)%2 == 0) {
			app.populateList()
		}
	})

	return searchInput
}
