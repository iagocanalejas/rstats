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
		switch event.Key() {
		// disable default behavior fot TAB key (next list item) as we use that to change focus
		case tcell.KeyTab:
			return nil
		case tcell.KeyEnter:
			selectedItem := app.racesList.GetCurrentItem()
			app.showDetailsView(app.races[selectedItem].ID)
		}
		return event
	})

	app.populateList()

	return app.racesList
}

func (app *Application) populateList() {
	app.racesList.Clear()

	races, err := app.service.SearchRaces(app.currentSearch)
	if err != nil {
		app.errorModal(err)
		return
	}

	app.races = races
	for _, race := range races {
		app.racesList.AddItem(fmt.Sprintf("%d (%s)", race.ID, race.Date), race.Name, 0, nil)
	}
}
