package gui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type httpTestView struct {
	*tview.TextView
}

func newHTTPTextView() *httpTestView {
	textView := &httpTestView{
		tview.NewTextView(),
	}
	textView.SetTitle(" HTTP ")
	textView.SetBorder(true)
	textView.SetScrollable(true)
	textView.SetText("GET")
	textView.SetTextColor(tcell.ColorGreen)
	textView.SetToggleHighlights(true)
	textView.SetTextAlign(tview.AlignCenter)

	return textView
}

func (t *httpTestView) setFunc(g *Gui) {
	t.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			g.NewHTTPListModal()
		case tcell.KeyTab:
			g.ToFocus()
		}
	})
}
