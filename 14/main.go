package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"regexp"
	"strconv"
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

// Vector type
type Vec2D struct {
	X, Y int
}

func (v1 *Vec2D) Add(v2 Vec2D) Vec2D {
	return Vec2D{v1.X + v2.X, v1.Y + v2.Y}
}

func (v *Vec2D) Scale(a int) Vec2D {
	return Vec2D{a * v.X, a * v.Y}
}

func (v *Vec2D) Mod(w Vec2D) Vec2D {
	return Vec2D{
		(v.X%w.X + w.X) % w.X,
		(v.Y%w.Y + w.Y) % w.Y,
	}
}

func (v Vec2D) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "(%d, %d)", v.X, v.Y)
	return builder.String()
}

// Robot type
type Robot struct {
	ID   int
	P, V Vec2D
}

func (r Robot) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Robot{ID=%d, p=%v, v=%v}", r.ID, r.P, r.V)
	return builder.String()
}

// Reading inputs
func readInput(filename string) ([]Robot, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	// Parse input
	data := string(buf)

	// ex. regex line: p=74,25 v=-62,4
	re := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)
	matches := re.FindAllStringSubmatch(data, -1)

	robots := make([]Robot, len(matches))
	for id, match := range matches {
		x, err1 := strconv.Atoi(match[1])
		if err1 != nil {
			return nil, fmt.Errorf("error parsing line %v: %w", match, err)
		}
		y, err2 := strconv.Atoi(match[2])
		if err2 != nil {
			return nil, fmt.Errorf("error parsing line %v: %w", match, err)
		}

		vx, err3 := strconv.Atoi(match[3])
		if err3 != nil {
			return nil, fmt.Errorf("error parsing line %v: %w", match, err)
		}
		vy, err4 := strconv.Atoi(match[4])
		if err4 != nil {
			return nil, fmt.Errorf("error parsing line %v: %w", match, err)
		}

		// Create robot
		robots[id] = Robot{
			ID: id,
			P:  Vec2D{X: x, Y: y},
			V:  Vec2D{X: vx, Y: vy},
		}
	}
	return robots, nil
}

func moveRobot(r *Robot, mapSize Vec2D) {
	pos, vel := r.P, r.V
	newPos := pos.Add(vel)
	r.P = newPos.Mod(mapSize)
}

func findQuadrant(r *Robot, mapSize Vec2D) int {
	middle := Vec2D{mapSize.X / 2, mapSize.Y / 2} // Adjust middle point to be center, not one step offset

	pos := r.P

	if pos.X > middle.X && pos.Y < middle.Y {
		// Top Right
		return 0
	}
	if pos.X > middle.X && pos.Y > middle.Y {
		// Bottom Right
		return 1
	}
	if pos.X < middle.X && pos.Y > middle.Y {
		// Bottom Left
		return 2
	}
	if pos.X < middle.X && pos.Y < middle.Y {
		// Top Left
		return 3
	}
	return -1
}

func simulateRobots(robots []Robot, mapSize Vec2D, n_steps int) {
	for i := 0; i < n_steps; i++ {
		for id := range robots {
			moveRobot(&robots[id], mapSize)
		}
	}
}

func getSafety(robots []Robot, mapSize Vec2D) int {
	quads := make([]int, 4)
	for id := range robots {
		q := findQuadrant(&robots[id], mapSize)
		if q == -1 {
			continue
		}
		quads[q]++
	}

	return quads[0] * quads[1] * quads[2] * quads[3]
}

func createImage(robots []Robot, mapSize Vec2D) image.Image {
	w, h := mapSize.X, mapSize.Y
	img := image.NewRGBA(image.Rect(0, 0, w, h))

	black := color.RGBA{0, 0, 0, 255}
	green := color.RGBA{0, 255, 0, 255}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, black)
		}
	}

	for _, robot := range robots {
		x, y := robot.P.X, robot.P.Y
		img.Set(x, y, green)
	}

	return img
}

func part1() {
	robots, err := readInput("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	mapSize := Vec2D{101, 103}
	simulateRobots(robots, mapSize, 100)
	fmt.Println("Safety: ", getSafety(robots, mapSize))
}

func part2() {
	robots, err := readInput("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	mapSize := Vec2D{101, 103}

	for i := 0; i < 101*103; i++ {
		simulateRobots(robots, mapSize, 1)
		// safety := getSafety(robots, mapSize)
		fileName := fmt.Sprintf("./dst/%d.png", i)
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Println("Error creating file: ", err)
		}
		defer file.Close()

		img := createImage(robots, mapSize)
		err = png.Encode(file, img)
		if err != nil {
			fmt.Println("Error encoding image:", err)
			return
		}
	}
}
func main() {
	part2()
}
