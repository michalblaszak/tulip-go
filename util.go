package tulip

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

func inClipRecion(x, y, x1, y1, x2, y2 int) bool {
	return (x >= x1) && (x <= x2) &&
		(y >= y1) && (y <= y2)
}

// Returns the number of screen cells ocupied by the string
func EmitStr(x1, y1, x2, y2 int, style tcell.Style, str string) int {
	cellX := x1

	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}

		if w == 1 {
			if inClipRecion(cellX, y1, x1, y1, x2, y2) {
				Screen.SetContent(cellX, y1, c, comb, style)
			}
			cellX += w
		} else { // w == 2
			inClip1 := inClipRecion(cellX, y1, x1, y1, x2, y2)   // Is 1st half of the rune visible?
			inClip2 := inClipRecion(cellX+1, y1, x1, y1, x2, y2) // Is 2nd half of the rune visible?

			if inClip1 && inClip2 { // Both rune parts are visible
				Screen.SetContent(cellX, y1, c, comb, style)
				cellX += w
			} else if inClip1 && !inClip2 { // 1st part is visible, 2nd is not
				// Just print ' ' instead of the entire double-width rune
				comb = []rune{c}
				c = ' '
				w = 1
				Screen.SetContent(cellX, y1, c, comb, style)
				cellX += w
			} else if !inClip1 && inClip2 { // 1st part isn't visible, 2nd is visible
				// Just print ' ' instead of the entire double-width rune
				comb = []rune{c}
				c = ' '
				w = 1
				Screen.SetContent(cellX+1, y1, c, comb, style)
				cellX += w
			} else { // Neither rune part is visible
				cellX += w
			}
		}
	}

	return cellX - x1
}

func StrCellWidth(s string) int {
	len := 0

	for _, c := range s {
		w := runewidth.RuneWidth(c)
		if w == 0 {
			w = 1
		}
		len += w
	}
	return len
}

func erase(screenX1, screenY1, screenX2, screenY2 int, st tcell.Style, bk rune) {
	for x := screenX1; x <= screenX2; x++ {
		for y := screenY1; y <= screenY2; y++ {
			Screen.SetContent(x, y, bk, nil, st)
		}
	}
}

// func minInt(a, b int) int {
//     if (a < b) {
//         return a
//     }

//     return b
// }

func maxInt(a, b int) int {
    if (a > b) {
        return a
    }

    return b
}

// // Checks if i is in [left; right]
// func in (i, left, right int) bool {
//     return i >= left && i <= right
// }

// // reverseInt - reverses the slice of ints
// func reverseInt(s []int) {
// 	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
//     		s[i], s[j] = s[j], s[i]
// 	}
// }

// // splitToGigits - splits the int number into a slice of its digits
// func splitToDigits(n int) []int{
// 	var _ret []int

// 	for n !=0 {
// 		_ret = append(_ret, n % 10)
// 		n /= 10
// 	}

// 	reverseInt(_ret)

// 	return _ret
// }
