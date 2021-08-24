package gui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type navigate struct {
	*tview.TextView
	keybindings map[string]string
}

func newNavigate() *navigate {
	return &navigate{
		TextView: tview.NewTextView().SetTextColor(tcell.ColorYellow),
		keybindings: map[string]string{
			"http":        " Tab: move url field\n Enter: change http method",
			"url":         " Tab: move Response Field\n Enter: http request",
			"resField":    " Tab: move params table\n Enter: Copy response text to clipboard",
			"paramsTable": " Tab: move body table\n Enter/Ctrl+C: set query paramater",
			"bodyTable":   " Tab: move params table\n Enter: set body paramater",
			"copied":      " Copied response text to clipboard!\n Tab: move params table  Enter: Copy response text to clipboard",
		},
	}
}

func (n *navigate) update(panel string) {
	n.SetText(n.keybindings[panel])
}
