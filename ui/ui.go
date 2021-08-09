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
	TextView *tview.TextView
	TableView *tview.Table
	App   *tview.Application
	Pages *tview.Pages
}

func New() *Gui {
	g := &Gui{
		TextView: TextView("Response"),
		TableView: nil,
		App:   tview.NewApplication(),
		Pages: tview.NewPages(),
	}
	return g
}

func (g *Gui) Run(i interface{}) error {
	app := g.App
	textView := g.TextView
	inputUrlField := g.Form("Request URL: ", "https://httpbin.org/get", "URL")
	tableView := g.Table()

	inputUrlField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			text := textView.GetText(true)
			req, _ := http.NewRequest("GET", text, nil)
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
			toFixBody := "{" + string(body) + "}}"
			textView.SetText(toFixBody)
		case tcell.KeyTab:
			app.SetFocus(tableView)
			tableView.SetSelectable(true, true)
			tableView.SetBordersColor(tcell.ColorPaleVioletRed)
			tableView.Select(1, 1)
			inputUrlField.SetFieldBackgroundColor(tcell.ColorGray)
		}
	})

	tableView.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			row, column := tableView.GetSelection()
			currentCell := tableView.GetCell(row, column)
			currentCell.SetTextColor(tcell.ColorBlue)
			currentCell.SetTransparency(true)
			tableView.SetSelectable(true, true)
		case tcell.KeyTab:
			app.SetFocus(inputUrlField)
			inputUrlField.SetFieldBackgroundColor(tcell.ColorPaleVioletRed)
			tableView.SetBordersColor(tcell.ColorWhite)
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


func TextView(title string) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetTitle(title)
	textView.SetBorder(true)

	return textView
}

func (g *Gui) Input(tableView *tview.Table, cell *tview.TableCell, label string, width int) {
	input := tview.NewInputField()
	text := cell.Text
	input.SetText(text)
	input.SetLabel(label)
	input.SetLabelWidth(width)
	input.Autocomplete()
	input.SetFieldBackgroundColor(tcell.ColorPaleVioletRed)
	input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			cell.Text = input.GetText()
			tableView.SetSelectable(false, false)
			g.Pages.RemovePage("input")
			g.App.SetFocus(tableView)
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

func (g *Gui) Form(label string, placeholder string, title string) *tview.InputField {
	field := tview.NewInputField()
	field.SetLabel(label)
	field.SetPlaceholder(placeholder)
	field.SetTitle(title)
	field.SetFieldTextColor(tcell.ColorMaroon)
	field.SetLabelColor(tcell.ColorBlue)
	field.SetFieldBackgroundColor(tcell.ColorPaleVioletRed)
	field.SetPlaceholderTextColor(tcell.ColorWhite)
	field.SetBorder(true)

	return field
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
		cell.SetTextColor(tcell.ColorBlue)

		g.Input(table, cell, "params", 10)
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