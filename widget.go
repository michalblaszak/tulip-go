package tulip

type IWidget interface {
	Paint(parentScreenX1, parentScreenY1, parentScreenX2, parentScreenY2 int)
}