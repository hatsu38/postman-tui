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
		UrlField: Form(" Request URL: ", "https://httpbin.org/get"),
		App:   tview.NewApplication(),
		Pages: tview.NewPages(),
	}
	return g
}

func Form(label, text string) *tview.InputField {
	field := tview.NewInputField()
	field.SetLabel(label)
	field.SetFieldTextColor(tcell.ColorWhite)
	field.SetLabelColor(tcell.ColorBlue)
	field.SetFieldBackgroundColor(tcell.ColorPaleVioletRed)
	field.SetBorder(true)
	field.SetText(text)

	return field
}


func (g *Gui) Run(i interface{}) error {
	app := g.App
	textView := g.TextView("Response")
	inputUrlField := g.UrlField
	tableView := g.Table()

	inputUrlField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			text := inputUrlField.GetText()
			params := g.GetParams(tableView)
			query := g.GetParamsText(params)
			reqUrl := text + query
			req, _ := http.NewRequest("GET", reqUrl, nil)
			client := new(http.Client)

			resp, err := client.Do(req)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				app.Stop()
				os.Exit(1)
			}
			defer resp.Body.Close()

			body, respErr := io.ReadAll(resp.Body)
			if respErr != nil {
				fmt.Fprintln(os.Stderr, err)
				app.Stop()
				os.Exit(1)
			}
			toFixBody := " " + string(body) + " "
			textView.SetText(toFixBody)
		case tcell.KeyTab:
			g.ToTableFocus(tableView)
			inputUrlField.SetFieldBackgroundColor(tcell.ColorGray)
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
	flex.AddItem(tableView, 0, 1, false)
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

	return textView
}

func (g *Gui) ToTableFocus(tableView *tview.Table) {
	g.App.SetFocus(tableView)
	tableView.SetSelectable(true, true)
	tableView.SetBordersColor(tcell.ColorPaleVioletRed)
	g.UrlField.SetFieldBackgroundColor(tcell.ColorGray)
}

func (g *Gui) ToUrlFieldFocus(tableView *tview.Table) {
	urlField := g.UrlField
	g.App.SetFocus(urlField)
	tableView.SetSelectable(false, false)
	tableView.SetBordersColor(tcell.ColorPaleVioletRed)
	urlField.SetFieldBackgroundColor(tcell.ColorPaleVioletRed)
	tableView.SetBordersColor(tcell.ColorWhite)
}

func (g *Gui) Input(tableView *tview.Table, cell *tview.TableCell) {
	text := cell.Text
	input := Form(" params", text)
	input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			cell.Text = input.GetText()
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
	grid.AddItem(p, 1, 1, 1, 1, 0, 0, true)

	return grid
}

func (g *Gui) Table() *tview.Table {
	table := tview.NewTable()
	table.SetBorders(true)
	table.SetCell(0, 0, g.TableCell("Params", 1, tcell.ColorYellow, false))
	table.SetCell(0, 1, g.TableCell("Key", 2, tcell.ColorYellow, false))
	table.SetCell(0, 2, g.TableCell("Value", 2, tcell.ColorYellow, false))
	table.SetCell(1, 0, g.TableCell("1", 1, tcell.ColorWhite, false))
	table.SetCell(1, 1, g.TableCell("", 2, tcell.ColorWhite, true))
	table.SetCell(1, 2, g.TableCell("", 2, tcell.ColorWhite, true))
	// 選択された状態でEnterされたとき
	table.SetSelectedFunc(func(row int, column int) {
		cell := table.GetCell(row, column)
		cell.SetTextColor(tcell.ColorWhite)

		g.Input(table, cell)
	})

	return table
}

func (g *Gui) TableCell(title string, width int, color tcell.Color, selectable bool) *tview.TableCell {
	tcell := tview.NewTableCell(title)
	tcell.SetExpansion(width)
	tcell.SetAlign(tview.AlignCenter)
	tcell.SetTextColor(color)
	tcell.SetSelectable(selectable)

	return tcell
}


func (g *Gui) GetParams(table *tview.Table) Params {
	var params Params
	key := table.GetCell(1, 1).Text
	value := table.GetCell(1, 2).Text
	param := Param {
		Key: key,
		Value: value,
	}
	s := append(params, param)
	return s
}

func (g *Gui) GetParamsText(params Params) string {
	query := "?"
	for _, v := range params {
		if v.Key == "" || v.Value == "" {
			continue
		}
		text := fmt.Sprintf("%s=%s", v.Key, v.Value)
		query += text
	}
	return query
}