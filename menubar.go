package tulip

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

const (
	MenuItemDefaultBkgColor          = tcell.ColorWhite
	MenuItemDefaultForeColor         = tcell.ColorBlack
	MenuItemDefaultBkgColorDisabled  = tcell.ColorWhite
	MenuItemDefaultForeColorDisabled = tcell.ColorGray
	MenuItemDefaultBkgColorSelected  = tcell.ColorBlack
	MenuItemDefaultForeColorSelected = tcell.ColorWhite
)

const (
    submenuExpandIcon = BlackRightPointingTriangle
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
	region            TRegion // will be calculated dynamically
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

	submenu *SubMenu

	Action func()
}

type SubMenu struct {
	region TRegion // will be calculated dynamically

	bkgColor   tcell.Color
	bkgPattern rune
	foreColor  tcell.Color

	menuItems []*SubMenuItem

	active bool
}

type SubMenuItem struct {
	region            TRegion // will be calculated dynamically
	label             string // the label as declared
    labelDisplay      string // the label to be displayed
	labelWidthCells   int
	enabled           bool
	selected          bool
	bkgColor          tcell.Color
	foreColor         tcell.Color
	bkgColorDisabled  tcell.Color
	foreColorDisabled tcell.Color
	bkgColorSelected  tcell.Color
	foreColorSelected tcell.Color

	submenu *SubMenu

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

		strlen := EmitStr(x, screenY1, screenX2, screenY2, st, " "+item.label+" ")

		item.region = TRegion{
			left: x,
			top:  screenY1,
			w:    strlen,
			h:    1,
		}

		x += strlen

		if item.submenu != nil && item.submenu.active {
			item.submenu.Paint(parentScreenX1, parentScreenY1, parentScreenX2, parentScreenY2)
		}
	}

	// Process submenus
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
		m.reset()
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

func (m *MenuBar) reset() {
	for _, item := range m.menuItems {
		item.unselect()
	}

	// repaint = true
	//Repaint()
}

func (m *MenuBar) selectFirstAvailable() *MenuItem {
	var retItem *MenuItem = nil

	for _, item := range m.menuItems {
		if item.enabled {
			item.selected = true
			retItem = item
			break
		}
	}

	// repaint = true
	Repaint()
	return retItem
}

func (m *MenuBar) selectLastAvailable() *MenuItem {
	rangeLen := len(m.menuItems) - 1
	var retItem *MenuItem = nil
	for i := range m.menuItems {
		item := m.menuItems[rangeLen-i]
		if item.enabled {
			item.selected = true
			retItem = item
			break
		}
	}

	Repaint()
	return retItem
}

func (m *MenuBar) getActiveMenuItem() *MenuItem {
	if m == nil {
		return nil
	}

	var _ret *MenuItem = nil

	for _, item := range m.menuItems {
		if item.selected {
			_ret = item
			break
		}
	}

	return _ret
}

func (m *MenuBar) activateNext() *MenuItem {
	selectNext := false
	var retItem *MenuItem = nil
	for _, item := range m.menuItems {
		if item.selected {
			item.unselect()
			//item.selected = false
			// item.closeSubMenu(false)
			selectNext = true
			continue
		}
		if selectNext && item.enabled {
			item.selected = true
			selectNext = false
			retItem = item
			break
		}
	}

	// 'selectNext' remains selected if there were no more items available for selection after the previous one.
	// We need to select the first available one (roll over)
	if selectNext {
		retItem = m.selectFirstAvailable()
	}

	Repaint()
	return retItem
}

func (m *MenuBar) activatePrevious() *MenuItem {
	rangeLen := len(m.menuItems) - 1
	selectNext := false
	var retItem *MenuItem = nil

	for i := range m.menuItems {
		item := m.menuItems[rangeLen-i]
		if item.selected {
			item.unselect()
			// item.selected = false
			// item.closeSubMenu(false)
			selectNext = true
			continue
		}
		if selectNext && item.enabled {
			item.selected = true
			selectNext = false
			retItem = item
			break
		}
	}

	// 'selectNext' remains selected if there were no more items available for selection before the previous one.
	// We need to select the last available one (roll over)
	if selectNext {
		retItem = m.selectLastAvailable()
	}

	Repaint()
	return retItem
}

