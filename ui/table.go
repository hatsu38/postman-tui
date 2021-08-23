package ui

import (
	"fmt"

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

func (t *table) SetTableCells(g *Gui, title string) {
	// 選択された状態でEnterされたとき
	t.SetSelectedFunc(func(row, column int) {
		g.NewInputModal(t)
	})
	t.AddTableHeader(title)
	t.AddParamsRow(1)
}

func (t *table) AddTableHeader(cellTxt string) {
	t.SetCell(0, 0, t.SetTableCell(cellTxt, 1, tcell.ColorIndianRed, false))
	t.SetCell(0, 1, t.SetTableCell("Key", 3, tcell.ColorIndianRed, false))
	t.SetCell(0, 2, t.SetTableCell("Value", 3, tcell.ColorIndianRed, false))
}

func (t *table) AddParamsRow(idx int) {
	t.SetCell(idx, 0, t.SetTableCell(fmt.Sprint(idx), 1, tcell.ColorWhite, false))
	t.SetCell(idx, 1, t.SetTableCell("", 3, tcell.ColorWhite, true))
	t.SetCell(idx, 2, t.SetTableCell("", 3, tcell.ColorWhite, true))
}

func (t *table) SetTableCell(title string, width int, color tcell.Color, selectable bool) *tview.TableCell {
	tcell := tview.NewTableCell(title)
	tcell.SetExpansion(width)
	tcell.SetAlign(tview.AlignCenter)
	tcell.SetTextColor(color)
	tcell.SetSelectable(selectable)
	tcell.SetTransparency(true)

	return tcell
}

func (t *table) GetParams() Params {
	var params Params

	rows := t.GetRowCount()
	for r := 1; r < rows; r++ {
		key := t.GetCell(r, 1).Text
		value := t.GetCell(r, 2).Text
		param := Param{
			Key:   key,
			Value: value,
		}
		params = append(params, param)
	}

	return params
}
