package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Vec struct {
	X, Y int
}

func (v1 Vec) Add(v2 Vec) Vec {
	return Vec{v1.X + v2.X, v1.Y + v2.Y}
}

type Item struct {
	value    Vec
	priority int
	index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value Vec, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

// Find index of item in a priority queue, returns -1, false if not found.
func (pq PriorityQueue) Contains(pos Vec) (int, bool) {
	for i, item := range pq {
		if item.value.X == pos.X && item.value.Y == pos.Y {
			return i, true
		}
	}
	return -1, false
}

// Parses input for all corrupted memory
func parseInput(filename string) ([]Vec, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var data []Vec
	for scanner.Scan() {
		line := string(scanner.Text())
		if line == "" {
			continue
		}
		digits := strings.Split(line, ",")
		if len(digits) != 2 {
			return nil, fmt.Errorf("invalid number of arguments, on line when reading input")
		}

		// Convert each digit to int
		x, err1 := strconv.Atoi(digits[0])
		if err1 != nil {
			return nil, fmt.Errorf("could not convert %v to int: %w", digits[0], err1)
		}
		y, err2 := strconv.Atoi(digits[1])
		if err2 != nil {
			return nil, fmt.Errorf("could not convert %v to int: %w", digits[1], err2)
		}

		vec := Vec{x, y}
		data = append(data, vec)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return data, nil
}

// Get adjacent nodes, checking for bounds and if are not corrupted
func getAdj(curPos Vec, corrupted map[Vec]struct{}, mapSize Vec) []Vec {
	directions := []Vec{
		{0, -1},
		{1, 0},
		{0, 1},
		{-1, 0},
	}

	var adjacent []Vec
	for _, dir := range directions {
		newPos := curPos.Add(dir)

		// Check if new position is in bounds
		if newPos.X < 0 || newPos.Y < 0 || newPos.X >= mapSize.X || newPos.Y >= mapSize.Y {
			continue
		}

		// Check if new position is not corrupted
		if _, contains := corrupted[newPos]; contains {
			continue
		}

		adjacent = append(adjacent, newPos)
	}

	return adjacent
}

// A star heuristic, manhattan distance from end
func heuristic(curPos, end Vec) int {
	dx := AbsInt(curPos.X - end.X)
	dy := AbsInt(curPos.Y - end.Y)
	return dx + dy
}

// Traverse the best_path map, to reconstruct the best path
func reconstructPath(start, end Vec, best_path map[Vec]Vec) []Vec {
	path := []Vec{end}
	curPos := end

	for curPos.X != start.X || curPos.Y != start.Y {
		curPos = best_path[curPos]
		path = append(path, curPos)
	}
	slices.Reverse(path)
	return path
}

// Search for the shortest path in a graph, with the help of a heuristic
func AStar(start, end Vec, occupied map[Vec]struct{}, mapSize Vec) ([]Vec, int) {
	// Initalise variables
	visited := make(map[Vec]struct{})
	best_path := make(map[Vec]Vec)
	best_cost := make(map[Vec]int)

	// Create queue
	pq := make(PriorityQueue, 0)

	pq.Push(&Item{
		value:    start,
		priority: heuristic(start, end),
	})

	for len(pq) > 0 {
		// Fetch new node from queue
		node := heap.Pop(&pq).(*Item)
		curPos := node.value
		curCost := best_cost[curPos]

		// Add to visited
		visited[curPos] = struct{}{}

		// Check if reached the end
		if curPos == end {
			return reconstructPath(start, end, best_path), best_cost[end]
		}

		// Visit adjacents
		for _, adj := range getAdj(curPos, occupied, mapSize) {
			newCost := curCost + 1
			newPriority := newCost + heuristic(adj, end)

			// Skip if already visited
			if _, exists := visited[adj]; exists {
				continue
			}

			// See if already in queue, check if found new best path to that node, if so update it.
			if i, exists := pq.Contains(adj); exists {
				if pq[i].priority > newPriority {
					pq.update(pq[i], adj, newPriority)
					best_path[adj] = curPos
					best_cost[adj] = newCost
				}
				continue
			}

			// Elsee, add new node to queue
			item := &Item{
				value:    adj,
				priority: newPriority,
			}
			heap.Push(&pq, item)
			best_path[adj] = curPos
			best_cost[adj] = newCost
		}
	}

	// No path found
	return nil, -1
}

func example() {
	data, err := parseInput("test_input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create set of n first corrupted memories
	corrupted := make(map[Vec]struct{})
	for _, v := range data[:12] {
		corrupted[v] = struct{}{}
	}

	// Search for path
	_, length := AStar(Vec{0, 0}, Vec{6, 6}, corrupted, Vec{7, 7})
	fmt.Println("example: ", length)
}

func part1() {
	data, err := parseInput("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create set of n first
	corrupted := make(map[Vec]struct{})
	for _, v := range data[:1024] {
		corrupted[v] = struct{}{}
	}

	_, length := AStar(Vec{0, 0}, Vec{70, 70}, corrupted, Vec{71, 71})
	fmt.Println("part I: ", length)
}

func part2() {
	data, err := parseInput("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create set of n first
	corrupted := make(map[Vec]struct{})

	for _, v := range data[:1024] {
		corrupted[v] = struct{}{}
	}

	// Initialise path, create a set for all positions in the path
	path, _ := AStar(Vec{0, 0}, Vec{70, 70}, corrupted, Vec{71, 71})
	pathSet := make(map[Vec]struct{})
	for _, pos := range path {
		pathSet[pos] = struct{}{}
	}

	fmt.Println("Part II...")

	for i := 1024; i < len(data); i++ { // First 1024 are already safe, we checked it.
		corrupted[data[i]] = struct{}{}

		// If new corrupted memory is not in path, continue
		if _, exists := pathSet[data[i]]; !exists {
			continue
		}

		// Recalculate path
		path, _ := AStar(Vec{0, 0}, Vec{70, 70}, corrupted, Vec{71, 71})

		if path == nil {
			fmt.Println("No path found, after corrupting:", data[i])
			return
		}

		// Recalculate in path set
		pathSet = make(map[Vec]struct{})
		for _, pos := range path {
			pathSet[pos] = struct{}{}
		}

	}
}

func main() {
	example()
	part1()
	part2()
}
