package main

import (
	"fmt"
	"image/color"
	"strings"
	"unicode/utf8"

	"github.com/e-kucheriavyi/wordle/pallete"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	la "github.com/laranatech/gorana/layout"
)

var backspaceMap = &[]byte{
	0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 1, 1, 1, 1, 1, 1,
	0, 0, 1, 0, 0, 0, 0, 0, 1,
	0, 1, 0, 0, 1, 0, 1, 0, 1,
	1, 0, 0, 0, 0, 1, 0, 0, 1,
	0, 1, 0, 0, 1, 0, 1, 0, 1,
	0, 0, 1, 0, 0, 0, 0, 0, 1,
	0, 0, 0, 1, 1, 1, 1, 1, 1,
	0, 0, 0, 0, 0, 0, 0, 0, 0,
}

func attemptItem(r, i int) *la.NodeItem {
	return la.Node(
		la.Id(fmt.Sprintf("attempt_%d_%d", r, i)),
		la.Width(la.Fix(attemptItemSide)),
		la.Height(la.Fix(attemptItemSide)),
	)
}

func attemptRow(r int) *la.NodeItem {
	return la.Node(
		la.Id(fmt.Sprintf("attempt-row_%d", r)),
		la.Row(),
		la.Gap(8),
		la.Children(
			attemptItem(r, 0),
			attemptItem(r, 1),
			attemptItem(r, 2),
			attemptItem(r, 3),
			attemptItem(r, 4),
		),
	)
}

func keyNode(key rune) *la.NodeItem {
	return la.Node(
		la.Id(fmt.Sprintf("key_%c", key)),
		la.Width(la.Fix(keySide)),
		la.Height(la.Fix(keySide)),
	)
}

func growKeyNode(key rune) *la.NodeItem {
	return la.Node(
		la.Id(fmt.Sprintf("key_%c", key)),
		la.Width(la.Grow(1)),
		la.Height(la.Fix(keySide)),
	)
}

func keyboardNode() *la.NodeItem {
	return la.Node(
		la.Id("keyboard"),
		la.Column(),
		la.Width(la.Grow(1)),
		la.Gap(keyRowGap),
		la.Children(
			la.Node(
				la.Id("keyboard_row_0"),
				la.Row(),
				la.Width(la.Grow(1)),
				la.Gap(keyGap),
				la.Children(
					spacer(1),
					keyNode('й'),
					keyNode('ц'),
					keyNode('у'),
					keyNode('к'),
					keyNode('е'),
					keyNode('н'),
					keyNode('г'),
					keyNode('ш'),
					keyNode('щ'),
					keyNode('з'),
					keyNode('х'),
					keyNode('ъ'),
					spacer(1),
				),
			),
			la.Node(
				la.Id("keyboard_row_1"),
				la.Row(),
				la.Width(la.Grow(1)),
				la.Gap(8),
				la.Children(
					spacer(1),
					keyNode('ф'),
					keyNode('ы'),
					keyNode('в'),
					keyNode('а'),
					keyNode('п'),
					keyNode('р'),
					keyNode('о'),
					keyNode('л'),
					keyNode('д'),
					keyNode('ж'),
					keyNode('э'),
					spacer(1),
				),
			),
			la.Node(
				la.Id("keyboard_row_2"),
				la.Row(),
				la.Width(la.Grow(1)),
				la.Gap(8),
				la.Children(
					growKeyNode('+'),
					keyNode('я'),
					keyNode('ч'),
					keyNode('с'),
					keyNode('м'),
					keyNode('и'),
					keyNode('т'),
					keyNode('ь'),
					keyNode('б'),
					keyNode('ю'),
					growKeyNode('-'),
				),
			),
		),
	)
}

func spacer(v float32) *la.NodeItem {
	return la.Node(
		la.Id("spacer"),
		la.Width(la.Grow(v)),
	)
}

func CreateLayout() *la.OutputItem {
	root := la.Node(
		la.Id("root"),
		la.Gap(8),
		la.Padding(8),
		la.Width(la.Fix(screenW)),
		la.Height(la.Fix(screenH)),
		la.Column(),
		la.Children(
			la.Node(
				la.Id("header"),
				la.Height(la.Fix(64)),
				la.Width(la.Grow(1)),
			),
			la.Node(
				la.Id("top"),
				la.Row(),
				la.Width(la.Grow(1)),
				la.Height(la.Fit()),
				la.Children(
					la.Node(
						la.Id("top-spacer-left"),
						la.Width(la.Grow(1)),
					),
					la.Node(
						la.Id("attempts"),
						la.Column(),
						la.Width(la.Fit()),
						la.Height(la.Fit()),
						la.Gap(8),
						la.Children(
							attemptRow(0),
							attemptRow(1),
							attemptRow(2),
							attemptRow(3),
							attemptRow(4),
							attemptRow(5),
						),
					),
					la.Node(
						la.Id("top-spacer-right"),
						la.Width(la.Grow(1)),
					),
				),
			),
			la.Node(
				la.Id("bottom"),
				la.Width(la.Grow(1)),
				la.Height(la.Grow(1)),
				la.Children(
					keyboardNode(),
				),
			),
		),
	)

	la.Layout(root)

	node := la.Export(root)

	return node
}

