package ui

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
			"url":         " Tab: move params table, Enter: http request",
			"http":        " Tab: move url field, Enter: change http method",
			"paramsTable": " Tab: move body table, Enter: set query paramater",
			"bodyTable":   " Tab: move params table, Enter: set body paramater",
		},
	}
}

func (n *navigate) update(panel string) {
	n.SetText(n.keybindings[panel])
}
