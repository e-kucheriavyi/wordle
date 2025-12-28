package main

import (
	_ "embed"
	"image/color"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	la "github.com/laranatech/gorana/layout"
)

const (
	screenW         = 720
	screenH         = 720
	keyGap          = 8
	keyRowGap       = 8
	keySide         = (screenW - (keyGap * 11)) / 12
	attemptItemSide = 58
	inputDebounce   = 250
)

type Stage byte

const (
	INTRO Stage = iota
	GAME
	SCORE
)

type LetterStatus byte

const (
	WRONG LetterStatus = iota
	PRESENT
	GUESSED
	PENDING
)

var bg = color.RGBA{100, 100, 100, 100}

type Game struct {
	Stage         Stage
	Word          []rune
	GuessedWords  [][]rune
	Node          *la.OutputItem
	Hovered       *la.OutputItem
	LastClickedAt time.Time
	LastSubmitted int
}

func NewGame() *Game {
	return &Game{
		Stage:         GAME,
		Word:          []rune(GetWord(time.Now())),
		GuessedWords:  make([][]rune, 0, 6),
		Node:          CreateLayout(),
		LastSubmitted: -1,
	}
}

func (g *Game) Update() error {
	if g.Stage == GAME {
		g.UpdateGame()
	}

	return nil
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

func Collide(node *la.OutputItem, x, y float32) bool {
	if node.X > x || node.X+node.W < x {
		return false
	}

	if node.Y > y || node.Y+node.H < y {
		return false
	}

	return true
}

func (g *Game) IsPressed() bool {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		return true
	}

	touches := ebiten.TouchIDs()

	return len(touches) > 0
}

func (g *Game) CursorPosition() (float32, float32) {
	x, y := ebiten.CursorPosition()
	if x != 0 && y != 0 {
		return float32(x), float32(y)
	}

	touches := ebiten.TouchIDs()

	if len(touches) == 0 {
		return 0, 0
	}

	x, y = ebiten.TouchPosition(touches[0])

	return float32(x), float32(y)
}

func (g *Game) UpdateGame() error {
	x, y := g.CursorPosition()

	g.Hovered = FindHovered(g.Node, float32(x), float32(y))

	isPressed := g.IsPressed()

	if !isPressed {
		return nil
	}

	if time.Since(g.LastClickedAt) < inputDebounce*time.Millisecond {
		return nil
	}

	g.LastClickedAt = time.Now()

	if g.Hovered == nil {
		return nil
	}

	tmp := strings.ReplaceAll(g.Hovered.Id, "key_", "")

	l := []rune(tmp)[0]

	if l == '-' {
		return g.HandleBackspace()
	}

	if l == '+' {
		return g.HandleSubmit()
	}

	return g.HandleLetterClick(l)
}

func (g *Game) HandleLetterClick(l rune) error {
	if len(g.GuessedWords) == 0 {
		g.GuessedWords = append(g.GuessedWords, []rune{l})
		return nil
	}

	lastIndex := len(g.GuessedWords) - 1

	w := g.GuessedWords[lastIndex]

	if len(w) < 5 {
		g.GuessedWords[lastIndex] = append(g.GuessedWords[lastIndex], l)
	}

	return nil
}

func (g *Game) HandleBackspace() error {
	if len(g.GuessedWords) == 0 {
		return nil
	}

	lastIndex := len(g.GuessedWords) - 1

	if lastIndex == g.LastSubmitted {
		return nil
	}

	if len(g.GuessedWords[lastIndex]) > 0 {
		g.GuessedWords[lastIndex] = g.GuessedWords[lastIndex][:len(g.GuessedWords[lastIndex])-1]
	}

	return nil
}

func (g *Game) HandleSubmit() error {
	lastIndex := len(g.GuessedWords) - 1

	if g.LastSubmitted == lastIndex {
		return nil
	}

	if len(g.GuessedWords[lastIndex]) != 5 {
		return nil
	}

	if !ValidateWord(string(g.GuessedWords[lastIndex])) {
		return nil
	}

	if g.IsWordGuessed() || len(g.GuessedWords) == 6 {
		g.Stage = SCORE
		return nil
	}

	g.LastSubmitted = lastIndex
	g.GuessedWords = append(g.GuessedWords, make([]rune, 0, 5))

	return nil
}

func (g *Game) IsWordGuessed() bool {
	lastIndex := len(g.GuessedWords) - 1

	w := g.GuessedWords[lastIndex]

	for i := range 6 {
		if i > 4 {
			return true
		}
		if w[i] != g.Word[i] {
			return false
		}
	}

	return false
}

func (g *Game) IsLetterGuessed(letter rune) bool {
	if g.LastSubmitted == -1 {
		return false
	}

	for i, w := range g.GuessedWords {
		if i > g.LastSubmitted {
			break
		}
		for _, c := range w {
			if c == letter {
				return true
			}
		}
	}
	return false
}

func (g *Game) IsLetterInWord(letter rune) bool {
	if g.LastSubmitted == -1 {
		return false
	}

	for _, c := range g.Word {
		if c == letter {
			return true
		}
	}
	return false
}

func (g *Game) GetLetterStatus(r, i int, letter rune) LetterStatus {
	if g.LastSubmitted < r {
		return PENDING
	}
	if g.Word[i] == letter {
		return GUESSED
	}
	if g.IsLetterInWord(letter) {
		return PRESENT
	}
	return WRONG
}

func ExtractIndecies(str string) (int, int) {
	tmp := strings.ReplaceAll(str, "attempt_", "")

	vals := strings.Split(tmp, "_")

	v0, _ := strconv.Atoi(vals[0])
	v1, _ := strconv.Atoi(vals[1])

	return v0, v1
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("Wordle")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err.Error())
	}
}
