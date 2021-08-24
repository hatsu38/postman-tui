package gui

import (
	"fmt"
	"net/url"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type table struct {
	*tview.Table
}
type Param struct {
	Key   string
	Value string
}
type Params []Param

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
	t.SetCell(0, 0, SetTableCell(cellTxt, 1, tcell.ColorIndianRed, false))
	t.SetCell(0, 1, SetTableCell("Key", 3, tcell.ColorIndianRed, false))
	t.SetCell(0, 2, SetTableCell("Value", 3, tcell.ColorIndianRed, false))
}

func (t *table) AddParamsRow(idx int) {
	t.SetCell(idx, 0, SetTableCell(fmt.Sprint(idx), 1, tcell.ColorWhite, false))
	t.SetCell(idx, 1, SetTableCell("", 3, tcell.ColorWhite, true))
	t.SetCell(idx, 2, SetTableCell("", 3, tcell.ColorWhite, true))
}

func SetTableCell(title string, width int, color tcell.Color, selectable bool) *tview.TableCell {
	tcell := tview.NewTableCell(title)
	tcell.SetExpansion(width)
	tcell.SetAlign(tview.AlignCenter)
	tcell.SetTextColor(color)
	tcell.SetSelectable(selectable)
	tcell.SetTransparency(true)

	return tcell
}

func (t *table) GetQuery() string {
	params := t.GetParams()
	return params.ToQueryString()
}

func (t *table) GetBodyParams() string {
	params := t.GetParams()
	return params.ToBodyParams()
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

func (p Params) ToQueryString() string {
	var query string
	for i, v := range p {
		if v.Key == "" || v.Value == "" {
			continue
		}
		if i == 0 {
			query += "?"
		} else {
			query += "&"
		}
		query += fmt.Sprintf("%s=%s", v.Key, v.Value)
	}

	return query
}

func (p Params) ToBodyParams() string {
	val := url.Values{}
	for i, v := range p {
		if v.Key == "" || v.Value == "" {
			continue
		}
		if i == 0 {
			val.Set(v.Key, v.Value)
		} else {
			val.Add(v.Key, v.Value)
		}
	}

	return val.Encode()
}
