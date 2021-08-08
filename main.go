package main

import (
	"github.com/rivo/tview"
)


func main() {
	// Flex
	app := tview.NewApplication()
	flex := tview.NewFlex().
		AddItem(tview.NewBox().SetBorder(true).SetTitle("HTTP METHODS"), 0, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("URL"), 0, 1, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Params"), 0, 2, false).
			AddItem(tview.NewBox().SetBorder(true).SetTitle("Response"), 0, 4, false), 0, 4, false)
	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}
