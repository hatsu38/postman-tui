package gui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type urlField struct {
	*tview.InputField
}

func newUrlField(label, text string) *urlField {
	field := &urlField{
		tview.NewInputField(),
	}
	field.SetLabel(label)
	field.SetFieldBackgroundColor(tcell.ColorBlack)
	field.SetBorder(true)
	field.SetBorderColor(tcell.ColorGreen)
	field.SetLabelColor(tcell.ColorIndianRed)
	field.SetText(text)
	return field
}

func (u *urlField) setFunc(g *Gui) {
	u.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			url := g.GetRequestUrl()
			resp := g.HttpRequest(url)
			defer resp.Body.Close()

			body := g.ParseResponse(resp)
			g.ResTextView.SetText(body)
		case tcell.KeyTab:
			g.ToFocus()
		}
	})
}
