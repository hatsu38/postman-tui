package ui

import (
	"net/http"
	"fmt"
	"os"
	"io"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Gui struct {
	UrlField *tview.InputField
	App   *tview.Application
	Pages *tview.Pages
}

type Param struct {
	Key string
	Value string
}
type Params []Param

func New() *Gui {
	g := &Gui{
		UrlField: NewForm(" Request URL: ", "https://httpbin.org/get"),
		App:   NewApplication(),
		Pages: tview.NewPages(),
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

func (g *Gui) GetRequestUrl(tableView *tview.Table) string {
	field := g.UrlField
	urlText := field.GetText()
	params := g.GetParams(tableView)
	query := g.GetParamsText(params)

	return urlText + query
}

func (g *Gui) HttpRequest(url string) *http.Response {
	req, _ := http.NewRequest("GET", url, nil)
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
	textView := g.TextView("Response")
	inputUrlField := g.UrlField
	tableView := g.Table()

	inputUrlField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			url := g.GetRequestUrl(tableView)
			resp := g.HttpRequest(url)
			defer resp.Body.Close()

			body := g.ParseResponse(resp)

			textView.SetText(body)
		case tcell.KeyTab:
			g.ToTableFocus(tableView)
		}
	})

	tableView.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyTab:
			g.ToUrlFieldFocus(tableView)
		}
	})

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(inputUrlField, 0, 1, true)
	flex.AddItem(tableView, 0, 3, false)
	flex.AddItem(textView, 0, 5, false)

	g.Pages.AddAndSwitchToPage("main", flex, true)

	if err := app.SetRoot(g.Pages, true).Run(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (g *Gui) TextView(title string) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetTitle(title)
	textView.SetBorder(true)
	textView.SetScrollable(true)
	textView.SetTextColor(tcell.ColorGreen)

	return textView
}

func (g *Gui) ToTableFocus(tableView *tview.Table) {
	g.App.SetFocus(tableView)
	tableView.SetSelectable(true, true)
	tableView.SetBordersColor(tcell.ColorGreen)

	g.UrlField.SetBorderColor(tcell.ColorWhite)
}

func (g *Gui) ToUrlFieldFocus(tableView *tview.Table) {
	urlField := g.UrlField
	g.App.SetFocus(urlField)
	tableView.SetSelectable(false, false)
	tableView.SetBordersColor(tcell.ColorWhite)

	urlField.SetBorderColor(tcell.ColorGreen)
}

func (g *Gui) Input(tableView *tview.Table) {
	row, col := tableView.GetSelection()
	cell := tableView.GetCell(row, col)
	cell.SetTextColor(tcell.ColorWhite)

	text := cell.Text
	labelCell := tableView.GetCell(0, col)
	labelIndexCell := tableView.GetCell(row, 0)
	label := fmt.Sprintf(" %s %s: ", labelCell.Text, labelIndexCell.Text)
	input := NewForm(label, text)
	input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			txt := input.GetText()
			cell.Text = txt
			if txt != "" {
				g.AddParamsRow(tableView, row + 1)
			}
			g.Pages.RemovePage("input")
			g.ToTableFocus(tableView)
		}
	})

	g.Pages.AddAndSwitchToPage("input", g.Modal(input, 0, 3), true).ShowPage("main")
}

func (g *Gui) Modal(p tview.Primitive, width, height int) tview.Primitive {
	grid := tview.NewGrid()
	grid.SetColumns(0, width, 0)
	grid.SetRows(0, height, 0)
	grid.AddItem(p, 1, 1, 1, 1, 1, 1, true)

	return grid
}

func (g *Gui) Table() *tview.Table {
	table := tview.NewTable()
	table.SetBorders(true)
	g.AddTableHeader(table)
	g.AddParamsRow(table, 1)
	table.SetFixed(1, 3)
	// 選択された状態でEnterされたとき
	table.SetSelectedFunc(func(row int, column int) {
		g.Input(table)
	})

	return table
}

func (g *Gui) AddTableHeader(table *tview.Table) {
	table.SetCell(0, 0, g.TableCell("Params", 1, tcell.ColorIndianRed, false))
	table.SetCell(0, 1, g.TableCell("Key", 2, tcell.ColorIndianRed, false))
	table.SetCell(0, 2, g.TableCell("Value", 2, tcell.ColorIndianRed, false))
}

func (g *Gui) AddParamsRow(table *tview.Table, idx int) {
	table.SetCell(idx, 0, g.TableCell(fmt.Sprint(idx), 1, tcell.ColorWhite, false))
	table.SetCell(idx, 1, g.TableCell("", 2, tcell.ColorWhite, true))
	table.SetCell(idx, 2, g.TableCell("", 2, tcell.ColorWhite, true))
}

func (g *Gui) TableCell(title string, width int, color tcell.Color, selectable bool) *tview.TableCell {
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