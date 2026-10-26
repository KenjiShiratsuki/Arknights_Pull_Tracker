package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func initializeState() *Log {
	log := Log{
		history: []State{},
	}
	initState := State{
		threeStarCount: 0,
		fourStarCount:  0,
		fiveStarCount:  9,
		sixStarCount:   0,
		pityCount:      0,
		totalPullCount: 0,
		threeStarNew:   0,
		fourStarNew:    0,
		fiveStarNew:    0,
		sixStarNew:     0,
		totalNew:       0,
	}
	log.history = append(log.history, initState)
	log.currentState = 0
	return &log
}

func (s State) sixStarChance() int {
	if s.pityCount < 50 {
		return 2
	}

	chance := 2 + (s.pityCount-49)*2

	if chance > 100 {
		return 100
	}

	return chance
}

type UI struct {
	totalPulls    *widget.Label
	threeStars    *widget.Label
	fourStars     *widget.Label
	fiveStars     *widget.Label
	sixStars      *widget.Label
	pityCounter   *widget.Label
	sixStarChance *widget.Label
	totalNew      *widget.Label
	threeNew      *widget.Label
	fourNew       *widget.Label
	fiveNew       *widget.Label
	sixNew        *widget.Label
}

func initializeApp(log *Log) fyne.App {
	a := app.New()
	w := a.NewWindow("Arknights Pull Tracker by PhantomTheKitsune")
	ui := &UI{
		totalPulls:    widget.NewLabel(fmt.Sprintf("Total Pulls: %d", 0)),
		threeStars:    widget.NewLabel(fmt.Sprintf("3-Stars: %d", 0)),
		fourStars:     widget.NewLabel(fmt.Sprintf("4-Stars: %d", 0)),
		fiveStars:     widget.NewLabel(fmt.Sprintf("5-Stars: %d", 0)),
		sixStars:      widget.NewLabel(fmt.Sprintf("6-Stars: %d", 0)),
		pityCounter:   widget.NewLabel(fmt.Sprintf("Pity Count: %d", 0)),
		sixStarChance: widget.NewLabel(fmt.Sprintf("6-Star Rate: %d%%", 2)),
		totalNew:      widget.NewLabel(fmt.Sprintf("(%d new!)", 0)),
		threeNew:      widget.NewLabel(fmt.Sprintf("(%d new!)", 0)),
		fourNew:       widget.NewLabel(fmt.Sprintf("(%d new!)", 0)),
		fiveNew:       widget.NewLabel(fmt.Sprintf("(%d new!)", 0)),
		sixNew:        widget.NewLabel(fmt.Sprintf("(%d new!)", 0)),
	}
	w.SetContent(container.NewGridWithColumns(4,
		container.NewVBox(
			ui.totalPulls,
			ui.threeStars,
			ui.fourStars,
			ui.fiveStars,
			ui.sixStars,
		),
		container.NewVBox(
			ui.pityCounter,
			//Three-star, not new.
			widget.NewButton("+", func() {
				newState := log.history[log.currentState]
				newState.threeStarCount++
				newState.pityCount++
				newState.totalPullCount++
				pushState(newState, log)
				updateUI(log, ui)
			},
			),
			//Four-star, not new.
			widget.NewButton("+", func() {
				newState := log.history[log.currentState]
				newState.fourStarCount++
				newState.pityCount++
				newState.totalPullCount++
				pushState(newState, log)
				updateUI(log, ui)
			},
			),
			//Five-star, not new.
			widget.NewButton("+", func() {
				newState := log.history[log.currentState]
				newState.fiveStarCount++
				newState.pityCount++
				newState.totalPullCount++
				pushState(newState, log)
				updateUI(log, ui)
			},
			),
			//Six-star, not new.
			widget.NewButton("+", func() {
				newState := log.history[log.currentState]
				newState.sixStarCount++
				newState.pityCount = 0
				newState.totalPullCount++
				pushState(newState, log)
				updateUI(log, ui)
			},
			),
		),
		container.NewVBox(
			ui.sixStarChance,
			//Three-star, new.
			widget.NewButton("+", func() {
				newState := log.history[log.currentState]
				newState.threeStarNew++
				newState.pityCount++
				newState.totalPullCount++
				newState.totalNew++
				pushState(newState, log)
				updateUI(log, ui)
			},
			),
			//Four-star, new.
			widget.NewButton("+", func() {
				newState := log.history[log.currentState]
				newState.fourStarNew++
				newState.pityCount++
				newState.totalPullCount++
				newState.totalNew++
				pushState(newState, log)
				updateUI(log, ui)
			},
			),
			//Five-star, new.
			widget.NewButton("+", func() {
				newState := log.history[log.currentState]
				newState.fiveStarNew++
				newState.pityCount++
				newState.totalPullCount++
				newState.totalNew++
				pushState(newState, log)
				updateUI(log, ui)
			},
			),
			//Six-star, new.
			widget.NewButton("+", func() {
				newState := log.history[log.currentState]
				newState.sixStarNew++
				newState.pityCount = 0
				newState.totalPullCount++
				newState.totalNew++
				pushState(newState, log)
				updateUI(log, ui)
			},
			),
		),
		container.NewVBox(
			ui.totalNew,
			ui.threeNew,
			ui.fourNew,
			ui.fiveNew,
			ui.sixNew,
		)))
	w.Show()
	return a
}
