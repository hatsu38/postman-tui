package ui

import (
	"net/http"
	"net/url"
	"strings"
	"fmt"
	"os"
	"io"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Gui struct {
	App   *tview.Application
	Pages *tview.Pages
	UrlField *tview.InputField
	ParamsTable *tview.Table
	BodyTable *tview.Table
	ResTextView *tview.TextView
	HTTPTextView *tview.TextView
}

type Param struct {
	Key string
	Value string
}
type Params []Param

func New() *Gui {
	g := &Gui{
		App:   NewApplication(),
		Pages: tview.NewPages(),
		UrlField: NewForm(" Request URL: ", "https://httpbin.org/get"),
		ParamsTable: NewTable(),
		BodyTable: NewTable(),
		ResTextView: NewTextView(" Response ", ""),
		HTTPTextView: NewTextView(" HTTP Method ", "GET"),
	}
	return g
}

func NewApplication() *tview.Application {
	return tview.NewApplication().EnableMouse(true)
}

func NewForm(label, text string) *tview.InputField {
	field := tview.NewInputField()
	field.SetLabel(label)
	field.SetFieldBackgroundColor(tcell.ColorBlack)
	field.SetBorder(true)
	field.SetBorderColor(tcell.ColorGreen)
	field.SetLabelColor(tcell.ColorIndianRed)
	field.SetText(text)

	return field
}

func NewTable() *tview.Table {
	table := tview.NewTable()
	table.SetBorders(true)
	table.SetFixed(1, 3)

	return table
}

func NewTextView(title, text string) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetTitle(title)
	textView.SetBorder(true)
	textView.SetScrollable(true)
	textView.SetText(text)
	textView.SetTextColor(tcell.ColorGreen)
	textView.SetToggleHighlights(true)

	return textView
}

func (g *Gui) GetRequestUrl() string {
	field := g.UrlField
	urlText := field.GetText()
	params := g.GetParams(g.ParamsTable)
	query := g.GetParamsText(params)

	return urlText + query
}

func (g *Gui) HttpRequest(url string) *http.Response {

	bodyParams := g.GetParams(g.BodyTable)
	value := g.GetBodyParamsText(bodyParams)

	method := g.HTTPTextView.GetText(true)
	req, _ := http.NewRequest(method, url, strings.NewReader(value))
	if method != "GET" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}

	client := new(http.Client)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		g.App.Stop()
		os.Exit(1)
	}
	return resp
}

func (g *Gui) ParseResponse(resp *http.Response) string {
	body, respErr := io.ReadAll(resp.Body)
	if respErr != nil {
		fmt.Fprintln(os.Stderr, respErr)
		g.App.Stop()
		os.Exit(1)
	}

	return " " + string(body) + " "
}

func (g *Gui) Run(i interface{}) error {
	app := g.App
	resTextView := g.ResTextView
	httpTextView := g.HTTPTextView
	inputUrlField := g.UrlField
	paramsTable := g.ParamsTable
	bodyTable := g.BodyTable

	g.SetTableCells(paramsTable, "Query Params")
	g.SetTableCells(bodyTable, "Request Body")

	httpTextView.SetTextAlign(tview.AlignCenter)
	httpTextView.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			g.NewHTTPListModal()
		case tcell.KeyTab:
			g.ToFocus()
		}
	})

	inputUrlField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			url := g.GetRequestUrl()
			resp := g.HttpRequest(url)
			defer resp.Body.Close()

			body := g.ParseResponse(resp)
			resTextView.SetText(body)
		case tcell.KeyTab:
			g.ToFocus()
		}
	})

	paramsTable.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyTab:
			g.ToFocus()
		}
	})

	bodyTable.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyTab:
			g.ToFocus()
		}
	})

	httpFlex := tview.NewFlex()
	httpFlex.SetDirection(tview.FlexColumn)
	httpFlex.AddItem(g.HTTPTextView, 0, 1, false)
	httpFlex.AddItem(inputUrlField, 0, 9, true)

	requestFlex := tview.NewFlex()
	requestFlex.SetDirection(tview.FlexRow)
	requestFlex.AddItem(httpFlex, 0, 1, true)
	requestFlex.AddItem(paramsTable, 0, 5, false)
	requestFlex.AddItem(bodyTable, 0, 5, false)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexColumn)
	flex.AddItem(requestFlex, 0, 5, true)
	flex.AddItem(resTextView, 0, 3, false)

	g.Pages.AddAndSwitchToPage("main", flex, true)

	if err := app.SetRoot(g.Pages, true).Run(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (g *Gui) ToParamsTableFocus() {
	g.App.SetFocus(g.ParamsTable)
	g.ParamsTable.SetSelectable(true, true)
	g.HTTPTextView.SetBorderColor(tcell.ColorWhite)
	g.UrlField.SetBorderColor(tcell.ColorWhite)
	g.BodyTable.SetBorderColor(tcell.ColorWhite)

	g.ParamsTable.SetBordersColor(tcell.ColorGreen)
}

func (g *Gui) ToUrlFieldFocus() {
	urlField := g.UrlField
	g.App.SetFocus(urlField)
	g.ParamsTable.SetBordersColor(tcell.ColorWhite)
	g.HTTPTextView.SetBorderColor(tcell.ColorWhite)
	g.BodyTable.SetBorderColor(tcell.ColorWhite)

	urlField.SetBorderColor(tcell.ColorGreen)
}

func (g *Gui) ToHTTPFieldFocus() {
	g.App.SetFocus(g.HTTPTextView)
	g.BodyTable.SetSelectable(false, false)
	g.ParamsTable.SetBordersColor(tcell.ColorWhite)
	g.BodyTable.SetBordersColor(tcell.ColorWhite)
	g.UrlField.SetBorderColor(tcell.ColorWhite)

	g.HTTPTextView.SetBorderColor(tcell.ColorGreen)
}

func (g *Gui) ToBodyTable() {
	g.App.SetFocus(g.BodyTable)
	g.ParamsTable.SetSelectable(false, false)
	g.BodyTable.SetSelectable(true, true)
	g.ParamsTable.SetBordersColor(tcell.ColorWhite)
	g.UrlField.SetBorderColor(tcell.ColorWhite)
	g.HTTPTextView.SetBorderColor(tcell.ColorWhite)

	g.BodyTable.SetBordersColor(tcell.ColorGreen)
}

func (g *Gui) ToFocus() {
	primitive := g.App.GetFocus()

	switch primitive {
	case g.UrlField:
		g.ToParamsTableFocus()
	case g.ParamsTable:
		g.ToBodyTable()
	case g.BodyTable:
		g.ToHTTPFieldFocus()
	case g.HTTPTextView:
		g.ToUrlFieldFocus()
	}
}

func (g *Gui) NewInputModal(table *tview.Table) {
	row, col := table.GetSelection()
	cell := table.GetCell(row, col)
	cell.SetTextColor(tcell.ColorWhite)

	text := cell.Text
	labelCell := table.GetCell(0, col)
	labelIndexCell := table.GetCell(row, 0)
	tableTitle := table.GetCell(0, 0)
	label := fmt.Sprintf(" %s %s %s: ", tableTitle.Text, labelCell.Text, labelIndexCell.Text)
	input := NewForm(label, text)
	input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			txt := input.GetText()
			cell.Text = txt
			if txt != "" {
				g.AddParamsRow(table, row + 1)
			}
			g.Pages.RemovePage("input")
			g.App.SetFocus(table)
		}
	})

	g.Pages.AddAndSwitchToPage("input", g.Modal(input, 0, 3), true).ShowPage("main")
}

