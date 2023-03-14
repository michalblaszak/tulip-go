package tulip

import "github.com/gdamore/tcell/v2"

const (
	MenuItemDefaultBkgColor          = tcell.ColorWhite
	MenuItemDefaultForeColor         = tcell.ColorBlack
	MenuItemDefaultBkgColorDisabled  = tcell.ColorWhite
	MenuItemDefaultForeColorDisabled = tcell.ColorGray
	MenuItemDefaultBkgColorSelected  = tcell.ColorBlack
	MenuItemDefaultForeColorSelected = tcell.ColorWhite
)

type MenuBar struct {
	region TRegion

	bkgColor   tcell.Color
	bkgPattern rune
	foreColor  tcell.Color

	menuItems []*MenuItem

	active bool
}

type MenuItem struct {
	label             string
	labelWidthCells   int
	enabled           bool
	selected          bool
	bkgColor          tcell.Color
	foreColor         tcell.Color
	bkgColorDisabled  tcell.Color
	foreColorDisabled tcell.Color
	bkgColorSelected  tcell.Color
	foreColorSelected tcell.Color

	Action func()
}

func (m *MenuBar) AddMenuItem(label string) *MenuItem {
	_menuItem := MenuItem{
		label:             label,
		enabled:           true,
		selected:          false,
		labelWidthCells:   StrCellWidth(label),
		bkgColor:          MenuItemDefaultBkgColor,
		foreColor:         MenuItemDefaultForeColor,
		bkgColorDisabled:  MenuItemDefaultBkgColorDisabled,
		foreColorDisabled: MenuItemDefaultForeColorDisabled,
		bkgColorSelected:  MenuItemDefaultBkgColorSelected,
		foreColorSelected: MenuItemDefaultForeColorSelected,
	}

	m.menuItems = append(m.menuItems, &_menuItem)

	//repaint = true;
	//Repaint()

	return &_menuItem
}

func (m *MenuBar) Paint(parentScreenX1, parentScreenY1, parentScreenX2, parentScreenY2 int) {
	// x1, _, x2, y2 := s.Widget.parent.getDeviceClientCoords(windowWithBorders)
	// _, parentClip := s.getDeviceClientCoords(windowWithBorders)

	screenX1, screenY1 := parentScreenX1+m.region.left, parentScreenY1+m.region.top
	screenX2 := screenX1 + m.region.w - 1
	screenY2 := screenY1 + m.region.h - 1

	if screenX2 > parentScreenX2 {
		screenX2 = parentScreenX2
	}

	if screenY2 > parentScreenY2 {
		screenY2 = parentScreenY2
	}

	st := tcell.StyleDefault
	st = st.Background(m.bkgColor)
	st = st.Foreground(m.foreColor)
	st = st.Bold(true)

	erase(screenX1, screenY1, screenX2, screenY2, st, m.bkgPattern)

	// Draw menu items
	st = tcell.StyleDefault
	x := screenX1
	for _, item := range m.menuItems {
		switch {
		case !item.enabled:
			st = st.Background(item.bkgColorDisabled)
			st = st.Foreground(item.foreColorDisabled)
		case m.active && item.enabled && item.selected:
			st = st.Background(item.bkgColorSelected)
			st = st.Foreground(item.foreColorSelected)
		case !m.active && item.enabled && item.selected:
			st = st.Background(item.bkgColor)
			st = st.Foreground(item.foreColor)
		case item.enabled && !item.selected:
			st = st.Background(item.bkgColor)
			st = st.Foreground(item.foreColor)
		}

		x += EmitStr(x, screenY1, screenX2, screenY2, st, " "+item.label+" ")
	}
}

func (m *MenuBar) resize(w, h int) {
	m.region.left, m.region.top = 0, 0
	m.region.w, m.region.h = w, 1

	App.repaint = true
}

func (m *MenuBar) ToggleActive() {
	if len(m.menuItems) == 0 {
		m.active = false
		// m.parent.setActiveWidget(nil)
		return
	}

	// The menubar contains some items
	if m.active {
		m.active = false
		// s.parent.setActiveWidget(nil)
	} else {
		m.active = true
		// s.parent.setActiveWidget(s)

		// found := false
		// for _, item := range m.menuItems {
		// 	if item.selected {
		// 		found = true
		// 		break
		// 	}
		// }

		// if !found {
		m.selectFirstAvailable()
		// }
	}

	//repaint = true;
	Repaint()
}

func (m *MenuBar) selectFirstAvailable() {
	for _, item := range m.menuItems {
		if item.enabled {
			item.selected = true
			break
		}
	}

	// repaint = true
	Repaint()
}

func (m *MenuBar) HandleEvent(ev IEvent) {
	switch ev.(type) {
	case *EventKey:
		event := ev.(*EventKey)
		switch event.Key {
		case tcell.KeyF10: // should be already catched ib Desktop
			m.ToggleActive()
			event.processed = true
		case tcell.KeyLeft:
			if m.active {
				m.activatePrevious()
				event.processed = true
			}
		case tcell.KeyRight:
			if m.active {
				m.activateNext()
				event.processed = true
			}
		case tcell.KeyEnter:
			if m.active {
				activeMenuItem := m.getActiveMenuItem()
				if activeMenuItem != nil && activeMenuItem.Action != nil {
					m.ToggleActive()
					activeMenuItem.Action()
				}
				event.processed = true
			}
		default:
			event.processed = true
		}

		// case EventTypeMouse:
		// case EventTypeConsole:
	}
}
