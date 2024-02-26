package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/iagocanalejas/regatas/internal/service"
	"github.com/rivo/tview"
)

type Application struct {
	App *tview.Application

	service *service.Service

	flex          *tview.Flex
	racesList     *tview.List
	currentSearch string
}

func BuildApp() *Application {
	app := &Application{
		App:           tview.NewApplication().EnableMouse(true),
		service:       service.Init(),
		currentSearch: "",
	}

	app.setupListeners()
	app.initFlex()

	app.App.SetRoot(app.flex, true)

	return app
}

func (app *Application) initFlex() {
	app.flex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(app.searchInput(), 0, 1, false).
		AddItem(app.listView(), 0, 15, true).
		AddItem(app.bottomLegend(), 3, 1, false)
}

func (app *Application) setupListeners() {
	app.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			// configure search on <CR> press
			app.populateList()
		case tcell.KeyTab:
			app.nextFocus()
		case tcell.KeyEsc:
			app.App.Stop()
		}
		return event
	})
}

func (app *Application) nextFocus() {
	if app.flex.GetItem(0).HasFocus() {
		app.App.SetFocus(app.flex.GetItem(1))
	} else {
		app.App.SetFocus(app.flex.GetItem(0))
	}
}