func (mi *MenuItem) unselect() {
	mi.selected = false
	mi.closeSubMenu(false)
	// TODO: close submenus if open
}

func (mi *MenuItem) openSubMenu() {
	mi.submenu.region = TRegion{
		left: mi.region.left,
		top:  mi.region.top + 1,
		// w:    will be udjusted in the adjustItems()
		// h:    will be udjusted in the adjustItems()
	}

    mi.submenu.adjustItems()

	mi.submenu.active = true
	mi.submenu.selectFirstAvailable()

	//Repaint() // Repaint is in selectFirstAvailable()
}

func (sm *SubMenu) adjustItems() {
	submenuHeight := 2 // 2 for the frame
	submenuWidth := 0
    hasSubmenu := false

    for _, i := range sm.menuItems {
		submenuHeight++
		submenuWidth = maxInt(submenuWidth, StrCellWidth(i.label))
        if i.submenu != nil {
            hasSubmenu = true
        }
	}

    if hasSubmenu {
        submenuWidth += 1+StrCellWidth(string(submenuExpandIcon)) // +1 for space
    }

    // Add trailing spaces and submenu expand icon to the items which need it
    for _, i := range sm.menuItems {
        if i.submenu == nil {
            i.labelDisplay = fmt.Sprintf("%-*s", submenuWidth, i.label)
        } else {
            numSpaces := submenuWidth - 1 // 1 for expand icon
            i.labelDisplay = fmt.Sprintf("%-*s%s", numSpaces, i.label, string(submenuExpandIcon))
        }
	}

    sm.region.w = submenuWidth + 2 // 2 for the frame
    sm.region.h = submenuHeight
}

func (mi *MenuItem) closeSubMenu(withRepaint bool) {
	mi.submenu.active = false

	if withRepaint {
		Repaint()
	}
}

func (m *MenuBar) HandleEvent(ev IEvent) {
	switch ev.(type) {
	case *EventKey:
		event := ev.(*EventKey)
		switch event.Key {
		// case tcell.KeyF10: // should be already catched ib Desktop
		// 	m.ToggleActive()    // TODO: deactivate all submenus
		// 	event.processed = true
		case tcell.KeyLeft:
			if m.active {
				activeMenuItem := m.getActiveMenuItem()

				if activeMenuItem != nil {
					hadSubmenuOpen := false
					if activeMenuItem.submenu != nil && activeMenuItem.submenu.active {
						hadSubmenuOpen = true
						activeMenuItem.submenu.HandleEvent(ev)
					}

					if !event.processed {
						activeMenuItem = m.activatePrevious()
						if hadSubmenuOpen {
							activeMenuItem.openSubMenu()
						}
						event.processed = true
					}
				} else {
					m.selectFirstAvailable()
					event.processed = true
				}
			}
		case tcell.KeyRight:
			if m.active {
				activeMenuItem := m.getActiveMenuItem()

				if activeMenuItem != nil {
					hadSubmenuOpen := false
					if activeMenuItem.submenu != nil && activeMenuItem.submenu.active {
						hadSubmenuOpen = true
						activeMenuItem.submenu.HandleEvent(ev)
					}

					if !event.processed {
						activeMenuItem = m.activateNext()
						if hadSubmenuOpen {
							activeMenuItem.openSubMenu()
						}
						event.processed = true
					}
				} else {
					m.selectFirstAvailable()
					event.processed = true
				}
			}
		case tcell.KeyEnter:
			if m.active {
				activeMenuItem := m.getActiveMenuItem()

				if activeMenuItem != nil {
					if activeMenuItem.submenu != nil {
						if activeMenuItem.submenu.active {
							activeMenuItem.submenu.HandleEvent((ev))
						}

						if !event.processed {
							activeMenuItem.openSubMenu()
							event.processed = true
						}
					} else if activeMenuItem.Action != nil {
						m.ToggleActive()
						activeMenuItem.Action()
						event.processed = true
					} else {
						event.processed = true
					}
				} else {
					event.processed = true
				}
			}
		default:
			event.processed = true
		}

		// case EventTypeMouse:
		// case EventTypeConsole:
	}
}

