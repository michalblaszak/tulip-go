package tulip

import "github.com/gdamore/tcell/v2"

type StatusBar struct {
	// NonBorderedWidget
    region TRegion
	bkgColor   tcell.Color
	bkgPattern rune
    foreColor  tcell.Color
    
    // TODO: it's temporary
    CommandLabel *Label
}

func (s *StatusBar) erase(screenX1, screenY1, screenX2, screenY2 int, st tcell.Style, bk rune) {
    for x := screenX1; x <= screenX2; x++ {
        for y := screenY1; y <= screenY2; y++ {
            Screen.SetContent(x, y, bk, nil, st)
        }
    }
}

func (s *StatusBar) Paint(parentScreenX1, parentScreenY1, parentScreenX2, parentScreenY2 int) {
	// x1, _, x2, y2 := s.Widget.parent.getDeviceClientCoords(windowWithBorders)
    // _, parentClip := s.getDeviceClientCoords(windowWithBorders)

    screenX1, screenY1 := parentScreenX1 + s.region.left, parentScreenY1 + s.region.top
    screenX2 := screenX1 + s.region.w-1
    screenY2 := screenY1 + s.region.h-1

    if screenX2 > parentScreenX2 {
        screenX2 = parentScreenX2
    }

    if screenY2 > parentScreenY2 {
        screenY2 = parentScreenY2
    }

    st := tcell.StyleDefault
    st = st.Background(s.bkgColor)
    st = st.Foreground(s.foreColor)
    st = st.Bold(true)

    s.erase(screenX1, screenY1, screenX2, screenY2, st, s.bkgPattern)

    if s.CommandLabel != nil {
        s.CommandLabel.Paint(screenX1, screenY1, screenX2, screenY2)
    }
}

func (s *StatusBar) resize(w, h int) {
    s.region.left, s.region.top = 0, h-1
	s.region.w, s.region.h = w, 1

	App.repaint = true
}

func (s *StatusBar) SetCommandLabelText(str *string) {
    if s.CommandLabel != nil {
        s.CommandLabel.SetLabel("cmd:" + *str)
    }
}

func (s *StatusBar) CreateCommandLabel() {
    _label := CreateLabel(s)
    _label.SetCoords(0, 0, 10, 1)
    _label.SetColors(tcell.ColorWhite, ' ', tcell.ColorGreen)
    s.SetCommandLabelText(&"")
    s.CommandLabel = _label

}

// func (s *StatusBar) HandleEvent(ev IEvent) {
//     switch ev.(type) {
//     case *EventKey:
//         if App.appMode == appModeTypedCommand {
//             s.CommandLabel.label = App.commandBuffer
//             Repaint()
//         }
//     case *EventTypedCommand:
//         s.CommandLabel.label = App.commandBuffer
//         Repaint()
//     }

//     return
// }

// func (s *StatusBar) getDeviceClientCoords(_ TClientAreaType) (region Region, clipRegion ClippingRegion) {
//     var parentRegion Region
//     var parentClip ClippingRegion

//     screen_w, screen_h := Screen.Size()
//     if s.parent == nil {
//         region = Region{
//             x1: 0,
//             y1: screen_h - 1,
//             x2: screen_w - 1,
//             y2: screen_h - 1,
//         }
//         parentClip = ClippingRegion{
//             x1: region.x1,
//             y1: region.y1,
//             x2: region.x2,
//             y2: region.y2,
//         }
//     } else {
//         parentRegion, parentClip = s.parent.getDeviceClientCoords(windowWithBorders)

//         region = Region{
//             x1: parentRegion.x1,
//             y1: parentRegion.y2,
//             x2: parentRegion.x2,
//             y2: parentRegion.y2,
//         }
//     } 

//     // Adjust clipRegion
//     if s.parent != nil {
//         clipRegion.x1 = maxInt(region.x1, parentClip.x1)
//         clipRegion.y1 = maxInt(region.y1, parentClip.y1)
//         clipRegion.x2 = minInt(region.x2, parentClip.x2)
//         clipRegion.y2 = minInt(region.y2, parentClip.y2)
//     }

//     return
// }
