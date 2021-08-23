package ui

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Gui struct {
	App          *tview.Application
	Pages        *tview.Pages
	UrlField     *urlField
	ParamsTable  *table
	BodyTable    *table
	ResTextView  *resTestView
	HTTPTextView *httpTestView
	NavTextView  *navigate
}

type Param struct {
	Key   string
	Value string
}
type Params []Param

func New() *Gui {
	g := &Gui{
		App:          tview.NewApplication(),
		Pages:        tview.NewPages(),
		UrlField:     newUrlField(" Request URL: ", "https://httpbin.org/get"),
		ParamsTable:  newTable(),
		BodyTable:    newTable(),
		ResTextView:  newResTextView(),
		HTTPTextView: newHTTPTextView(),
		NavTextView:  newNavigate(),
	}
	return g
}

func (g *Gui) GetRequestUrl() string {
	field := g.UrlField
	urlText := field.GetText()
	params := g.ParamsTable.GetParams()
	query := g.GetParamsText(params)

	return urlText + query
}

func (g *Gui) HttpRequest(url string) *http.Response {

	bodyParams := g.BodyTable.GetParams()
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

	paramsTable.SetTableCells(g, "Query Params")
	bodyTable.SetTableCells(g, "Request Body")

	httpTextView.setFunc(g)
	inputUrlField.setFunc(g)
	paramsTable.setFunc(g)
	bodyTable.setFunc(g)

	httpFlex := tview.NewFlex()
	httpFlex.SetDirection(tview.FlexColumn)
	httpFlex.AddItem(g.HTTPTextView, 0, 1, false)
	httpFlex.AddItem(inputUrlField, 0, 9, true)

	requestFlex := tview.NewFlex()
	requestFlex.SetDirection(tview.FlexRow)
	requestFlex.AddItem(httpFlex, 0, 1, true)
	requestFlex.AddItem(paramsTable, 0, 5, false)
	requestFlex.AddItem(bodyTable, 0, 5, false)

	reqResflex := tview.NewFlex()
	reqResflex.SetDirection(tview.FlexColumn)
	reqResflex.AddItem(requestFlex, 0, 5, true)
	reqResflex.AddItem(resTextView, 0, 3, false)

	appFlex := tview.NewFlex()
	appFlex.SetDirection(tview.FlexRow)
	appFlex.AddItem(reqResflex, 0, 9, true)
	appFlex.AddItem(g.NavTextView, 1, 1, false)

	g.ToUrlFieldFocus()

	g.Pages.AddAndSwitchToPage("main", appFlex, true)

	if err := app.SetRoot(g.Pages, true).Run(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (g *Gui) ToUrlFieldFocus() {
	urlField := g.UrlField
	g.App.SetFocus(urlField)
	g.ParamsTable.SetBordersColor(tcell.ColorWhite)
	g.HTTPTextView.SetBorderColor(tcell.ColorWhite)
	g.BodyTable.SetBorderColor(tcell.ColorWhite)

	urlField.SetBorderColor(tcell.ColorGreen)
	g.NavTextView.update("url")
}

func (g *Gui) ToHTTPFieldFocus() {
	g.App.SetFocus(g.HTTPTextView)
	g.BodyTable.SetSelectable(false, false)
	g.ParamsTable.SetBordersColor(tcell.ColorWhite)
	g.BodyTable.SetBordersColor(tcell.ColorWhite)
	g.UrlField.SetBorderColor(tcell.ColorWhite)

	g.HTTPTextView.SetBorderColor(tcell.ColorGreen)
	g.NavTextView.update("http")
}

func (g *Gui) ToParamsTableFocus() {
	g.App.SetFocus(g.ParamsTable)
	g.ParamsTable.SetSelectable(true, true)
	g.HTTPTextView.SetBorderColor(tcell.ColorWhite)
	g.UrlField.SetBorderColor(tcell.ColorWhite)
	g.BodyTable.SetBorderColor(tcell.ColorWhite)

	g.ParamsTable.SetBordersColor(tcell.ColorGreen)
	g.NavTextView.update("paramsTable")

}

func (g *Gui) ToBodyTable() {
	g.App.SetFocus(g.BodyTable)
	g.ParamsTable.SetSelectable(false, false)
	g.BodyTable.SetSelectable(true, true)
	g.ParamsTable.SetBordersColor(tcell.ColorWhite)
	g.UrlField.SetBorderColor(tcell.ColorWhite)
	g.HTTPTextView.SetBorderColor(tcell.ColorWhite)

	g.BodyTable.SetBordersColor(tcell.ColorGreen)
	g.NavTextView.update("bodyTable")
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
	default:
		g.ToUrlFieldFocus()
	}
}

func (g *Gui) NewInputModal(table *table) {
	row, col := table.GetSelection()
	cell := table.GetCell(row, col)
	cell.SetTextColor(tcell.ColorWhite)

	text := cell.Text
	labelCell := table.GetCell(0, col)
	labelIndexCell := table.GetCell(row, 0)
	tableTitle := table.GetCell(0, 0)
	label := fmt.Sprintf(" %s %s %s: ", tableTitle.Text, labelCell.Text, labelIndexCell.Text)
	input := newUrlField(label, text)
	input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			txt := input.GetText()
			cell.Text = txt
			if txt != "" {
				table.AddParamsRow(row + 1)
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