func (g *Gui) NewHTTPListModal() {
	list := tview.NewList()
	list.SetBorder(true)
	list.SetBorderColor(tcell.ColorGreen)
	list.SetTitle(" HTTP Methods")
	list.AddItem("GET", "", 'a', nil)
	list.AddItem("POST", "", 'b', nil)
	list.AddItem("PUT", "", 'c', nil)
	list.AddItem("PATCH", "", 'd', nil)
	list.AddItem("DELETE", "", 'e', nil)

	txt := g.HTTPTextView.GetText(true)
	indices := list.FindItems(txt, "", true, true)
	list.SetCurrentItem(indices[0])

	list.SetSelectedFunc(func(idx int, mainTxt, subtxt string, key rune) {
		g.HTTPTextView.SetText(mainTxt)
		g.Pages.RemovePage("list")
		g.ToHTTPFieldFocus()
	})

	g.Pages.AddAndSwitchToPage("list", g.Modal(list, 40, 13), true).ShowPage("main")
}

func (g *Gui) Modal(p tview.Primitive, width, height int) tview.Primitive {
	grid := tview.NewGrid()
	grid.SetColumns(0, width, 0)
	grid.SetRows(0, height, 0)
	grid.AddItem(p, 1, 1, 1, 1, 1, 1, true)

	return grid
}

func (g *Gui) SetTableCells(table *tview.Table, title string) {
	// 選択された状態でEnterされたとき
	table.SetSelectedFunc(func(row, column int) {
		g.NewInputModal(table)
	})
	g.AddTableHeader(table, title)
	g.AddParamsRow(table, 1)
}

func (g *Gui) AddTableHeader(table *tview.Table, cellTxt string) {
	table.SetCell(0, 0, g.SetTableCell(cellTxt, 1, tcell.ColorIndianRed, false))
	table.SetCell(0, 1, g.SetTableCell("Key", 3, tcell.ColorIndianRed, false))
	table.SetCell(0, 2, g.SetTableCell("Value", 3, tcell.ColorIndianRed, false))
}

func (g *Gui) AddParamsRow(table *tview.Table, idx int) {
	table.SetCell(idx, 0, g.SetTableCell(fmt.Sprint(idx), 1, tcell.ColorWhite, false))
	table.SetCell(idx, 1, g.SetTableCell("", 3, tcell.ColorWhite, true))
	table.SetCell(idx, 2, g.SetTableCell("", 3, tcell.ColorWhite, true))
}

func (g *Gui) SetTableCell(title string, width int, color tcell.Color, selectable bool) *tview.TableCell {
	tcell := tview.NewTableCell(title)
	tcell.SetExpansion(width)
	tcell.SetAlign(tview.AlignCenter)
	tcell.SetTextColor(color)
	tcell.SetSelectable(selectable)
	tcell.SetTransparency(true)

	return tcell
}

func (g *Gui) GetParams(table *tview.Table) Params {
	var params Params

	rows := table.GetRowCount()
	for r := 1; r < rows; r++ {
		key := table.GetCell(r, 1).Text
		value := table.GetCell(r, 2).Text
		param := Param {
			Key: key,
			Value: value,
		}
		params = append(params, param)
	}

	return params
}

func (g *Gui) GetParamsText(params Params) string {
	var query string
	for i, v := range params {
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

func (g *Gui) GetBodyParamsText(params Params) string {
	val := url.Values{}
	for i, v := range params {
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
