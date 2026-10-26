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
		fiveStarCount:  0,
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
	pityCounter   *widget.Label
	sixStarChance *widget.Label
	totalNew      *widget.Label
	threeStars    *widget.Button
	fourStars     *widget.Button
	fiveStars     *widget.Button
	sixStars      *widget.Button
	threeNew      *widget.Button
	fourNew       *widget.Button
	fiveNew       *widget.Button
	sixNew        *widget.Button
}

func initializeApp(log *Log) fyne.App {
	a := app.New()
	w := a.NewWindow("Arknights Pull Tracker by PhantomTheKitsune")
	ui := &UI{
		totalPulls:    widget.NewLabel(fmt.Sprintf("Total Pulls: %d", 0)),
		pityCounter:   widget.NewLabel(fmt.Sprintf("Pity Count: %d", 0)),
		sixStarChance: widget.NewLabel(fmt.Sprintf("6-Star Rate: %d%%", 2)),
		totalNew:      widget.NewLabel(fmt.Sprintf("(%d new!)", 0)),
	}
	ui.sixStars = widget.NewButton("6-Stars: 0", func() {
		newState := log.history[log.currentState]
		newState.sixStarCount++
		newState.pityCount = 0
		newState.totalPullCount++
		pushState(newState, log)
		updateUI(log, ui)
	})
	ui.sixNew = widget.NewButton("(0 new!)", func() {
		newState := log.history[log.currentState]
		newState.sixStarCount++
		newState.sixStarNew++
		newState.pityCount = 0
		newState.totalPullCount++
		newState.totalNew++
		pushState(newState, log)
		updateUI(log, ui)
	})
	ui.fiveStars = widget.NewButton("5-Stars: 0", func() {
		newState := log.history[log.currentState]
		newState.fiveStarCount++
		newState.pityCount++
		newState.totalPullCount++
		pushState(newState, log)
		updateUI(log, ui)
	})
	ui.fiveNew = widget.NewButton("(0 new!)", func() {
		newState := log.history[log.currentState]
		newState.fiveStarCount++
		newState.fiveStarNew++
		newState.pityCount++
		newState.totalPullCount++
		newState.totalNew++
		pushState(newState, log)
		updateUI(log, ui)
	})
	ui.fourStars = widget.NewButton("4-Stars: 0", func() {
		newState := log.history[log.currentState]
		newState.fourStarCount++
		newState.pityCount++
		newState.totalPullCount++
		pushState(newState, log)
		updateUI(log, ui)
	})
	ui.fourNew = widget.NewButton("(0 new!)", func() {
		newState := log.history[log.currentState]
		newState.fourStarCount++
		newState.fourStarNew++
		newState.pityCount++
		newState.totalPullCount++
		newState.totalNew++
		pushState(newState, log)
		updateUI(log, ui)
	})
	ui.threeStars = widget.NewButton("3-Stars: 0", func() {
		newState := log.history[log.currentState]
		newState.threeStarCount++
		newState.pityCount++
		newState.totalPullCount++
		pushState(newState, log)
		updateUI(log, ui)
	})
	ui.threeNew = widget.NewButton("(0 new!)", func() {
		newState := log.history[log.currentState]
		newState.threeStarCount++
		newState.threeStarNew++
		newState.pityCount++
		newState.totalPullCount++
		newState.totalNew++
		pushState(newState, log)
		updateUI(log, ui)
	})
	w.SetContent(container.NewVBox(
		container.NewHBox(ui.pityCounter, ui.sixStarChance),
		container.NewHBox(ui.sixStars, ui.sixNew),
		container.NewHBox(ui.fiveStars, ui.fiveNew),
		container.NewHBox(ui.fourStars, ui.fourNew),
		container.NewHBox(ui.threeStars, ui.threeNew),
		container.NewHBox(ui.totalPulls, ui.totalNew),
	))
	updateUI(log, ui)
	w.Show()
	return a
}
