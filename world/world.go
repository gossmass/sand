package world

import (
	"cmp"
	"image/color"
	"math/rand/v2"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CellID int

const CELL_VOID CellID = 0

type Move struct {
	src  int
	dst  int
	swap bool
}

type Cell struct {
	Color color.RGBA
	Type  CellID
}

func NewCell(color color.RGBA, id CellID) *Cell {
	return &Cell{
		Color: color,
		Type:  id,
	}
}

type Action func(w *World, cell *Cell, x, y int)

type Rule struct {
	action Action
}

func NewRule(action Action) *Rule {
	return &Rule{
		action: action,
	}
}

var EmptyAction Action = func(w *World, cell *Cell, x, y int) {}

type World struct {
	buffer  []color.RGBA
	cells   []*Cell
	changes []Move
	rules   map[CellID]*Rule
	rt      rl.RenderTexture2D

	width    int
	height   int
	scale    float32
	position rl.Vector2
}

func New(width, height int, scale float32) *World {
	w := &World{
		width:    width,
		height:   height,
		scale:    scale,
		position: rl.Vector2{X: 0, Y: 0},
	}
	w.buffer = make([]color.RGBA, width*height)
	w.cells = make([]*Cell, width*height)
	w.rules = make(map[CellID]*Rule)
	w.rt = rl.LoadRenderTexture(int32(width), int32(height))

	w.AddRule(NewRule(EmptyAction), CELL_VOID)

	return w
}

func (w *World) Destroy() {
	rl.UnloadRenderTexture(w.rt)
}
func (w *World) AddRule(rule *Rule, id CellID) {
	w.rules[id] = rule
}

func (w *World) ToGlobalSpace(x, y int) (int, int) {
	return int(float32(x)*w.scale + w.position.X), int(float32(y)*w.scale + w.position.Y)
}

func (w *World) ToLocalSpace(x, y int) (int, int) {
	return int(float32(x)/w.scale - w.position.X), int(float32(y)/w.scale - w.position.Y)
}

func (w *World) InBounds(x, y int) bool {
	return !(x < 0 || x > w.width-1 || y < 0 || y > w.height-1)
}

func (w *World) To1D(x, y int) int {
	return y*w.width + x
}

func (w *World) To2D(index int) (int, int) {
	return index % w.width, index / w.width
}

func (w *World) Put(x, y int, cell *Cell) {
	if !w.InBounds(x, y) {
		return
	}

	w.cells[w.To1D(x, y)] = cell
}

func (w *World) PutBlob(x, y, size int, cellToCopy Cell) {
	sx := max(0, x-size)
	ex := min(w.width, x+size)
	sy := max(0, y-size)
	ey := min(w.height, y+size)

	for ix := sx; ix < ex; ix++ {
		for iy := sy; iy < ey; iy++ {
			w.cells[w.To1D(ix, iy)] = NewCell(cellToCopy.Color, cellToCopy.Type)
		}
	}
}

func (w *World) Get(x, y int) *Cell {
	if !w.InBounds(x, y) {
		return nil
	}
	return w.cells[w.To1D(x, y)]
}

func (w *World) GetID(x, y int) CellID {
	if !w.InBounds(x, y) {
		return CELL_VOID
	}
	cell := w.cells[w.To1D(x, y)]

	if cell == nil {
		return CELL_VOID
	}
	return cell.Type

}

func (w *World) IsVoid(x, y int) bool {
	cell := w.Get(x, y)
	if cell == nil {
		return false
	}
	return cell.Type == CELL_VOID
}

func (w *World) move(srcX, srcY, dstX, dstY int, swap bool) {
	src := w.To1D(srcX, srcY)
	dst := w.To1D(dstX, dstY)

	if src < 0 || src > len(w.cells)-1 || dst < 0 || dst > len(w.cells)-1 {
		return
	}

	w.changes = append(w.changes, Move{src: src, dst: dst, swap: swap})
}

func (w *World) Move(srcX, srcY, dstX, dstY int) {
	w.move(srcX, srcY, dstX, dstY, false)
}

func (w *World) Swap(x1, y1, x2, y2 int) {
	w.move(x1, y1, x2, y2, true)
}

func (w *World) ApplyMoves() {
	w.changes = slices.DeleteFunc(w.changes, func(move Move) bool {
		dst := w.cells[move.dst]
		if dst == nil {
			return false
		}
		return dst.Type != CELL_VOID
	})

	slices.SortFunc(w.changes, func(a, b Move) int {
		return cmp.Compare(a.src, b.dst)
	})

	w.changes = append(w.changes, Move{src: -1, dst: -1})

	iprev := 0
	for i := 0; i < len(w.changes)-1; i++ {
		if w.changes[i+1].dst != w.changes[i].dst {
			rand := iprev + rand.IntN(i-iprev+1)
			cmd := w.changes[rand]

			if cmd.swap {
				w.cells[cmd.dst], w.cells[cmd.src] = w.cells[cmd.src], w.cells[cmd.dst]
			} else {
				w.cells[cmd.dst] = w.cells[cmd.src]
				w.cells[cmd.src] = nil //NewCell(rl.Blank, CELL_VOID)
			}

			iprev = i + 1
		}
	}

	w.changes = w.changes[:0]
}

func (w *World) Update() {
	for y := 0; y < w.height; y++ {
		for x := 0; x < w.width; x++ {
			cell := w.cells[w.To1D(x, y)]
			if cell == nil {
				continue
			}
			if rule, ok := w.rules[cell.Type]; ok {
				rule.action(w, cell, x, y)
			}
		}
	}

	w.ApplyMoves()
	w.UpdateTexture()
}

func (w *World) UpdateTexture() {
	for i := 0; i < len(w.cells); i++ {
		cell := w.cells[i]
		if cell == nil {
			continue
		}
		w.buffer[i] = w.cells[i].Color
	}
	rl.UpdateTexture(w.rt.Texture, w.buffer)
}
func (w *World) Render() {
	rl.DrawTextureEx(w.rt.Texture, w.position, 0, w.scale, rl.White)
}
