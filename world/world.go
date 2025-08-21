package world

import (
	"image/color"
	"slices"
	"cmp"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CellID int

const CELL_VOID CellID = 0

type Move struct {
	src int
	dst int
}

type Cell struct {
	Color color.RGBA
	Type  CellID
}

type Action func(cell *Cell, x, y int)

type Rule struct {
	action Action
}

type World struct {
	buffer      []color.RGBA
	cells       []Cell
	changes     []Move
	rules       []Rule
	rt          rl.RenderTexture2D

	width       int
	height      int
	scale       float32
	position    rl.Vector2
}

func New(width, height int, scale float32) *World {
	w := &World{
		width: width,
		height: height,
		scale:  scale,
		position: rl.Vector2{X: 0, Y: 0},
	}
	w.buffer = make([]color.RGBA, width * height)
	w.cells = make([]Cell, width * height)

	return w
}

func (w *World) Destroy() {
	rl.UnloadRenderTexture(w.rt)
}
func (w* World) AddRule(rule *Rule, id CellID) {

}

func (w *World) ToGlobalSpace(x, y int) (int, int) {
	return int(float32(x) * w.scale + w.position.X), int(float32(y) * w.scale + w.position.Y)
}

func (w *World) ToLocalSpace(x, y int) (int, int) {
	return int(float32(x) / w.scale - w.position.X), int(float32(y) / w.scale - w.position.Y)
}

func (w *World) InBounds(x, y int) bool {
	return !(x < 0 || x > w.width - 1 || y < 0 || y > w.height - 1)
}

func (w *World) To1D(x, y int) int {
	return y * w.width + x
}

func (w *World) To2D(index int) (int, int) {
	return index % w.width, index / w.width
}

func (w *World) Put(x, y int, cell Cell) {
	if !w.InBounds(x, y) {
		return
	}

	w.cells[w.To1D(x, y)] = cell
}

func (w *World) Get(x, y int) *Cell {
	if !w.InBounds(x, y) {
		return CELL_VOID
	}
	return &w.cells[w.To1D(x, y)]
}

func (w *World) IsVoid(x, y int) bool {
	return w.Get(x, y) == CELL_VOID
}

func (w *World) Move(srcX, srcY, dstX, dstY) {
	w.changes = append(w.changes, Move{src: w.To1D(srcX, srcY), dst: w.To1D(dstX, dstY)})
}

func (w *World) ApplyMoves() {
	w.changes = slices.DeleteFunc(w.changes, func(move Move){
		return w.cells[move.dst].Type != CELL_VOID
	})

	slices.SortFunc(w.changes, func(a, b Move) {
		return cmp.Compare(a.src, b.dst)
	})
	
	w.changes = append(w.changes, Move{src: -1, dst: -1})

	iprev := 0
	for i := 0; i < len(w.changes); i++ {
		if w.changes[i + 1].dst != w.changes[i].dst {
			rand := iprev + math.IntN(i - iprev)

			dst := w.changes[rand].dst
			src := w.changes[rand].src

			w.cells[dst] = w.cells[src]
			w.cells[src] = Cell{Type:CELL_VOID}

			iprev = i + 1
		}
	}

	w.changes = w.changes[:0]
}

func (w *World) Swap(x1, y1, x2, y2 int) {
	//w.grid.Swap(x1, y1, x2, y2)
}

func (w *World) Update() {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			cell := &w.cells[w.To1D(x, y)]
			if rule, ok := w.rules[cell.Type]; ok {
				rule.action(cell, x, y)
			}
		}
	}
}

func (w *World) Render() {
	rl.DrawTextureEx(w.rt.Texture, w.position, 0, w.scale, rl.White)
}
