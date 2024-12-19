package main

import (
  "bufio"
  "os"
  "strings"
  "strconv"
  "fmt"
)
type Vec struct {
  X, Y int
}

type Item struct {
  value Vec
  priority int
  index int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
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
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}



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


func main() {
  data, err := parseInput("input.txt") 
  if err != nil {
    fmt.Println(err)
    return
  }

  // Create set of n first 
  corrupted := make(map[Vec] struct{})
  for _, v := range data[:1000] {
    corrupted[v] = struct{}{}
  }
}
