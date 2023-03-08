package tulip

type TDesktop struct {
	region TRegion
}

func (d *TDesktop) Paint() {
	// d.Window.Paint()
	// for _, win := range d.childWindows {
	//     win.Paint()
	// }
}

func (d *TDesktop) resize(w, h int) {
	d.region.w, d.region.h = w, h
	repaint = true
}
