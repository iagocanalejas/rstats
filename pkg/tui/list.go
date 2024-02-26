package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (app *Application) listView() *tview.List {
	app.racesList = tview.NewList()
	app.racesList.Box.SetBorder(true)

	app.racesList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// disable default behavior fot TAB key (next list item) as we use that to change focus
		if event.Key() == tcell.KeyTab {
			return nil
		}
		return event
	})

	app.populateList()

	return app.racesList
}

func (app *Application) populateList() {
	races := app.service.SearchRaces(app.currentSearch)

	app.racesList.Clear()
	for _, race := range races {
		app.racesList.AddItem(fmt.Sprintf("%d", race.ID), race.Name, 0, nil)
	}
}
