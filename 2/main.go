package main

import (
  "bufio"
  "fmt"
  "os"
  "strings"
  "strconv"
)

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}


func parseInput(filename string) ([][]int, error) {
  file, err := os.Open(filename)

  if err != nil {
    return nil, err
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  
  var res [][]int
  for scanner.Scan() {
    numStrs := strings.Fields(scanner.Text())
    var numsInt []int
    for _, numStr := range numStrs {
      num,  err := strconv.Atoi(numStr)
      
      if err != nil {
        return nil, fmt.Errorf("error converting '%s' to integer: %v", numStr, err) 
      }
      numsInt = append(numsInt, num)
    }
    res = append(res, numsInt)
  }

  if err := scanner.Err(); err != nil {
    return nil, fmt.Errorf("Could not read file: %w", err)
  }

  return res, nil
}

func isSafe(data [][]int) int {
  counter := 0

  for _, row := range data {
    is_decreasing := false
    
    for i := 0; i < len(row) -1; i++ {
      first, second := row[i], row[i+1]
      diff := first - second
      diffAbs := AbsInt(diff)
      
      if i == 0 {
        is_decreasing = (diff < 0)
      }
        
      switch {
      case diff != 0 && is_decreasing == (diff > 0):
      case diffAbs < 1 || diffAbs > 2:
      default:
        counter++
      }
    }
  }
  return counter
}

func main() {
  data, err := parseInput("input.txt")

  if err != nil {
    fmt.Println(err)
    return
  }

  count := isSafe(data)
  fmt.Println("Count: ", count)
}
