package world

import (
	"image/color"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Chunk struct {
	buffer []rl.Color
}

type TileID int

type Action func(w *World, x, y int)

type Rule struct {
	action Action
	color  color.RGBA
}

func NewRule(action Action, color color.RGBA) *Rule {
	return &Rule{
		action: action,
		color:  color,
	}
}

func EmptyRule(w *World, x, y int) {}

type World struct {
	chunks      []Chunk
	buffer      []color.RGBA
	rules       map[TileID]*Rule
	rulesColors map[color.RGBA]TileID
	width       int
	height      int
	rt          rl.RenderTexture2D
	scale       float32
}

func New(width, height int, scale float32) *World {
	w := &World{
		width:  width,
		height: height,
		scale:  scale,
	}
	w.rules = make(map[TileID]*Rule)
	w.rulesColors = make(map[color.RGBA]TileID)
	w.buffer = make([]color.RGBA, width*height)
	w.rt = rl.LoadRenderTexture(int32(width), int32(height))

	return w
}

func (w *World) Destroy() {
	rl.UnloadRenderTexture(w.rt)
}

func (w *World) AddRule(rule *Rule, tileID TileID) {
	w.rules[tileID] = rule
	w.rulesColors[rule.color] = tileID
}

func (w *World) to1D(x, y int) int {
	return y*w.width + x
}

func (w *World) to1DCheck(x, y int) (int, bool) {
	index := w.to1D(x, y)
	return index, !(index < 0 || index > len(w.buffer)-1)
}

func (w *World) rawPut(index int, tileID TileID) {
	if rule, ok := w.rules[tileID]; ok {
		w.buffer[index] = rule.color
	}

}

func (w *World) ToGlobalSpace(x, y int) (int, int) {
	return int(float32(x) * w.scale), int(float32(y) * w.scale)
}

func (w *World) ToLocalSpace(x, y int) (int, int) {
	return int(float32(x) / w.scale), int(float32(y) / w.scale)
}

func (w *World) Put(x, y int, tileID TileID) {
	if index, ok := w.to1DCheck(x, y); ok {
		w.rawPut(index, tileID)
	}
}

func (w *World) PutBlob(x, y, size int, tileID TileID) {
	sx := max(0, x-size)
	ex := min(w.width-1, x+size)
	sy := max(0, y-size)
	ey := min(w.height-1, y+size)

	for x := sx; x < ex; x++ {
		for y := sy; y < ey; y++ {
			index := w.to1D(x, y)
			w.rawPut(index, tileID)
		}
	}
}

func (w *World) Get(x, y int) TileID {
	if index, ok := w.to1DCheck(x, y); ok {
		if ruleID, ok := w.rulesColors[w.buffer[index]]; ok {
			return ruleID
		}
	}
	return 0
}

func (w *World) Swap(x1, y1, x2, y2 int) {
	index1, ok := w.to1DCheck(x1, y1)
	if !ok {
		return
	}
	index2, ok := w.to1DCheck(x2, y2)
	if !ok {
		return
	}
	w.buffer[index1], w.buffer[index2] = w.buffer[index2], w.buffer[index1]
}

func (w *World) Update() {
	for x := 0; x < w.width; x++ {
		for y := w.height - 1; y >= 0; y-- {
			index := w.to1D(x, y)
			col := w.buffer[index]

			if ruleID, ok := w.rulesColors[col]; ok {
				w.rules[ruleID].action(w, x, y)
			} else {
				log.Printf("Failed to get rule for color: %v\n", col)
			}
		}
	}

	rl.UpdateTexture(w.rt.Texture, w.buffer)
}

func (w *World) Render() {
	rl.DrawTextureEx(w.rt.Texture, rl.Vector2{}, 0, w.scale, rl.White)
}

// func (w *World) Render() {
// }
