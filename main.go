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

	inputField := tview.NewInputField()
	inputField.SetLabel("Request URL: ")
	inputField.SetPlaceholder("https://httpbin.org/get")
	inputField.SetTitle("URL")
	inputField.SetFieldTextColor(tcell.ColorMaroon)
	inputField.SetLabelColor(tcell.ColorBlue)
	inputField.SetFieldBackgroundColor(tcell.ColorGray)
	inputField.SetPlaceholderTextColor(tcell.ColorWhite)
	inputField.SetBorder(true)

	textView := tview.NewTextView()
	textView.SetTitle("Response")
	textView.SetBorder(true)

	inputField.SetChangedFunc(func(text string) {
		textView.SetText(text)
	})
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
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
		}
	})


	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow).
		AddItem(inputField, 3, 0, true).
		AddItem(textView, 0, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
