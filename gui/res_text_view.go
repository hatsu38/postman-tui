package gui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	defaultText = `
                _                                _         _
_ __   ___  ___| |_ _ __ ___   __ _ _ __        | |_ _   _(_)
| '_ \ / _ \/ __| __| '_ ' _ \ / _' | '_ \ _____| __| | | | |
| |_) | (_) \__ \ |_| | | | | | (_| | | | |_____| |_| |_| | |
| .__/ \___/|___/\__|_| |_| |_|\__,_|_| |_|      \__|\__,_|_|
|_|`
)

type resTestView struct {
	*tview.TextView
}

func newResTextView() *resTestView {
	textView := &resTestView{
		tview.NewTextView(),
	}
	textView.SetTitle(" Response ")
	textView.SetBorder(true)
	textView.SetScrollable(true)
	textView.SetText(defaultText)
	textView.SetTextColor(tcell.ColorGreen)
	textView.SetToggleHighlights(true)

	return textView
}

func (t *resTestView) setFunc(g *Gui) {
	t.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyTab:
			g.ToFocus()
		}
	})
}
