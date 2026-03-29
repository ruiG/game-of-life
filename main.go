package main

import (
	"image"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

// Game of Life rules:
// Any live cell with 2-3 neighbors survives
// Any dead cell with exactly 3 neighbors becomes alive
// All other cells die or stay dead

type GameGrid struct {
	rows  int
	cols  int
	cells [][]bool
}

func NewGameGrid(rows int, cols int) *GameGrid {
	game := &GameGrid{
		rows:  rows,
		cols:  cols,
		cells: createEmptyGrid(rows, cols),
	}
	return game
}

func createEmptyGrid(rows int, cols int) [][]bool {
	grid := make([][]bool, rows)
	for i := range cols {
		grid[i] = make([]bool, cols)
	}
	return grid
}

// Count living neighbors
func countNeighbors(cells [][]bool, x int, y int) int {
	count := 0

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue // Skip the cell itself
			}
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < len(cells[0]) && ny >= 0 && ny < len(cells) { // Check bounds
				if cells[ny][nx] { // Note: grid is [row][col] so it's [y][x]
					count++
				}
			}
		}
	}

	return count
}

// Compute next generation
func nextGeneration(currentCells [][]bool) [][]bool {
	newCells := createEmptyGrid(len(currentCells), len(currentCells[0]))
	// Apply rules to each cell
	for y := range currentCells {
		for x := range currentCells[y] {
			switch countNeighbors(currentCells, x, y) {
			case 2:
				// Stay the same
			case 3:
				newCells[y][x] = true // Become alive
			default:
				newCells[y][x] = false // Die or stay dead
			}
		}
	}
	return newCells
}

func main() {
	a := app.New()
	w := a.NewWindow("Game of Life")
	width := 200
	height := 200

	game := NewGameGrid(width, height)

	// Create a raster that draws pixels
	raster := canvas.NewRaster(func(w, h int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, w, h))

		// Draw pixel by pixel
		for x := range w {
			for y := range h {
				if x < game.cols && y < game.rows && game.cells[y][x] { // Note: grid is [row][col] so it's [y][x]
					img.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
				} else {
					img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 255})
				}
			}
		}

		return img
	})

	// Animation loop
	go func() {
		for {
			time.Sleep(time.Millisecond * 33)       // ~30fps
			game.cells = nextGeneration(game.cells) // Compute next state
			fyne.Do(func() {
				raster.Refresh() // Redraw
			})
		}
	}()

	w.SetContent(raster)
	w.Resize(fyne.NewSize(float32(width), float32(height)))
	w.ShowAndRun()
}
