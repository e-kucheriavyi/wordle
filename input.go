package main

import (
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	la "github.com/laranatech/gorana/layout"
)

const (
	clickInputDebounce = 250
)

var keymap map[string]rune = map[string]rune{
	"q":         'й',
	"w":         'ц',
	"e":         'у',
	"r":         'к',
	"t":         'е',
	"y":         'н',
	"u":         'г',
	"i":         'ш',
	"o":         'щ',
	"p":         'з',
	"[":         'х',
	"]":         'ъ',
	"a":         'ф',
	"s":         'ы',
	"d":         'в',
	"f":         'а',
	"g":         'п',
	"h":         'р',
	"j":         'о',
	"k":         'л',
	"l":         'д',
	";":         'ж',
	"'":         'э',
	"z":         'я',
	"x":         'ч',
	"c":         'с',
	"v":         'м',
	"b":         'и',
	"n":         'т',
	"m":         'ь',
	",":         'б',
	".":         'ю',
	"Enter":     '+',
	"Backspace": '-',
}

func (g *Game) InputKey() rune {
	keys := inpututil.AppendJustReleasedKeys(nil)

	for _, key := range keys {
		return MapInputToRune(key)
	}

	return ' '
}

func MapInputToRune(k ebiten.Key) rune {
	name := ""

	switch k {
	case ebiten.KeyEnter:
		name = "Enter"
	case ebiten.KeyBackspace:
		name = "Backspace"
	case ebiten.KeyQ:
		name = "q"
	case ebiten.KeyW:
		name = "w"
	case ebiten.KeyE:
		name = "e"
	case ebiten.KeyR:
		name = "r"
	case ebiten.KeyT:
		name = "t"
	case ebiten.KeyY:
		name = "y"
	case ebiten.KeyU:
		name = "u"
	case ebiten.KeyI:
		name = "i"
	case ebiten.KeyO:
		name = "o"
	case ebiten.KeyP:
		name = "p"
	case ebiten.KeyBracketLeft:
		name = "["
	case ebiten.KeyBracketRight:
		name = "]"
	case ebiten.KeyA:
		name = "a"
	case ebiten.KeyS:
		name = "s"
	case ebiten.KeyD:
		name = "d"
	case ebiten.KeyF:
		name = "f"
	case ebiten.KeyG:
		name = "g"
	case ebiten.KeyH:
		name = "h"
	case ebiten.KeyJ:
		name = "j"
	case ebiten.KeyK:
		name = "k"
	case ebiten.KeyL:
		name = "l"
	case ebiten.KeySemicolon:
		name = ";"
	case ebiten.KeyQuote:
		name = "'"
	case ebiten.KeyZ:
		name = "z"
	case ebiten.KeyX:
		name = "x"
	case ebiten.KeyC:
		name = "c"
	case ebiten.KeyV:
		name = "v"
	case ebiten.KeyB:
		name = "b"
	case ebiten.KeyN:
		name = "n"
	case ebiten.KeyM:
		name = "m"
	case ebiten.KeyComma:
		name = ","
	case ebiten.KeyPeriod:
		name = "."
	}

	r, ok := keymap[name]

	if !ok {
		return ' '
	}
	return r
}

func (g *Game) IsOkToClick() bool {
	return time.Since(g.LastClickedAt) > clickInputDebounce*time.Millisecond
}

func FindHovered(node *la.OutputItem, x, y float32) *la.OutputItem {
	if strings.HasPrefix(node.Id, "key_") {
		if Collide(node, x, y) {
			return node
		}
	}

	for _, child := range node.Children {
		hovered := FindHovered(child, x, y)
		if hovered != nil {
			return hovered
		}
	}

	return nil
}

func (g *Game) IsPressed() bool {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		return true
	}

	touches := ebiten.AppendTouchIDs(nil)

	return len(touches) > 0
}

func (g *Game) CursorPosition() (float32, float32) {
	x, y := ebiten.CursorPosition()
	if x != 0 && y != 0 {
		return float32(x), float32(y)
	}

	touches := ebiten.AppendTouchIDs(nil)

	if len(touches) == 0 {
		return 0, 0
	}

	x, y = ebiten.TouchPosition(touches[0])

	return float32(x), float32(y)
}
