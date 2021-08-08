package main

import (
	"net/http"
	"fmt"
	"os"
	"io"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {

	app := tview.NewApplication()

	inputUrlField := tview.NewInputField()
	inputUrlField.SetLabel("Request URL: ")
	inputUrlField.SetPlaceholder("https://httpbin.org/get")
	inputUrlField.SetTitle("URL")
	inputUrlField.SetFieldTextColor(tcell.ColorMaroon)
	inputUrlField.SetLabelColor(tcell.ColorBlue)
	inputUrlField.SetFieldBackgroundColor(tcell.ColorGray)
	inputUrlField.SetPlaceholderTextColor(tcell.ColorWhite)
	inputUrlField.SetBorder(true)

	// Display Response Text View
	textView := tview.NewTextView()
	textView.SetTitle("Response")
	textView.SetBorder(true)

	// Params Key Field
	inputParamsKeyField := tview.NewInputField()
	inputParamsKeyField.SetLabel("Params Key: ")
	inputParamsKeyField.SetPlaceholder("key")
	inputParamsKeyField.SetTitle("Request Pararms Key")
	inputParamsKeyField.SetFieldTextColor(tcell.ColorMaroon)
	inputParamsKeyField.SetLabelColor(tcell.ColorBlue)
	inputParamsKeyField.SetFieldBackgroundColor(tcell.ColorGray)
	inputParamsKeyField.SetPlaceholderTextColor(tcell.ColorWhite)
	inputParamsKeyField.SetBorder(true)

	inputUrlField.SetChangedFunc(func(text string) {
		textView.SetText(text)
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

	inputParamsKeyField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyTab:
			app.SetFocus(inputUrlField)
		}
	})


	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow).
		AddItem(inputUrlField, 0, 1, true).
		AddItem(inputParamsKeyField, 0, 2, true).
		AddItem(textView, 0, 4, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
