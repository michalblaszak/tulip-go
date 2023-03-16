package main

import (
	"github.com/michalblaszak/tulip"
)

// "time"

// "github.com/gdamore/tcell/v2"

func main() {
	_statusBar := tulip.Desktop.AddStatusBar()

	// _label := tulip.CreateLabel(_statusBar)
	// _label.SetCoords(0, 0, 10, 1)
	// _label.SetColors(tcell.ColorWhite, ' ', tcell.ColorGreen)
	// _label.SetLabel("Ala ma kota")

	// _statusBar.CommandLabel = _label
	_statusBar.CreateCommandLabel()

	_menubar := tulip.Desktop.AddMenuBar()

	_fileMenu := _menubar.AddMenuItem("File")
	_editMenu := _menubar.AddMenuItem("Edit")

	/*_fileSubMenu :=*/
	_fileMenu.AddSubMenu()

	/*_editSubMenu :=*/
	_editMenu.AddSubMenu()

	tulip.App.Run()
}
