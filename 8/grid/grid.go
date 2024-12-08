package grid

import (
	"fmt"
	"strings"
)

// TextGrid represents a 2D grid of characters.
type TextGrid struct {
	width, height int
	grid          [][]rune
}

// NewTextGrid creates a new TextGrid with specified width and height.
func NewTextGrid(width, height int) *TextGrid {
	grid := make([][]rune, height)
	for i := range grid {
		grid[i] = make([]rune, width)
		for j := range grid[i] {
			grid[i][j] = '.' // Initialize with spaces
		}
	}
	return &TextGrid{width: width, height: height, grid: grid}
}

// Print displays the TextGrid.
func (tg *TextGrid) Print() {
	for _, row := range tg.grid {
		fmt.Println(strings.Join(strings.Split(string(row), ""), ""))
	}
}

// SetChar sets a character at the given (x, y) position.
func (tg *TextGrid) SetChar(x, y int, char rune) error {
	if x < 0 || x >= tg.width || y < 0 || y >= tg.height {
		return fmt.Errorf("coordinates out of bounds")
	}
	tg.grid[y][x] = char
	return nil
}
