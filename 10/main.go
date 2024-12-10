package main

import (
  "os"
  "bufio"
  "fmt"
)


func loadInput(filename string) ([][] byte, error){
  file, err := os.Open(filename)  
  if err != nil {
    return nil, fmt.Errorf("could not open file: %w", err)
  }
  
  scanner := bufio.NewScanner(file)
  var res [][] byte
  for scanner.Scan() {
    line := scanner.Text()
    lineHeight := make([] byte, len(line))
    for i, char := range line {
      lineHeight[i] = byte(char) - '0'
    }
    res = append(res, lineHeight)
  }
  
  return res, nil
} 

func checkBounds(heightMap [][] byte, x, y int) bool {
    return x >= 0 && 
           y >= 0 && 
           x < len(heightMap[0]) && 
           y < len(heightMap)
}

func getValidAdjacent(heightMap [][] byte, x, y, diff int) [][2]int {
  var directions = [4][2]int{
    {0, -1},  //left 
    {1, 0},   // Down
    {0, 1},   // Right
    {-1, 0}, // Up
  } 

  var adjNodes [][2] int

  for _, dir := range directions {
    x_adj := x + dir[0]
    y_adj := y + dir[1]
    
    // Check bounds
    if !checkBounds(heightMap, x_adj, y_adj) {
      continue
    }
    
    // Check height difference
    if heightMap[x_adj][y_adj] - heightMap[x][y] > byte(diff) {
      continue
    }
    
    node := [2]int{x_adj, y_adj}
    adjNodes = append(adjNodes, node)
  }
  
  return adjNodes
}


func main() {
  heightMap, err := loadInput("test_input.txt")
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Println(heightMap)
}
