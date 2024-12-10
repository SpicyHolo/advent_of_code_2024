package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// Stores 2d coordinates
type pos struct {
	x, y int
}

func AddPos(p1, p2 pos) pos {
	return pos{p1.x + p2.x, p1.y + p2.y}
}

// Implement a set on a map
type Set[T comparable] map[T]struct{}

func (s Set[T]) add(element T) {
	s[element] = struct{}{}
}

func (s Set[T]) remove(element T) {
	delete(s, element)
}

func (s Set[T]) contains(element T) bool {
	_, exists := s[element]
	return exists
}

// Loads input form text file
func loadInput(filename string) ([][]uint8, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}

	scanner := bufio.NewScanner(file)
	var res [][]uint8
	for scanner.Scan() {
		line := scanner.Text()
		lineHeight := make([]uint8, len(line))
		for i, char := range line {
			lineHeight[i] = uint8(byte(char) - '0')
		}
		res = append(res, lineHeight)
	}

	return res, nil
}

func getHeight(heightMap [][]uint8, position pos) (height uint8) {
	height = heightMap[position.y][position.x]
	return
}

func checkBounds(heightMap [][]uint8, position pos) bool {
	return position.x >= 0 &&
		position.y >= 0 &&
		position.x < len(heightMap[0]) &&
		position.y < len(heightMap)
}

// Get all valid adjacent positions in heightMap, according to trail rules
/*
	The height difference between them should be 1
*/
func getValidAdjacent(heightMap [][]uint8, position pos, diff uint8) []pos {
	var directions = [4]pos{
		{0, -1}, //left
		{1, 0},  // Down
		{0, 1},  // Right
		{-1, 0}, // Up
	}

	var adjNodes []pos

	for _, dir := range directions {
		adjNode := AddPos(position, dir)

		// Check bounds
		if !checkBounds(heightMap, adjNode) {
			continue
		}

		// Check height difference
		if currDiff := getHeight(heightMap, adjNode) - getHeight(heightMap, position); currDiff != diff {
			continue
		}

		adjNodes = append(adjNodes, adjNode)
	}

	return adjNodes
}

// using BFS, counts the number of trails, from starting position
/*
	Part I
	Unique trails are defined as trails with unique destination.
*/
func countUniqueDestTrails(heightMap [][]uint8, start pos) int {
	visited := make(Set[pos])
	numPaths := 0

	queue := make(chan pos, len(heightMap)*len(heightMap))
	queue <- start

	for len(queue) > 0 {
		// Pop from queue
		node := <-queue

		// Check if reached the end of path (height == 9)
		if getHeight(heightMap, node) == 9 {
			numPaths++
			continue
		}

		// Get Valid adjacent
		adjacent := getValidAdjacent(heightMap, node, uint8(1))
		for _, adjNode := range adjacent {
			if !visited.contains(adjNode) {
				visited.add(adjNode)
				queue <- adjNode
			}
		}

	}

	return numPaths
}

func serializePath(path []pos) string {
	var builder strings.Builder
	for _, node := range path {
		fmt.Fprintf(&builder, "(%d, %d),", node.x, node.y)
	}
	return builder.String()
}

// Searches for node in a path
func isInPath(path []pos, node pos) bool {
	for _, path_node := range path {
		if path_node == node {
			return true
		}
	}
	return false
}

// using BFS, counts the number of trails, from starting position
/*
	Part II
	Here unique trails are defined, as the list of traversed nodes (so one destination can have multiple unique trails)
*/
func countUniqueTrails(heightMap [][]uint8, start pos) int {
	numPaths := 0

	queue := make(chan []pos, len(heightMap)*len(heightMap))
	queue <- []pos{start}

	for len(queue) > 0 {
		// Pop from queue
		path := <-queue
		node := path[len(path)-1]

		// Check if reached the end of path (height == 9)
		if getHeight(heightMap, node) == 9 {
			numPaths++
			continue
		}

		// Get Valid adjacent
		adjacent := getValidAdjacent(heightMap, node, uint8(1))
		for _, adjNode := range adjacent {
			if !isInPath(path, adjNode) {
				newPath := append(append([]pos(nil), path...), adjNode)
				queue <- newPath
			}
		}

	}

	return numPaths
}

// Finds all trailhead values, for each finds the number of unique trails, and accumulates the result
func countPathsAtTrailheads(heightMap [][]uint8, trailheadValue uint8, countTrails func(heightMap [][]uint8, start pos) int) (int, time.Duration) {
	startTime := time.Now()
	sum := 0
	for y, row := range heightMap {
		for x, height := range row {
			if height == trailheadValue {
				sum += countTrails(heightMap, pos{x, y})
			}
		}
	}

	return sum, time.Now().Sub(startTime)
}

func part1() {
	// Read input
	heightMap, err := loadInput("./inputs/input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	sum, executionTime := countPathsAtTrailheads(heightMap, 0, countUniqueDestTrails)
	fmt.Printf("sum: %v in %v\n", sum, executionTime)
}

func part2() {
	// Read input
	heightMap, err := loadInput("./inputs/input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	sum, executionTime := countPathsAtTrailheads(heightMap, 0, countUniqueTrails)
	fmt.Printf("sum: %v in %v\n", sum, executionTime)
}

func main() {
	part1()
	part2()
}
