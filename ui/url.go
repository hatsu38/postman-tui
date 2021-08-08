package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)
type Url struct {
	*tview.InputField
}

func (s *tview.InputField) Done() string {
	// if key == tcell.KeyEscape {
	// 	s.SetText("")
	// 	return 
	// }

	// if key != tcell.KeyEnter {
	// 	return
	// }
	println(s)
	text := s.GetText()
	return text
}

func NewRequestUrl() *Url {
	s := &Url{
		InputField: tview.NewInputField(),
	}
	// s := tview.NewInputField().
	s.SetLabel("Enter a URL: ")
	s.SetPlaceholder("https://httpbin.org/get")
	s.SetFieldWidth(10)
	s.SetLabelColor(tcell.Color252)
	s.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			s.SetText("")
			return
		}
		if key != tcell.KeyEnter {
			return
		}
		println(s)
		text := s.GetText()
		println(text)
	})

	return s
}