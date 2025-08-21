package main

import (
	"fmt"
	"sand/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	CELL_SAND world.CellID = iota + 1
	CELL_WATER

	CELL_MAX
)

func sandRule(w *world.World, x, y int) {
	if w.Get(x, y+1) == world.CELL_VOID {
		w.Swap(x, y, x, y+1)
	} else if w.Get(x-1, y+1) == world.CELL_VOID {
		w.Swap(x, y, x-1, y+1)
	} else if w.Get(x+1, y+1) == world.CELL_VOID {
		w.Swap(x, y, x+1, y+1)
	}
}

func waterRule(w *world.World, x, y int) {
	if w.Get(x, y+1) == world.CELL_VOID {
		w.Swap(x, y, x, y+1)
	} else {
		if w.Get(x-1, y) == world.CELL_VOID {
			w.Swap(x, y, x-1, y)
		}
		if w.Get(x+1, y) == world.CELL_VOID {
			w.Swap(x, y, x+1, y)
		}
	}
}

func getCellName(cellID world.CellID) string {
	switch cellID {
	case world.CELL_VOID:
		return "VOID"
	case CELL_SAND:
		return "SAND"
	case CELL_WATER:
		return "WATER"
	}

	return "Undefined"
}

func main() {
	var boardWidth int32 = 128
	var boardHeight int32 = 96
	var boardScale float32 = 8.0
	var cursorSize int = 3
	var cursorCellID world.CellID = CELL_SAND

	rl.InitWindow(int32(float32(boardWidth)*boardScale), int32(float32(boardHeight)*boardScale), "Title")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	w := world.New(int(boardWidth), int(boardHeight), boardScale)
	defer w.Destroy()

	// w.AddRule(world.NewRule(sandRule, rl.Green), CELL_SAND)
	// w.AddRule(world.NewRule(waterRule, rl.Blue), CELL_WATER)
	
	bgImg := rl.GenImageChecked(int(float32(boardWidth) * boardScale), int(float32(boardHeight) * boardScale), int(boardScale) * 4, int(boardScale) * 4, rl.DarkGray, rl.Black)
	bgTex := rl.LoadTextureFromImage(bgImg)
	rl.UnloadImage(bgImg)
	defer rl.UnloadTexture(bgTex)

	for !rl.WindowShouldClose() {

		mx := int(rl.GetMouseX())
		my := int(rl.GetMouseY())
		lx, ly := w.ToLocalSpace(mx, my)

		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			w.PutBlob(lx, ly, cursorSize, cursorCellID)
		}

		if w := rl.GetMouseWheelMove(); w != 0 {
			dir := 1
			if w < 0 {
				dir = -1
			}

			if rl.IsKeyDown(rl.KeyLeftControl) {
				cursorSize = min(max(3, cursorSize + dir * 2), 51)
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

		w.Update()
		
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawTexture(bgTex, 0, 0, rl.White)
		w.Render()

		rl.DrawText(fmt.Sprintf("Size: %d\nTile: %s (%d)", cursorSize, getCellName(cursorCellID), cursorCellID), 10, 40, 20, rl.RayWhite)

		size := int32(float32(cursorSize*2) * boardScale)
		x, y := w.ToGlobalSpace(lx-cursorSize, ly-cursorSize)
		rl.DrawRectangleLines(int32(x) - 1, int32(y) - 1, size + 2, size + 2, rl.White)

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}
