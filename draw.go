package main

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	la "github.com/laranatech/gorana/layout"
)

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
					spacer(1),
					keyNode('я'),
					keyNode('ч'),
					keyNode('с'),
					keyNode('м'),
					keyNode('и'),
					keyNode('т'),
					keyNode('ь'),
					keyNode('б'),
					keyNode('ю'),
					spacer(1),
				),
			),
			la.Node(
				la.Id("keyboard_row_3"),
				la.Row(),
				la.Width(la.Grow(1)),
				la.Gap(8),
				la.Children(
					spacer(1),
					keyNode('-'),
					keyNode('+'),
					spacer(1),
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
				la.Height(la.Fix(16)),
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
	g.DrawNode(screen, g.Node)
}

func (g *Game) DrawKey(screen *ebiten.Image, node *la.OutputItem) {
	c := color.RGBA{150, 150, 150, 255}

	tmp := strings.Replace(node.Id, "key_", "", 1)

	id := ' '

	for _, v := range tmp {
		id = v
	}

	if g.IsLetterGuessed(id) {
		if g.IsLetterInWord(id) {
			c = color.RGBA{255, 150, 0, 255}
		} else {
			c = color.RGBA{80, 80, 80, 255}
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

	DrawLetter(
		screen,
		id,
		node.X+(keySide/2)-((LetterWidth*s)/2),
		node.Y+(keySide/2)-((LetterWidth*s)/2),
		s,
		color.RGBA{255, 255, 255, 255},
	)

	if g.Hovered != nil && g.Hovered.Id == node.Id {
		vector.StrokeRect(
			screen,
			node.X,
			node.Y,
			node.W,
			node.H,
			2,
			color.RGBA{255, 255, 255, 255},
			false,
		)
	}
}

func (g *Game) DrawAttemptItem(screen *ebiten.Image, node *la.OutputItem) {
	c := color.RGBA{100, 100, 100, 255}

	vector.StrokeRect(
		screen,
		node.X,
		node.Y,
		node.W,
		node.H,
		2,
		c,
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

	c1 := getColorByStatus(status)

	vector.FillRect(screen, node.X, node.Y, node.W, node.H, c1, false)

	s := float32(4)

	DrawLetter(
		screen,
		w[i],
		node.X+(attemptItemSide/2)-((LetterWidth*s)/2),
		node.Y+(attemptItemSide/2)-((LetterWidth*s)/2),
		s,
		color.RGBA{255, 255, 255, 255},
	)
}

func getColorByStatus(status LetterStatus) color.Color {
	switch status {
	case GUESSED:
		return color.RGBA{0, 255, 0, 255}
	case PRESENT:
		return color.RGBA{255, 150, 0, 255}
	case WRONG:
		return color.RGBA{80, 80, 80, 255}
	}
	return color.RGBA{200, 200, 200, 255}
}

func (g *Game) DrawNode(screen *ebiten.Image, node *la.OutputItem) {
	if strings.HasPrefix(node.Id, "key_") {
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
