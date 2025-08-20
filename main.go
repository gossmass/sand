package main

import (
	"fmt"
	"sand/world"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	TILE_SAND world.TileID = iota + 1
	TILE_WATER

	TILE_MAX
)

func sandRule(w *world.World, x, y int) {
	if w.Get(x, y+1) == world.TILE_VOID {
		w.Swap(x, y, x, y+1)
	} else if w.Get(x-1, y+1) == world.TILE_VOID {
		w.Swap(x, y, x-1, y+1)
	} else if w.Get(x+1, y+1) == world.TILE_VOID {
		w.Swap(x, y, x+1, y+1)
	}
}

func waterRule(w *world.World, x, y int) {
	if w.Get(x, y+1) == world.TILE_VOID {
		w.Swap(x, y, x, y+1)
	} else {
		if w.Get(x-1, y) == world.TILE_VOID {
			w.Swap(x, y, x-1, y)
		}
		if w.Get(x+1, y) == world.TILE_VOID {
			w.Swap(x, y, x+1, y)
		}
	}
}

func main() {
	var boardSize int32 = 256
	var boardScale float32 = 2.0
	var cursorSize int = 3
	var cursorTileID world.TileID = TILE_SAND

	rl.InitWindow(int32(float32(boardSize)*boardScale), int32(float32(boardSize)*boardScale), "Title")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	w := world.New(int(boardSize), int(boardSize), boardScale)
	w.AddRule(world.NewRule(sandRule, rl.Green), TILE_SAND)
	w.AddRule(world.NewRule(waterRule, rl.Blue), TILE_WATER)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		mx := int(rl.GetMouseX())
		my := int(rl.GetMouseY())
		lx, ly := w.ToLocalSpace(mx, my)

		if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			w.PutBlob(lx, ly, cursorSize, cursorTileID)
		}

		if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
			if cursorTileID+1 == TILE_MAX {
				cursorTileID = world.TILE_VOID
			} else {
				cursorTileID++
			}
		}

		if w := rl.GetMouseWheelMove(); w != 0 {
			if w < 0 {
				cursorSize = max(3, cursorSize-1)
			} else {
				cursorSize += 2
			}
		}

		w.Update()
		w.Render()

		rl.DrawText(fmt.Sprintf("Size: %d\nTile: %d", cursorSize, cursorTileID), 10, 40, 20, rl.Black)

		size := int32(float32(cursorSize*2) * boardScale)
		x, y := w.ToGlobalSpace(lx-cursorSize, ly-cursorSize)
		rl.DrawRectangleLines(int32(x), int32(y), size, size, rl.Red)

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}

	w.Destroy()
}
