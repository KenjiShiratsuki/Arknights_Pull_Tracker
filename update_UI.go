package main

import (
	"fmt"
)

func updateUI(log *Log, ui *UI) {
	state := log.history[log.currentState]
	ui.totalPulls.SetText(fmt.Sprintf("Total Pulls: %d", state.totalPullCount))
	ui.threeStars.SetText(fmt.Sprintf("3-Stars: %d", state.threeStarCount))
	ui.fourStars.SetText(fmt.Sprintf("4-Stars: %d", state.fourStarCount))
	ui.fiveStars.SetText(fmt.Sprintf("5-Stars: %d", state.fiveStarCount))
	ui.sixStars.SetText(fmt.Sprintf("6-Stars: %d", state.sixStarCount))
	ui.pityCounter.SetText(fmt.Sprintf("6-Star Pity: %d", state.pityCount))
	ui.threeNew.SetText(fmt.Sprintf("(%d new!)", state.threeStarNew))
	ui.fourNew.SetText(fmt.Sprintf("(%d new!)", state.fourStarNew))
	ui.fiveNew.SetText(fmt.Sprintf("(%d new!)", state.fiveStarNew))
	ui.sixNew.SetText(fmt.Sprintf("(%d new!)", state.sixStarNew))
	ui.sixStarChance.SetText(fmt.Sprintf("6* rate %d%%", state.sixStarChance()))
}
