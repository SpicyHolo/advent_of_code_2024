package main

import (
	"bytes"
	"container/heap"
	"container/list"
	"fmt"
	"os"
	"slices"
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

func (v1 Vec) AddVec(v2 Vec) Vec {
	return Vec{v1.X + v2.X, v1.Y + v2.Y}
}

// Manhattan dist betwen two vectors
func dist(v1, v2 Vec) int {
	dx := AbsInt(v1.X - v2.X)
	dy := AbsInt(v1.Y - v2.Y)
	return dx + dy
}

type PriorityQueue []*Item

type Item struct {
	value    Vec
	priority int
	index    int
}

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
func (pq PriorityQueue) Contains(s Vec) (int, bool) {
	for i, item := range pq {
		if item.value == s {
			return i, true
		}
	}
	return -1, false
}

// Parses input for all corrupted memory
func parseInput(filename string) ([][]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	occupancyGrid := bytes.Split(data, []byte("\n"))
	return occupancyGrid, nil
}

// Get adjacent nodes, checking for bounds and if are not corrupted, also tries to cheat
func getAdj(curVec Vec, occupancy [][]byte) []Vec {
	mapSize := Vec{len(occupancy[0]), len(occupancy)}
	directions := []Vec{
		{0, -1},
		{1, 0},
		{0, 1},
		{-1, 0},
	}

	var adjacent []Vec

	for _, dir := range directions {
		newVec := curVec.AddVec(dir)

		// Check if new position is in bounds
		if newVec.X < 0 || newVec.Y < 0 || newVec.X >= mapSize.X || newVec.Y >= mapSize.Y {
			continue
		}

		// Check if new position is a wall
		if occupancy[newVec.Y][newVec.X] == '#' {
			continue
		}

		adjacent = append(adjacent, newVec)
	}

	return adjacent
}

// Traverse the best_path map, to reconstruct the best path (as a slice of positions on the path)
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
// MaxCost can be -1, to return any shorest path or a positive value, to set maximum allowed path cost
func Djikstra(start, end Vec, occupancy [][]byte, maxCost int) ([]Vec, int) {
	// Initalise variables
	visited := make(map[Vec]struct{})
	best_path := make(map[Vec]Vec)
	best_cost := make(map[Vec]int)

	// Create queue
	pq := make(PriorityQueue, 0)

	pq.Push(&Item{
		value:    start,
		priority: 0,
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
			return reconstructPath(start, curPos, best_path), best_cost[curPos]
		}

		// Visit adjacents
		for _, adj := range getAdj(curPos, occupancy) {
			newCost := curCost + 1

			// Skip if already visited
			if _, exists := visited[adj]; exists {
				continue
			}

			if newCost > maxCost && maxCost != -1 {
				continue
			}
			// See if already in queue, check if found new best path to that node, if so update it.
			if i, exists := pq.Contains(adj); exists {
				if pq[i].priority > newCost {
					pq.update(pq[i], adj, newCost)
					best_path[adj] = curPos
					best_cost[adj] = newCost
				}
				continue
			}

			// Elsee, add new node to queue
			item := &Item{
				value:    adj,
				priority: newCost,
			}
			heap.Push(&pq, item)
			best_path[adj] = curPos
			best_cost[adj] = newCost
		}
	}

	// No path found
	return nil, -1
}

// BFS, returns the map of all reachable positionsfrom the start as keys, and distance to them as values
func BFS(start Vec, occupancyGrid [][]byte) map[Vec]int {
	visited := make(map[Vec]struct{})

	queue := list.New()

	visited[start] = struct{}{}
	queue.PushBack(start)

	costMap := make(map[Vec]int)
	costMap[start] = 0
	for queue.Len() > 0 {
		element := queue.Front()
		curPos := element.Value.(Vec)
		queue.Remove(element)

		for _, adj := range getAdj(curPos, occupancyGrid) {
			newCost := costMap[curPos] + 1
			if _, seen := visited[adj]; seen {
				continue
			}
			costMap[adj] = newCost
			visited[adj] = struct{}{}
			queue.PushBack(adj)
		}
	}
	return costMap
}

// Find start and end positions from the input grid
func findStartEnd(occupancyGrid [][]byte) (Vec, Vec) {
	var start, end Vec
	for y := 0; y < len(occupancyGrid); y++ {
		for x := 0; x < len(occupancyGrid[0]); x++ {
			if pos := occupancyGrid[y][x]; pos == 'S' {
				occupancyGrid[y][x] = '.'
				start = Vec{x, y}
			} else if pos == 'E' {
				occupancyGrid[y][x] = '.'
				end = Vec{x, y}
			}
		}
	}
	return start, end
}

func part1() {
	if len(os.Args) != 2 {
		return
	}

	// Load Input
	occupancyGrid, err := parseInput(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	/* First find the base path length */
	// Find start, end
	start, end := findStartEnd(occupancyGrid)

	_, path_length := Djikstra(start, end, occupancyGrid, -1)

	if path_length == -1 {
		return
	}

	/* Now find all cheating paths, that save at least 100 picoseconds */
	maxCost := path_length - 100
	numPaths := 0

	// Try removing each wall, and check if path in the new map is short enough.
	for y := 1; y < len(occupancyGrid)-1; y++ {
		for x := 1; x < len(occupancyGrid[0])-1; x++ {
			if occupancyGrid[y][x] == '#' {
				occupancyGrid[y][x] = '.'
				_, path_length := Djikstra(start, end, occupancyGrid, maxCost)
				if path_length != -1 {
					numPaths++
				}
				occupancyGrid[y][x] = '#'
			}
		}
	}
	fmt.Println("Part I: ", numPaths)
}

func part2() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go file_path")
		return
	}

	// Load Input
	occupancyGrid, err := parseInput(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	/* First find the base path length */
	// Find start, end
	start, end := findStartEnd(occupancyGrid)

	_, path_length := Djikstra(start, end, occupancyGrid, -1)

	if path_length == -1 {
		return
	}

	/* Now find all cheating paths, that save at least 100 picoseconds */
	maxCost := path_length - 100
	numPaths := 0

	// with BFS get all posisble reachable positions from start and end, with their cost
	fromstart := BFS(start, occupancyGrid)
	fromend := BFS(end, occupancyGrid)

	// For all possible combinations of those, check if manhattan distance between them is less than cheating time.
	// If the total distance with cheating is less than maximum allowable cost, count it as a solution.
	// @Neil Thistlethwaite
	for pos1 := range fromstart {
		for pos2 := range fromend {
			if d := dist(pos1, pos2); d <= 20 {
				if fromstart[pos1]+d+fromend[pos2] <= maxCost {
					numPaths++
				}
			}
		}
	}

	fmt.Println("Part II:", numPaths)
}

func main() {
	part2()
	fmt.Println("calculating part I, this may take a while!")
	part1() // Takes a while!
}