func (g *Game) Draw(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, screenW, screenH, pallete.BG, false)

	switch g.Stage {
	case GAME:
		g.DrawNode(screen, g.Node)
	case SCORE:
		g.DrawScore(screen)
	}
}

func (g *Game) DrawScore(screen *ebiten.Image) {
	txt := string(g.Word)
	s := float32(5)
	if g.IsWordGuessed() {
		s = 8
		txt = fmt.Sprintf("%d / %d", len(g.GuessedWords), 6)
	}

	DrawText(
		screen,
		txt,
		screenW/2-(float32(utf8.RuneCountInString(txt)-2)*(LetterWidth*s)),
		screenH/2-LetterWidth*s,
		s,
		pallete.FG,
	)
}

func (g *Game) DrawKey(screen *ebiten.Image, node *la.OutputItem) {
	c := pallete.PASSIVE

	tmp := strings.Replace(node.Id, "key_", "", 1)

	id := ' '

	for _, v := range tmp {
		id = v
	}

	if g.IsLetterGuessed(id) {
		if g.IsLetterInWord(id) {
			c = pallete.PRESENT
		} else {
			c = pallete.MISS
		}
	}

	vector.FillRect(
		screen,
		node.X,
		node.Y,
		node.W,
		node.H,
		c,
		false,
	)

	s := float32(4)

	x := node.X+(node.W/2)-((LetterWidth*s)/2)
	y := node.Y+(node.H/2)-((LetterWidth*s)/2)

	if id == '-' {
		DrawBitmap(
			screen,
			backspaceMap,
			x,
			y,
			s,
			9,
			pallete.FG,
		)
	} else if id == '+' {
		txt := "enter"
		s := float32(1.5)
		x := node.X + (node.W / 2) - (float32(len(txt)) * LetterWidth * s) / 2
		y := node.Y + (node.H / 2) - (LetterWidth * 2) / 2
		DrawText(
			screen,
			"enter",
			x,
			y,
			s,
			pallete.FG,
		)
	} else {
		DrawLetter(
			screen,
			id,
			x,
			y,
			s,
			pallete.FG,
		)
	}

	if g.Hovered != nil && g.Hovered.Id == node.Id {
		vector.StrokeRect(
			screen,
			node.X,
			node.Y,
			node.W,
			node.H,
			2,
			pallete.FG,
			false,
		)
	}
}

func (g *Game) DrawAttemptItem(screen *ebiten.Image, node *la.OutputItem) {
	vector.StrokeRect(
		screen,
		node.X,
		node.Y,
		node.W,
		node.H,
		2,
		pallete.FG,
		false,
	)

	r, i := ExtractIndecies(node.Id)

	if r > len(g.GuessedWords)-1 {
		return
	}

	w := g.GuessedWords[r]

	if i > len(w)-1 {
		return
	}

	status := g.GetLetterStatus(r, i, w[i])

	c := getColorByStatus(status)

	vector.FillRect(screen, node.X, node.Y, node.W, node.H, c, false)

	s := float32(4)

	DrawLetter(
		screen,
		w[i],
		node.X+(attemptItemSide/2)-((LetterWidth*s)/2),
		node.Y+(attemptItemSide/2)-((LetterWidth*s)/2),
		s,
		pallete.FG,
	)
}

func getColorByStatus(status LetterStatus) color.Color {
	switch status {
	case GUESSED:
		return pallete.MATCH
	case PRESENT:
		return pallete.PRESENT
	case WRONG:
		return pallete.MISS
	}
	return pallete.PASSIVE
}

func (g *Game) DrawHeader(screen *ebiten.Image, node *la.OutputItem) {
	v := len(g.GuessedWords) - 1
	if v < 0 {
		v = 0
	}
	s := float32(4)

	txt := fmt.Sprintf("%d / %d", v, 6)
	DrawText(
		screen,
		txt,
		node.X+node.W/2-(float32(len(txt))*LetterWidth*s)/2,
		node.Y+node.Y/2+(LetterWidth*s)/2,
		s,
		pallete.FG,
	)
}

func (g *Game) DrawNode(screen *ebiten.Image, node *la.OutputItem) {
	if node.Id == "header" {
		g.DrawHeader(screen, node)
	} else if strings.HasPrefix(node.Id, "key_") {
		g.DrawKey(screen, node)
	} else if strings.HasPrefix(node.Id, "attempt_") {
		g.DrawAttemptItem(screen, node)
	}

	for _, child := range node.Children {
		g.DrawNode(screen, child)
	}
}

func (g *Game) Layout(a, b int) (int, int) {
	return screenW, screenH
}
