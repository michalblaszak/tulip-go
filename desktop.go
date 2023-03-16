package tulip

import (
	"github.com/gdamore/tcell/v2"
)

type TDesktop struct {
	region     TRegion
	bkgColor   tcell.Color
	bkgPattern rune
	foreColor  tcell.Color
	statusbar  *StatusBar
	menubar    *MenuBar
}

// func (d *TDesktop) erase(screenX1, screenY1, screenX2, screenY2 int, st tcell.Style, bk rune) {
// 	for x := screenX1; x <= screenX2; x++ {
// 		for y := screenY1; y <= screenY2; y++ {
// 			Screen.SetContent(x, y, bk, nil, st)
// 		}
// 	}
// }

func (d *TDesktop) Paint(parentScreenX1, parentScreenY1, parentScreenX2, parentScreenY2 int) {
	screenX1, screenY1 := parentScreenX1+d.region.left, parentScreenY1+d.region.top
	screenX2 := screenX1 + d.region.w - 1
	screenY2 := screenY1 + d.region.h - 1

	if screenX2 > parentScreenX2 {
		screenX2 = parentScreenX2
	}

	if screenY2 > parentScreenY2 {
		screenY2 = parentScreenY2
	}

	st := tcell.StyleDefault
	st = st.Bold(false)
	st = st.Background(d.bkgColor)
	st = st.Foreground(d.foreColor)

	erase(screenX1, screenY1, screenX2, screenY2, st, d.bkgPattern)

	if d.menubar != nil {
		d.menubar.Paint(screenX1, screenY1, screenX2, screenY2)
	}

	if d.statusbar != nil {
		d.statusbar.Paint(screenX1, screenY1, screenX2, screenY2)
	}

	// d.Window.Paint()
	// for _, win := range d.childWindows {
	//     win.Paint()
	// }
}

func (d *TDesktop) resize(w, h int) {
	d.region.w, d.region.h = w, h

	if d.menubar != nil {
		d.menubar.resize(w, h)
	}

	if d.statusbar != nil {
		d.statusbar.resize(w, h)
	}

	App.repaint = true
}

// HandleEvent processes system level events and delegates them to the appropriate widget for further processing.
// Keys reserved for the desktop only:
// Alt-F10             : activate and deactivate the desktop's menubar
// Alt-x           : exit the application
// Ctrl-Tab        : move to the next window
// Other keystrokes: will be delegate to the currently active window or decktop's menubar (if active)
func (d *TDesktop) HandleEvent(ev IEvent) {
	switch ev.(type) {
	case *EventTypedCommand:
		event := ev.(*EventTypedCommand)
		switch event.Command {
		case typedCommandAppMenu:
			if d.menubar != nil {
				d.menubar.ToggleActive()
				event.processed = true
			}

		case typedCommandUnknown:
			if d.menubar != nil && d.menubar.active {
				d.menubar.HandleEvent(ev)
			}
		}
	case *EventKey:
		if d.menubar != nil && d.menubar.active {
			d.menubar.HandleEvent(ev)
		}
	}
}

func (d *TDesktop) AddStatusBar() *StatusBar {
	_statusBar := StatusBar{
		region: TRegion{
			top:  0, // actually it will be always calculated for the status bar
			left: 0,
			w:    -1, // actually it will be always calculated for the status bar
			h:    1,
		},
		bkgColor:   tcell.ColorWhite,
		bkgPattern: ' ',
		foreColor:  tcell.ColorBlack,
	}

	d.statusbar = &_statusBar

	return &_statusBar
}

func (d *TDesktop) SetCommandLabelText(s string) {
	if d.statusbar != nil {
		d.statusbar.SetCommandLabelText(s)
	}
}

func (d *TDesktop) AddMenuBar() *MenuBar {
	_menuBar := MenuBar{
		region: TRegion{
			top:  0,
			left: 0,
			w:    -1, // actually it will be always calculated for the status bar
			h:    1,
		},
		bkgColor:   tcell.ColorWhite,
		bkgPattern: ' ',
		foreColor:  tcell.ColorBlack,
		active:     false,
	}

	d.menubar = &_menuBar

	return &_menuBar
}
