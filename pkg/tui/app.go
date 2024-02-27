package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/iagocanalejas/regatas/internal/service"
	"github.com/iagocanalejas/regatas/internal/types/races"
	"github.com/rivo/tview"
)

type Application struct {
	App *tview.Application

	service *service.Service

	race           *races.Race
	races          []races.Race
	currentSearch  string // current search keywords
	hasError       bool   // if the error modal is showing or not
	showingDetails bool   // if the details view is in display

	flex        *tview.Flex
	searchInput *tview.InputField
	racesList   *tview.List
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
		AddItem(app.searchView(), 4, 0, false).
		AddItem(app.listView(), 0, 1, true).
		AddItem(app.bottomLegend(), 3, 0, false)
}

func (app *Application) setupListeners() {
	app.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if app.hasError {
			return event
		}
		switch event.Key() {
		case tcell.KeyEnter:
			if app.searchInput.HasFocus() {
				// configure search on <CR> press
				app.populateList()
			}
		case tcell.KeyTab:
			app.nextFocus()
		case tcell.KeyEsc:
			if app.showingDetails {
				app.App.SetRoot(app.flex, true)
				app.showingDetails = false
			} else {
				app.App.Stop()
			}
		}
		return event
	})
}

func (app *Application) nextFocus() {
	if app.searchInput.HasFocus() {
		app.App.SetFocus(app.racesList)
	} else {
		app.App.SetFocus(app.searchInput)
	}
}

func (app *Application) errorModal(err error) {
	if app.hasError {
		return
	}

	modal := tview.NewModal().
		SetText(err.Error()).
		AddButtons([]string{"Continue"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.App.SetRoot(app.flex, true).SetFocus(app.flex)
			app.hasError = false
		})

	modal.
		SetBackgroundColor(tcell.ColorDarkRed).
		SetTextColor(tcell.ColorYellow).
		SetBorder(true).
		SetBorderColor(tcell.ColorWhite).
		SetBorderPadding(2, 2, 2, 2)

	app.App.SetRoot(modal, true).SetFocus(modal)
	app.hasError = true
}