func (mi *MenuItem) AddSubMenu() *SubMenu {
	_subMenu := SubMenu{
		region: TRegion{ // will be calculated dynamically
			top:  0,
			left: 0,
			w:    -1,
			h:    -1,
		},
		bkgColor:   tcell.ColorWhite,
		bkgPattern: ' ',
		foreColor:  tcell.ColorBlack,
		active:     false,
	}

	mi.submenu = &_subMenu

	return &_subMenu
}

func (smi *SubMenuItem) AddSubMenu() *SubMenu {
	_subMenu := SubMenu{
		region: TRegion{ // will be calculated dynamically
			top:  0,
			left: 0,
			w:    -1,
			h:    -1,
		},
		bkgColor:   tcell.ColorWhite,
		bkgPattern: ' ',
		foreColor:  tcell.ColorBlack,
		active:     false,
	}

	smi.submenu = &_subMenu

	return &_subMenu
}

func (sm *SubMenu) Paint(parentScreenX1, parentScreenY1, parentScreenX2, parentScreenY2 int) {
	// x1, _, x2, y2 := s.Widget.parent.getDeviceClientCoords(windowWithBorders)
	// _, parentClip := s.getDeviceClientCoords(windowWithBorders)
	screenX1, screenY1 := parentScreenX1+sm.region.left, parentScreenY1+sm.region.top
	screenX2 := screenX1 + sm.region.w - 1
	screenY2 := screenY1 + sm.region.h - 1

	if screenX2 > parentScreenX2 {
		screenX2 = parentScreenX2
	}

	if screenY2 > parentScreenY2 {
		screenY2 = parentScreenY2
	}

	st := tcell.StyleDefault
	st = st.Background(sm.bkgColor)
	st = st.Foreground(sm.foreColor)
	st = st.Bold(true)

	erase(screenX1, screenY1, screenX2, screenY2, st, sm.bkgPattern)

	// Draw top border
	Screen.SetContent(screenX1, screenY1, BoxDrawingsLightDownAndRight, nil, st) // left-top corner
	Screen.SetContent(screenX2, screenY1, BoxDrawingsLightDownAndLeft, nil, st)  // right-top corner

	Screen.SetContent(screenX1, screenY2, BoxDrawingsLightUpAndRight, nil, st) // left-bottom corner
	Screen.SetContent(screenX2, screenY2, BoxDrawingsLightUpAndLeft, nil, st)  // right-bottom corner

	for x := 1; x < sm.region.w-1; x++ { // 1 and -1 as corners are drawn above
		Screen.SetContent(screenX1+x, screenY1, BoxDrawingsLightHorizontal, nil, st)
		Screen.SetContent(screenX1+x, screenY2, BoxDrawingsLightHorizontal, nil, st)
	}
	// Draw menu items
	st_item := tcell.StyleDefault
	y := screenY1
	for _, item := range sm.menuItems {
		switch {
		case !item.enabled:
			st_item = st_item.Background(item.bkgColorDisabled)
			st_item = st_item.Foreground(item.foreColorDisabled)
		case sm.active && item.enabled && item.selected:
			st_item = st_item.Background(item.bkgColorSelected)
			st_item = st_item.Foreground(item.foreColorSelected)
		case !sm.active && item.enabled && item.selected:
			st_item = st_item.Background(item.bkgColor)
			st_item = st_item.Foreground(item.foreColor)
		case item.enabled && !item.selected:
			st_item = st_item.Background(item.bkgColor)
			st_item = st_item.Foreground(item.foreColor)
		}
		Screen.SetContent(screenX1, y+1, BoxDrawingsLightVertical, nil, st) // left-border
		Screen.SetContent(screenX2, y+1, BoxDrawingsLightVertical, nil, st) // right-border

		EmitStr(screenX1+1, y+1, screenX2-1, y+1, st_item, item.labelDisplay) // +1 for the border
		y++
	}
}

