package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type table struct {
	*tview.Table
}

func newTable() *table {
	table := &table{
		tview.NewTable(),
	}
	table.SetBorders(true)
	table.SetFixed(1, 3)

	return table
}

func (t *table) setFunc(g *Gui) {
	t.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyTab:
			g.ToFocus()
		}
	})
}
