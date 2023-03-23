package tulip

import (
	"fmt"
	"os"
	"time"

	"strconv"
	//	"time"

	"github.com/gdamore/tcell/v2"
)

type ttypedCommand int

const (
	typedCommandNone    ttypedCommand = iota // Not a command at all (Alt not pressed)
	typedCommandUnknown                      // In a command typing mode (Alt pressed) but command not recognized
	typedCommandResize
	typedCommandMove
	typedCommandAppQuit
	typedCommandAppMenu
)

type tAppMode int

const (
	appModeNormal tAppMode = iota
	appModeTypedCommand
)

type SApp struct {
	appMode       tAppMode
	commandBuffer string
	repaint       bool
}

var App SApp = SApp{
	appMode: appModeNormal,
	repaint: false,
}

var Screen tcell.Screen
var Desktop = TDesktop{
	region: TRegion{
		top:  0,
		left: 0,
		w:    10,
		h:    10,
	},
	bkgColor:   tcell.ColorBlue,
	bkgPattern: tcell.RuneBoard,
	foreColor:  tcell.ColorWhite,
}

func InitScreen() {
	// encoding.Register()

	// Console initialization
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	Screen = s

	if e := Screen.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	Screen.SetStyle(defStyle)

	Screen.Clear()

	// Initialize desktop window
	Desktop.region.w, Desktop.region.h = Screen.Size()
}

func Paint() {
	w, h := Screen.Size()

	Desktop.Paint(0, 0, w-1, h-1) // We pass a screen coordinates of Desktop

	// midocui.EmitStr(w/2-7, h/2, tcell.StyleDefault, "Hello, Michałl!")
	// midocui.EmitStr(w/2-9, h/2+1, tcell.StyleDefault, "Press ESC to exit.")

	Screen.Show()
}

func typedCommand(ev *tcell.EventKey) ttypedCommand {
	if ev.Modifiers()&tcell.ModAlt == tcell.ModAlt {
		App.appMode = appModeTypedCommand

        switch ev.Key() {
		case tcell.KeyF10:
			App.appMode = appModeNormal
			App.commandBuffer = ""
			return typedCommandAppMenu
        case tcell.KeyEnter:
			App.appMode = appModeNormal
			App.commandBuffer = ""
			return typedCommandUnknown
        }

		App.commandBuffer += string(ev.Rune())

		switch App.commandBuffer {
		case "resize":
			App.appMode = appModeNormal
			App.commandBuffer = ""
			return typedCommandResize
		case "move":
			App.appMode = appModeNormal
			App.commandBuffer = ""
			return typedCommandMove
		case "menu":
			App.appMode = appModeNormal
			App.commandBuffer = ""
			return typedCommandAppMenu
		case "X":
			App.appMode = appModeNormal
			App.commandBuffer = ""
			return typedCommandAppQuit
		default:
			return typedCommandUnknown
		}
	} else {
		App.appMode = appModeNormal
		App.commandBuffer = ""
		return typedCommandNone
	}
}

func (a *SApp) Run() {
	InitScreen()
	Paint()

	// TODO: Remove: debug
	f, err := os.Create("dat2")
	if err != nil {
		panic(err)
	}
	// TODO: End remove

	for {
		// On this level we handle App level events
		switch ev := Screen.PollEvent().(type) {
		case *tcell.EventResize:
			// TODO: Remove: debug
			xx, yy := ev.Size()
			_, err := f.WriteString(strconv.Itoa(xx) + "," + strconv.Itoa(yy) + "\n")
			if err != nil {
				panic(err)
			}
			// TODO: End remove
			Screen.Sync()
			Desktop.resize(Screen.Size())
		case *tcell.EventKey:
			typedCommand := typedCommand(ev)
			var event IEvent

			switch typedCommand {
			case typedCommandNone: // The user just types without Alt key
				event = &EventKey{
					Event: &Event{
						timestamp: time.Now(),
						processed: false,
					},
					Key:       ev.Key(),
					Rune:      ev.Rune(),
					Modifiers: ev.Modifiers(),
				}
				Desktop.HandleEvent(event)

			case typedCommandUnknown: // This command is unknown to App. Delegate it to Desktop for potential processing
				event = &EventTypedCommand{
					Event: &Event{
						timestamp: time.Now(),
						processed: false,
					},
					Command: typedCommandUnknown,
				}
				Desktop.HandleEvent(event)

				if event.Processed() {
					a.appMode = appModeNormal
					a.commandBuffer = ""
				}

			case typedCommandAppMenu:
				event = &EventTypedCommand{
					Event: &Event{
						timestamp: time.Now(),
						processed: false,
					},
					Command: typedCommand,
				}
				Desktop.HandleEvent(event)

			// case typedCommandMove, typedCommandResize, typedCommandAppMenu:
			// 	event = &EventTypedCommand{
			// 		Event: &Event{
			// 			timestamp: time.Now(),
			// 			processed: false,
			// 		},
			// 		Command: typedCommand,
			// 	}
			// 	Desktop.HandleEvent(event)
			case typedCommandAppQuit:
				Screen.Fini()
				os.Exit(0)
			} // ~TypedCommand

			Desktop.SetCommandLabelText(a.commandBuffer)
			App.repaint = true
			// if a.appMode == appModeTypedCommand && event.Processed() {
			//     a.appMode = appModeNormal
			//     a.commandBuffer = ""
			// }
		case *SysEventQuit:
			Screen.Fini()
			os.Exit(0)

		// case *AppEventCloseCurrentWin:
		//     Desktop.CloseCurrentWin()
		case *AppEventRepaint:
			App.repaint = true

        case *AppEventDeactivateMenu:
            if Desktop.menubar != nil {
                Desktop.menubar.active = false
                App.repaint = true
            }
		}


		if App.repaint {
			App.repaint = false
			Paint()
		}
	}
}

func Repaint() {
	appEv := &AppEventRepaint{}
	appEv.SetEventNow()
	if Screen == nil {
		fmt.Println("Repaint: Screen==nil")
	}
	_ = Screen.PostEvent(appEv)
}

func DeactivateMenu() {
    appEv := &AppEventDeactivateMenu{}
	appEv.SetEventNow()
	if Screen == nil {
		fmt.Println("DeactivateMenu: Screen==nil")
	}
	_ = Screen.PostEvent(appEv)
}