package main

import (
	"bufio" // For reading file line by line
	"fmt"   // For printing
	"os"    // For opening files
  "regexp"
  "strconv"
  "strings"
)

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Reads an input file, returns an array of strings for each line
func loadInput(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
  var res []string
	for scanner.Scan() {
    res = append(res, scanner.Text())
  }

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("Could not read file: %w", err)
	}

	return strings.Join(res, "\n"), nil
}

func parseMatches(matches [][]string) (int, error) {
  var res int
  for _, match := range matches {
    if len(match) < 3 {
      return 0, fmt.Errorf("Not a valid match, should have 3 groups")
    }
    
    first, err1 := strconv.Atoi(match[1])
    second, err2 := strconv.Atoi(match[2])
    if err1 != nil {
      return 0, fmt.Errorf("Couldn't convert one of the operands to int: %v", err1)
    }
    
    if err2 != nil {
      return 0, fmt.Errorf("Couldn't convert one of the operands to int: %v", err2)
    }

    res += first * second 
  }
  return res, nil
}

func part_one(data string) {
  var validMultiExpr = regexp.MustCompile(`mul\((\d+),(\d+)\)`)
  
  // Parse each line and regex, return one array of all matches
  matches := validMultiExpr.FindAllStringSubmatch(data, -1)

  // Parsing matches, calculate the sum from all operations 
  res, err := parseMatches(matches)
  if err != nil {
      fmt.Println("Error", err)
  }

  // Print solution
  fmt.Println("Sum: ", res)
}

// Given matched mul(int, int) statement, return the operation result
func calculate_mul(match []string) (int, error) {
  if len(match) != 3 {
    return 0, fmt.Errorf("Not enough elements in found mul(int, int) match.")
  }

  first, err1 := strconv.Atoi(match[1])
  second, err2 := strconv.Atoi(match[2])

  if err1 != nil {
    return 0, fmt.Errorf("Couldn't convert one of the operands to int: %v", err1)
  }
  
  if err2 != nil {
    return 0, fmt.Errorf("Couldn't convert one of the operands to int: %v", err2)
  }

  return first * second, nil
}

func part_two(data string) {
	mulRegex := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	dontRegex := regexp.MustCompile(`don't\(\)`)
	doRegex := regexp.MustCompile(`do\(\)`)

  result := 0
  mulEnabled := true
  
  // Find all tokens: {don't(), do() or mul(int, int)}
  // Use state machine to determine whether to count found mul operations.
  tokenRegex := regexp.MustCompile(`(?:(don't\(\))|(do\(\))|(mul\(\d+,\d+\)))`)
  tokens := tokenRegex.FindAllString(data, -1)
  for _, token := range tokens {
    switch {
    case dontRegex.MatchString(token):
      mulEnabled = false
    case doRegex.MatchString(token):
      mulEnabled = true
    case mulRegex.MatchString(token):
      if mulEnabled {
        match := mulRegex.FindStringSubmatch(token)
        num, err := calculate_mul(match)
        if err != nil {
          panic(err)
        }
        result += num
      }
    }
  }

  fmt.Println("Sum2: ", result)
}

func main() {
  // Load input data
  filePath := "input.txt"
  data, err := loadInput(filePath)
  if err != nil {
    fmt.Println("Error:", err)
    return
  }
  
  part_one(data)
  part_two(data)
  }
