package tulip

type TRegion struct {
	top, left, w, h int
}

type ClippingRegion struct {
    x1, y1, x2, y2 int
}

func (clip *ClippingRegion) inClipRecion(x, y int) bool {
    return (x >= clip.x1) && (x <= clip.x2) && 
           (y >= clip.y1) && (y <= clip.y2)
}