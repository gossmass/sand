package main

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"sand/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	CELL_SAND world.CellID = iota + 1
	CELL_WATER
	CELL_OIL
	CELL_SOLID
	CELL_FIRE
	CELL_ASH

	CELL_MAX
)

func randSideDir() int {
	dir := rand.IntN(2)
	if dir == 1 {
		return 1
	}

	return -1

}

func sandRule(w *world.World, cell *world.Cell, x, y int) {
	if w.IsVoid(x, y+1) {
		w.Move(x, y, x, y+1)
	} else if w.IsVoid(x-1, y+1) {
		w.Move(x, y, x-1, y+1)
	} else if w.IsVoid(x+1, y+1) {
		w.Move(x, y, x+1, y+1)
	}
}

func waterRule(w *world.World, cell *world.Cell, x, y int) {
	down := w.GetID(x, y+1)

	if down == world.CELL_VOID {
		w.Move(x, y, x, y+1)
	} else if down == CELL_SAND || down == CELL_OIL {
		w.Swap(x, y, x, y+1)
	} else {
		dir := randSideDir()
		if w.IsVoid(x+dir, y) {
			w.Move(x, y, x+dir, y)
		}
	}
}

func oilRule(w *world.World, cell *world.Cell, x, y int) {

}

func fireRule(w *world.World, cell *world.Cell, x, y int) {

}

func ashRule(w *world.World, cell *world.Cell, x, y int) {

}

func getCellName(cellID world.CellID) string {
	switch cellID {
	case world.CELL_VOID:
		return "VOID"
	case CELL_SAND:
		return "SAND"
	case CELL_WATER:
		return "WATER"
	case CELL_OIL:
		return "OIL"
	case CELL_SOLID:
		return "SOLID"
	case CELL_FIRE:
		return "FIRE"
	case CELL_ASH:
		return "ASH"
	}

	return "Undefined"
}

func getCellColor(cellID world.CellID) color.RGBA {
	switch cellID {
	case CELL_SAND:
		return rl.Green
	case CELL_WATER:
		return rl.Blue
	case CELL_SOLID:
		return rl.Gray
	case CELL_OIL:
		return rl.NewColor(24, 24, 24, 255)
	case CELL_FIRE:
		return rl.NewColor(200, 10, 10, 255)
	case CELL_ASH:
		return rl.NewColor(190, 190, 190, 255)
	}
	return rl.Blank

}

func main() {
	var boardWidth int32 = 128
	var boardHeight int32 = 96
	var boardScale float32 = 8.0
	var cursorSize int = 3
	var cursorCellID world.CellID = CELL_SAND
	var isPaused bool = false

	rl.InitWindow(int32(float32(boardWidth)*boardScale), int32(float32(boardHeight)*boardScale), "Title")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	w := world.New(int(boardWidth), int(boardHeight), boardScale)
	defer w.Destroy()

	w.AddRule(world.NewRule(sandRule), CELL_SAND)
	w.AddRule(world.NewRule(waterRule), CELL_WATER)
	w.AddRule(world.NewRule(world.EmptyAction), CELL_SOLID)
	w.AddRule(world.NewRule(oilRule), CELL_OIL)
	w.AddRule(world.NewRule(fireRule), CELL_FIRE)
	w.AddRule(world.NewRule(ashRule), CELL_ASH)

	bgImg := rl.GenImageChecked(int(float32(boardWidth)*boardScale), int(float32(boardHeight)*boardScale), int(boardScale)*4, int(boardScale)*4, rl.DarkGray, rl.Black)
	bgTex := rl.LoadTextureFromImage(bgImg)
	rl.UnloadImage(bgImg)
	defer rl.UnloadTexture(bgTex)

	for !rl.WindowShouldClose() {

		mx := int(rl.GetMouseX())
		my := int(rl.GetMouseY())
		lx, ly := w.ToLocalSpace(mx, my)

		if rl.IsMouseButtonDown(rl.MouseButtonRight) {
			cell := world.NewCell(getCellColor(cursorCellID), cursorCellID)
			w.Put(lx, ly, cell)

			if isPaused {
				w.UpdateTexture()
			}
		}

		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			w.PutBlob(lx, ly, cursorSize, world.Cell{Color: getCellColor(cursorCellID), Type: cursorCellID})

			if isPaused {
				w.UpdateTexture()
			}
		}

		if w := rl.GetMouseWheelMove(); w != 0 {
			dir := 1
			if w < 0 {
				dir = -1
			}

			if rl.IsKeyDown(rl.KeyLeftControl) {
				cursorSize = min(max(3, cursorSize+dir*2), 51)
			} else {
				next := cursorCellID + world.CellID(dir)

				if next < world.CELL_VOID {
					cursorCellID = CELL_MAX - 1
				} else if next >= CELL_MAX {
					cursorCellID = world.CELL_VOID
				} else {
					cursorCellID = next
				}
			}
		}

		if rl.IsKeyPressed(rl.KeyP) {
			isPaused = !isPaused
		}

		if rl.IsKeyPressed(rl.KeyS) && isPaused {
			w.Update()
		}

		if !isPaused {
			w.Update()
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawTexture(bgTex, 0, 0, rl.White)
		w.Render()

		rl.DrawText(fmt.Sprintf("Pause: P (%v) | Step: S\nSize: %d\nTile: %s (%d)", isPaused, cursorSize, getCellName(cursorCellID), cursorCellID), 10, 40, 20, rl.RayWhite)

		size := int32(float32(cursorSize*2) * boardScale)
		x, y := w.ToGlobalSpace(lx-cursorSize, ly-cursorSize)
		rl.DrawRectangleLines(int32(x)-1, int32(y)-1, size+2, size+2, rl.White)

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}
