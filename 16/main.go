package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
	"strings"
)

func AbsInt(a int) int {
	if a < 0 {
		a = -a
	}
	return a
}

// Direction constants
const (
	NORTH = 0
	EAST  = 1
	SOUTH = 2
	WEST  = 3
)

// Direction vectors
var DIRECTIONS = map[int]Vector{
	NORTH: {0, -1},
	EAST:  {1, 0},
	SOUTH: {0, 1},
	WEST:  {-1, 0},
}

// Labirynth structure
type Labirynth struct {
	Map        [][]byte
	Start, End Vector
}

func (lab Labirynth) String() string {
	var builder strings.Builder
	rows := bytes.Join(lab.Map, []byte("\n"))
	builder.WriteString(string(rows))
	return builder.String()
}

func isFree(lab Labirynth, pos Vector) bool {
	return lab.Map[pos.Y][pos.X] == '.'
}

// Position vector
type Vector struct {
	X, Y int
}

func (v1 Vector) Add(v2 Vector) Vector {
	return Vector{v1.X + v2.X, v1.Y + v2.Y}
}

func (v1 Vector) Sub(v2 Vector) Vector {
	return Vector{v1.X - v2.X, v1.Y - v2.Y}
}

// State represents an entity's position and direction
type State struct {
	Pos Vector
	Dir int
}

// Priority Queue for Dijkstra's
type Item struct {
	Value    State
	Parents  []*Item // Array of all parent nodes (to keep the shortest way to get to current node)
	Priority int     // this is the cost in Dijkstra's algorithm
	Index    int
}

// Smaller priority is first.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *Item, parents []*Item, cost int) {
	item.Parents = parents
	item.Priority = cost
	heap.Fix(pq, item.Index)
}

func (pq PriorityQueue) Contains(value State) (int, bool) {
	for i, item := range pq {
		if item.Value.Pos == value.Pos && item.Value.Dir == value.Dir {
			return i, true
		}
	}
	return -1, false
}

// Set implementation
type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(element T) {
	s[element] = struct{}{}
}

func (s Set[T]) Contains(element T) bool {
	_, exists := s[element]
	return exists
}

// Reads input from file
func readInput(filename string) (Labirynth, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Labirynth{}, fmt.Errorf("could not read file: %w", err)
	}
	Map := bytes.Split(data, []byte("\r\n"))

	start, end := Vector{}, Vector{}

	// Find starting and end position
	for y := 0; y < len(Map); y++ {
		for x := 0; x < len(Map[y]); x++ {
			switch Map[y][x] {
			case 'S':
				start.X, start.Y = x, y
				Map[y][x] = '.'
			case 'E':
				end.X, end.Y = x, y
				Map[y][x] = '.'
			}
		}
	}

	// Create labirynth object
	l := Labirynth{
		Map:   Map,
		Start: start,
		End:   end,
	}

	return l, nil
}

// A* Helpers
func getCost(state1, state2 State) (cost int) {
	pos1, pos2 := state1.Pos, state2.Pos
	dir1, dir2 := state1.Dir, state2.Dir

	dp := pos1.Sub(pos2)
	dx, dy := AbsInt(dp.X), AbsInt(dp.Y)

	dirDiff := AbsInt(dir1 - dir2)

	// Check if jumped rotations (distance between dir=0, dir=3 should be 1, since 3 is next to 0)
	if dirDiff == len(DIRECTIONS)-1 {
		dirDiff = 1
	}

	cost = dx + dy + dirDiff*1000
	return
}

func getAdjacent(lab Labirynth, state State) []State {
	var adj []State
	pos, dir := state.Pos, state.Dir

	newStates := map[string]State{
		"fw":  {pos.Add(DIRECTIONS[dir]), dir},
		"cw":  {pos, (dir + 1) % len(DIRECTIONS)},
		"ccw": {pos, (dir - 1 + len(DIRECTIONS)) % len(DIRECTIONS)},
	}

	for _, newState := range newStates {
		if isFree(lab, newState.Pos) {
			adj = append(adj, newState)
		}
	}

	return adj
}

func Dijkstra(lab Labirynth) (int, [][]State) {
	// Initialise priority queue
	pq := make(PriorityQueue, 0)

	var ends []*Item // Store all possible end states (4 orientations)
	var allPaths [][]State

	// Keep set of visited nodes
	visited := make(Set[State])

	minCost := -1

	// Add start node to queue, default facing EAST
	start_state := State{lab.Start, EAST}

	heap.Push(&pq, &Item{
		Value:    start_state,
		Priority: 0,
	})

	// Process the queue
	for pq.Len() > 0 {
		// Fetch next item
		cur_item := heap.Pop(&pq).(*Item)
		cur_state := cur_item.Value
		cur_cost := cur_item.Priority

		// Add as visited
		visited.Add(cur_state)

		// Path is Found!
		if cur_state.Pos == lab.End {
			if minCost != -1 && cur_cost > minCost {
				break
			}
			minCost = cur_cost
			ends = append(ends, cur_item)
		}

		// Add adjacent states to the queue
		for _, adj := range getAdjacent(lab, cur_state) {
			cost := cur_cost + getCost(cur_state, adj)

			// Skip if already in visited list
			if visited.Contains(adj) {
				continue
			}

			// If in queue, check if new path is better or equal
			if i, ok := pq.Contains(adj); ok {
				if pq[i].Priority > cost {
					// Update path if found a better one
					pq.update(pq[i], []*Item{cur_item}, cost)
				} else if pq[i].Priority == cost { // Add multiple parents, if few paths have the same cost
					new_parents := append(pq[i].Parents, cur_item)
					pq.update(pq[i], new_parents, cost)
				}
				continue
			}

			// If not in the queue, create new item
			item := &Item{
				Value:    adj,
				Parents:  []*Item{cur_item},
				Priority: cost,
			}
			heap.Push(&pq, item)

		}
	}

	// Reconstruct all paths for the end state
	if minCost != -1 {
		for _, end := range ends {
			allPaths = append(allPaths, reconstructPaths(end)...)
		}
	}

	return minCost, allPaths
}

// Reconstruct paths from stored best paths (as parent nodes)
// need DFS to generate all possible paths (since some nodes mayh have multiple parents)
func reconstructPaths(endItem *Item) [][]State {
	var result [][]State
	var dfs func(*Item, []State)

	// Recursive DFS function
	dfs = func(current *Item, path []State) {
		// Prepend the current item's value to the path
		newPath := append([]State{current.Value}, path...)

		// If there are no parents, we've reached the beginning of a path
		if len(current.Parents) == 0 {
			result = append(result, newPath)
			return
		}

		// Recurse into each parent
		for _, parent := range current.Parents {
			dfs(parent, newPath)
		}
	}

	// Start DFS from the end item
	dfs(endItem, []State{})

	return result
}

func solve(filename string) {
	// Load input
	labirynth, err := readInput(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	cost, allPaths := Dijkstra(labirynth)
	if cost == -1 {
		fmt.Println("No path found.")
		return
	}

	// Update the map with the paths found
	seats := 0
	for _, path := range allPaths {
		for _, state := range path {
			x, y := state.Pos.X, state.Pos.Y
			if labirynth.Map[y][x] != 'x' {
				seats++
			}
			labirynth.Map[y][x] = 'x'
		}
	}

	fmt.Println("Labirynth with paths:")
	fmt.Println(labirynth)
	fmt.Println("Cost: ", cost)
	fmt.Println("Seats: ", seats)
}

func main() {
	solve("test_input.txt")
}
