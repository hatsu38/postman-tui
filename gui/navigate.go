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
			"paramsTable": " Tab: move body table,   Enter: set query paramater,   h/left arrow: Move left by one column,   l/right arrow: Move right by one column,\n j/down arrow: Move down by one row,   k/up arrow: Move up by one row,   g/home: Move to the top,   G/end: Move to the bottom",
			"bodyTable":   " Tab: move params table,   Enter: set body paramater,   h/left arrow: Move left by one column,   l/right arrow: Move right by one column,\n j/down arrow: Move down by one row,   k/up arrow: Move up by one row,   g/home: Move to the top,   G/end: Move to the bottom",
			"copied":      " Copied response text to clipboard!\n Tab: move params table  Enter: Copy response text to clipboard",
		},
	}
}

func (n *navigate) update(panel string) {
	n.SetText(n.keybindings[panel])
}
