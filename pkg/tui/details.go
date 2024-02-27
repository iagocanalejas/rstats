package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (app *Application) showDetailsView(raceID int64) {
	app.showingDetails = true
	race, err := app.service.GetRaceByID(raceID)
	if err != nil {
		app.errorModal(err)
		return
	}

	app.race = race

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(app.detailHeader(), 3, 0, false).
		AddItem(app.participantDetails(), 0, 1, true)

	app.App.SetRoot(flex, true)
}

func (app *Application) detailHeader() *tview.TextView {
	header := tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetTextAlign(tview.AlignLeft).
		SetText(fmt.Sprintf("%d (%s) || %s", app.race.ID, app.race.Date, app.race.Name))

	header.Box.SetBorder(true)

	return header
}

func (app *Application) participantDetails() *tview.Table {
	table := tview.NewTable().
		SetBorders(true)

	table.SetCell(0, 0, &tview.TableCell{Text: "Club Name", Align: tview.AlignCenter, Color: tcell.ColorYellow})
	table.SetCell(0, 1, &tview.TableCell{Text: "Serie", Align: tview.AlignCenter, Color: tcell.ColorYellow})
	table.SetCell(0, 2, &tview.TableCell{Text: "Lane", Align: tview.AlignCenter, Color: tcell.ColorYellow})
	for i := 0; i < int(app.race.Laps.Int64); i++ {
		colIndex := 3 + i
		table.SetCell(0, colIndex, &tview.TableCell{Text: fmt.Sprintf("Lap %d", i+1), Align: tview.AlignCenter, Color: tcell.ColorYellow})
	}
	table.SetCell(0, int(app.race.Laps.Int64)+2, &tview.TableCell{Text: "Time", Align: tview.AlignCenter, Color: tcell.ColorYellow})

	for i, participant := range app.race.Participants {
		rowIndex := i + 1
		table.SetCell(rowIndex, 0, &tview.TableCell{Text: participant.Club.Name, Align: tview.AlignLeft})
		table.SetCell(rowIndex, 1, &tview.TableCell{Text: fmt.Sprintf("%d", participant.Series.Int64), Align: tview.AlignCenter})
		if participant.Lane.Valid {
			table.SetCell(rowIndex, 2, &tview.TableCell{Text: fmt.Sprintf("%d", participant.Lane.Int64), Align: tview.AlignCenter})
		}
		for j, lap := range participant.Laps {
			colIndex := 3 + j
			table.SetCell(rowIndex, colIndex, &tview.TableCell{Text: fmt.Sprintf("%v", lap), Align: tview.AlignLeft})
		}
	}

	return table
}