func (m *SubMenu) AddMenuItem(label string) *SubMenuItem {
	_menuItem := SubMenuItem{
		label:             label,
        labelDisplay:      "",
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

// func (smi *SubMenuItem) getPrintableLabel() string {
// 	commandIcon := ""
// 	if smi.submenu != nil {
// 		commandIcon = " " + string(BlackRightPointingTriangle /*BlackMediumRightPointingTriangleCentered*/)
// 	}
// 	return smi.label + commandIcon
// }

// TODO: Make it generic (same as for MenuBar)
func (m *SubMenu) selectFirstAvailable() {
	for _, item := range m.menuItems {
		if item.enabled {
			item.selected = true
			break
		}
	}

	// repaint = true
	Repaint()
}

// TODO: Make it generic (same as for MenuBar)
func (m *SubMenu) selectLastAvailable() {
	rangeLen := len(m.menuItems) - 1
	for i := range m.menuItems {
		item := m.menuItems[rangeLen-i]
		if item.enabled {
			item.selected = true
			break
		}
	}

	Repaint()
}

// TODO: Make it generic (same as for MenuBar)
func (m *SubMenu) getActiveMenuItem() *SubMenuItem {
	if m == nil {
		return nil
	}

	var _ret *SubMenuItem = nil

	for _, item := range m.menuItems {
		if item.selected {
			_ret = item
			break
		}
	}

	return _ret
}

func (m *SubMenu) activateNext() {
	selectNext := false
	for _, item := range m.menuItems {
		if item.selected {
			item.selected = false
			selectNext = true
			continue
		}
		if selectNext && item.enabled {
			item.selected = true
			selectNext = false
			break
		}
	}

	// 'selectNext' remains selected if there were no more items available for selection after the previous one.
	// We need to select the first available one (roll over)
	if selectNext {
		m.selectFirstAvailable()
	}

	Repaint()
}

func (m *SubMenu) activatePrevious() {
	rangeLen := len(m.menuItems) - 1
	selectNext := false
	for i := range m.menuItems {
		item := m.menuItems[rangeLen-i]
		if item.selected {
			item.selected = false
			selectNext = true
			continue
		}
		if selectNext && item.enabled {
			item.selected = true
			selectNext = false
		}
	}

	// 'selectNext' remains selected if there were no more items available for selection before the previous one.
	// We need to select the last available one (roll over)
	if selectNext {
		m.selectLastAvailable()
	}

	Repaint()
}

func (sm *SubMenu) HandleEvent(ev IEvent) {
	switch ev.(type) {
	case *EventKey:
		event := ev.(*EventKey)
		switch event.Key {
		// case tcell.KeyF10: // should be already catched ib Desktop
		// 	m.ToggleActive()
		// 	event.processed = true
		case tcell.KeyRight:
			if sm.active {
				activeMenuItem := sm.getActiveMenuItem()
				if activeMenuItem != nil {
					if activeMenuItem.submenu != nil {
						activeMenuItem.openSubMenu()
						event.processed = true
					}
				}
			}
		// case tcell.KeyRight:
		// 	if m.active {
		// 		m.activateNext()
		// 		event.processed = true
		// 	}
		case tcell.KeyEnter:
			if sm.active {
				activeMenuItem := sm.getActiveMenuItem()
				if activeMenuItem != nil {
					if activeMenuItem.submenu != nil {
						activeMenuItem.openSubMenu()
						event.processed = true
					} else if activeMenuItem.Action != nil {
						DeactivateMenu()
						activeMenuItem.Action()
						event.processed = true
					}
				}
			}
		default:
			// event.processed = true
		}

		// case EventTypeMouse:
		// case EventTypeConsole:
	}
}

func (smi *SubMenuItem) openSubMenu() {
	smi.submenu.region = TRegion{
		left: smi.region.left + smi.region.w,
		top:  smi.region.top + 1,
		// w:    will be udjusted in the adjustItems()
		// h:    will be udjusted in the adjustItems()
	}

    smi.submenu.adjustItems()

	smi.submenu.active = true
	smi.submenu.selectFirstAvailable()

	//Repaint() // Repaint is in selectFirstAvailable()
}

func (smi *SubMenuItem) Paint() {

}
