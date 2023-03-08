package tulip

import (
	"fmt"
	"os"

	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/encoding"
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
}

func InitScreen() {
	encoding.Register()

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
	// w, h := Screen.Size()

	Desktop.Paint()

	// midocui.EmitStr(w/2-7, h/2, tcell.StyleDefault, "Hello, Michałl!")
	// midocui.EmitStr(w/2-9, h/2+1, tcell.StyleDefault, "Press ESC to exit.")

	Screen.Show()
}

func typedCommand(ev *tcell.EventKey) ttypedCommand {
	if ev.Modifiers()&tcell.ModAlt == tcell.ModAlt {
		App.appMode = appModeTypedCommand

		if ev.Key() == tcell.KeyF10 {
			App.appMode = appModeNormal
			App.commandBuffer = ""
			return typedCommandAppMenu
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
		case "x":
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
			case typedCommandNone:
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
			case typedCommandUnknown: // This command is unknown to Desktop. Delegate it to the active window for potential processing
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
			case typedCommandMove, typedCommandResize, typedCommandAppMenu:
				event = &EventTypedCommand{
					Event: &Event{
						timestamp: time.Now(),
						processed: false,
					},
					Command: typedCommand,
				}
				Desktop.HandleEvent(event)
			case typedCommandAppQuit:
				Screen.Fini()
				os.Exit(0)
			}

			// if a.appMode == appModeTypedCommand && event.Processed() {
			//     a.appMode = appModeNormal
			//     a.commandBuffer = ""
			// }
		case *SysEventQuit:
			Screen.Fini()
			os.Exit(0)
		case *AppEventCloseCurrentWin:
			Desktop.CloseCurrentWin()
		case *AppEventRepaint:
			App.repaint = true
		}

		if App.repaint {
			App.repaint = false
			Paint()
		}
	}
}

// Hello returns a greeting for the named person.
func Hello(name string) string {
	// Return a greeting that embeds the name in a message.
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
	return message
}
