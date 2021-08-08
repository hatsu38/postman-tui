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
	App   *tview.Application
	Pages *tview.Pages
}

func New() *Gui {
	g := &Gui{
		App:   tview.NewApplication(),
		Pages: tview.NewPages(),
	}
	return g
}

func (g *Gui) Run(i interface{}) error {
	app := g.App
	textView := g.TextView("Response")
	inputUrlField := g.Form("Request URL: ", "https://httpbin.org/get", "URL")
	inputParamsKeyField := g.Form("Params Key: ", "key", "Request Pararms Key")

	inputParamsKeyField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyTab:
			app.SetFocus(inputUrlField)
		}
	})

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
			app.SetFocus(inputParamsKeyField)
		}
	})

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow).
		AddItem(inputUrlField, 0, 1, true).
		AddItem(inputParamsKeyField, 0, 2, true).
		AddItem(textView, 0, 4, false)

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

func (g *Gui) Form(label string, placeholder string, title string) *tview.InputField {
	field := tview.NewInputField()
	field.SetLabel(label)
	field.SetPlaceholder(placeholder)
	field.SetTitle(title)
	field.SetFieldTextColor(tcell.ColorMaroon)
	field.SetLabelColor(tcell.ColorBlue)
	field.SetFieldBackgroundColor(tcell.ColorGray)
	field.SetPlaceholderTextColor(tcell.ColorWhite)
	field.SetBorder(true)

	return field
}

