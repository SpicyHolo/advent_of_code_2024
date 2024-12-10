package main

import (
	"bufio"
	"fmt"
	"os"
)

type position struct {
	X, Y int
}

// Parses inputs from a file, returns an map, for each character an array of all position where they were found (and map dimensions.)
func readFile(filename string) (map[byte][]position, int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("could not open file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	y := 0
	maxX := 0
	res := make(map[byte][]position)
	for scanner.Scan() {
		line := scanner.Text()
		// Add each character to a map
		for x, char := range line {
			if char == '.' {
				continue
			}
			res[byte(char)] = append(res[byte(char)], position{x, y})
		}
		y++
		maxX = len(line)
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, 0, fmt.Errorf("error reading file: %w", err)
	}

	maxY := y

	return res, maxX, maxY, nil
}

// Returns two antinodes of two antennas, according to frequency rules (part I)
func pointsOnLine(pos1, pos2 position, maxX, maxY int) []position {
	dx := pos2.X - pos1.X
	dy := pos2.Y - pos1.Y

	points := []position{
		{pos1.X - dx, pos1.Y - dy},
		{pos2.X + dx, pos2.Y + dy},
	}

	var pointsFiltered []position
	for _, pos := range points {
		if pos.X >= 0 && pos.X < maxX && pos.Y >= 0 && pos.Y < maxY {
			pointsFiltered = append(pointsFiltered, pos)
		}
	}
	return pointsFiltered
}

// Returns two antinodes of two antennas, according to frequency rules (part II)
func pointsOnLine2(pos1, pos2 position, maxX, maxY int) []position {
	dx := pos2.X - pos1.X
	dy := pos2.Y - pos1.Y

	// Iterate through points on a line, defined by the input pair of points, one loop for each direction (positive and negative)
	var points []position
	for i := 0; true; i++ {
		point := position{pos1.X - i*dx, pos1.Y - i*dy}
		// Checks bounds
		if point.X < 0 || point.X >= maxX || point.Y < 0 || point.Y >= maxY {
			break
		}
		points = append(points, point)
	}

	for i := 0; true; i++ {
		point := position{pos2.X + i*dx, pos2.Y + i*dy}
		if point.X < 0 || point.X >= maxX || point.Y < 0 || point.Y >= maxY {
			break
		}
		points = append(points, point)
	}

	return points
}

// Given an array of T, return an array of all possible pairs (permutations)
func getPairs[T any](values []T) [][2]T {
	var pairs [][2]T
	n := len(values)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			pairs = append(pairs, [2]T{values[i], values[j]})
		}
	}
	return pairs
}

func part1(antennas map[byte][]position, maxX, maxY int) {
	uniquePointsSet := make(map[position]struct{})
	// grid := grid.NewTextGrid(maxX, maxY)
	for _, v := range antennas {
		pairs := getPairs(v)
		for _, pair := range pairs {
			new_points := pointsOnLine(pair[0], pair[1], maxX, maxY)
			for _, point := range new_points {
				uniquePointsSet[point] = struct{}{}
				// grid.SetChar(point.X, point.Y, '#')
			}
		}
	}

	// grid.Print() // print antinodes for debuggin
	fmt.Println("part1: ", len(uniquePointsSet))
}

func part2(antennas map[byte][]position, maxX, maxY int) {
	uniquePointsSet := make(map[position]struct{})
	// grid := grid.NewTextGrid(maxX, maxY)
	for _, v := range antennas {
		pairs := getPairs(v)
		for _, pair := range pairs {
			new_points := pointsOnLine2(pair[0], pair[1], maxX, maxY)
			for _, point := range new_points {
				uniquePointsSet[point] = struct{}{}
				// grid.SetChar(point.X, point.Y, '#')
			}
		}
	}

	// grid.Print() // print antinodes for debuggin
	fmt.Println("part2: ", len(uniquePointsSet))
}
func main() {
	// Parse inputs
	antennas, maxX, maxY, err := readFile("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	part1(antennas, maxX, maxY)
	part2(antennas, maxX, maxY)

}
