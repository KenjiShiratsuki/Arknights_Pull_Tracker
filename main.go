package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Arknights Pull Tracker by PhantomTheKitsune")
	t := 0
	three := 0
	four := 0
	five := 0
	six := 0
	pity := 0
	totals := widget.NewLabel(fmt.Sprintf("Pull Totals: %d", t))
	threeStar := widget.NewLabel(fmt.Sprintf("3-Star: %d", three))
	fourStar := widget.NewLabel(fmt.Sprintf("4-Star: %d", four))
	fiveStar := widget.NewLabel(fmt.Sprintf("5-Star: %d", five))
	sixStar := widget.NewLabel(fmt.Sprintf("6-Star: %d", six))
	pityCount := widget.NewLabel(fmt.Sprintf("6* Pity: %d", pity))
	w.SetContent(container.NewVBox(
		totals,
		threeStar,
		fourStar,
		fiveStar,
		sixStar,
		pityCount,
		widget.NewButton("3-Star", func() {
			three++
			pity++
			threeStar.SetText(fmt.Sprintf("3-Star: %d", three))
			pityCount.SetText(fmt.Sprintf("6* Pity: %d", pity))
		}),
		widget.NewButton("4-Star", func() {
			four++
			pity++
			fourStar.SetText(fmt.Sprintf("4-Star: %d", four))
			pityCount.SetText(fmt.Sprintf("6* Pity: %d", pity))
		}),
		widget.NewButton("5-Star", func() {
			five++
			pity++
			fiveStar.SetText(fmt.Sprintf("5-Star: %d", five))
			pityCount.SetText(fmt.Sprintf("6* Pity: %d", pity))
		}),
		widget.NewButton("6-Star", func() {
			six++
			pity = 0
			sixStar.SetText(fmt.Sprintf("6-Star: %d", six))
			pityCount.SetText(fmt.Sprintf("6* Pity: %d", pity))
		}),
	))

	w.ShowAndRun()
}
