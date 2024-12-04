package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

// Crossword type, representing its grid
type crossword [][]byte

// Reverses a string
func reverseString(str string) string {
	runes := []rune(str)
	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Provides a string representation of the crossword
func (c crossword) String() string {
	var builder strings.Builder

	builder.WriteString("Crossword: \n")
	for i, row := range c {
		builder.WriteString(string(row))

		if i < len(c)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

// Reads a crossword from a file
func read_file(filename string) (crossword, error) {
	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read line by line
	var res crossword
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		byte_text := []byte(scanner.Text())
		res = append(res, byte_text)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Could not read file: %w", err)
	}
	return res, nil
}

// Find a string pattern in a crossword
func findInCrossword(c crossword, text string) int {
	target := []byte(text)
	targetRev := []byte(reverseString(text))
	count := 0
	rows, cols := len(c), len(c[0])

	// Helper function to check if a slice matches either normal or reversed target
	checkMatch := func(slice []byte) bool {
		return bytes.Equal(slice, target) || bytes.Equal(slice, targetRev)
	}

	// Traverse array by rows
	for row := 0; row < rows; row++ {
		for col := 0; col <= cols-len(text); col++ {
			if slice := c[row][col : col+len(text)]; checkMatch(slice) {
				count++
			}
		}
	}

	// Traverse array by cols
	for col := 0; col < cols; col++ {
		for row := 0; row < rows-len(text)+1; row++ {
			// Create a slice for each column
			slice := make([]byte, len(text))
			for i := 0; i < len(text); i++ {
				slice[i] = c[row+i][col]
			}
			if bytes.Equal(slice, target) || bytes.Equal(slice, targetRev) {
				count++
			}
		}
	}

	// Helper function for traversing diagonals
	checkDiagonal := func(startRow, startCol, deltaRow, deltaCol int) {
		diagonal := []byte{}

		for i, j := startRow, startCol; i >= 0 && i < rows && j >= 0 && j < cols; i, j = i+deltaRow, j+deltaCol {
			diagonal = append(diagonal, c[i][j])
		}

		for i := 0; i <= len(diagonal)-len(text); i++ {
			if slice := diagonal[i : i+len(text)]; checkMatch(slice) {
				count++
			}
		}
	}

	// Traverse diagonals (top-left -> bottom-right)
	for col := 0; col <= cols-len(text); col++ {
		checkDiagonal(0, col, 1, 1)
	}

	for row := 1; row <= rows-len(text); row++ {
		checkDiagonal(row, 0, 1, 1)
	}

	// Traverse diagonals (top-right -> bottom-left)
	for col := cols - 1; col > len(text); col-- {
		checkDiagonal(0, col, 1, -1)
	}

	for row := 1; row <= rows-len(text); row++ {
		checkDiagonal(row, cols-1, 1, -1)
	}

	return count
}

// Given a crossword, check if it matches a pattern
/*
M.S
.A.
M.S
Where each MAS on a diagonal, can be in either in normal or reversed direction.
*/
func (c crossword) match() bool {
	// A has to be at the cetner, always
	if c[1][1] != 'A' {
		return false
	}
	// Check for the four possible corner configurations
	return (c[0][0] == 'S' && c[0][2] == 'S' && c[2][0] == 'M' && c[2][2] == 'M') ||
		(c[0][0] == 'M' && c[0][2] == 'M' && c[2][0] == 'S' && c[2][2] == 'S') ||
		(c[0][0] == 'S' && c[0][2] == 'M' && c[2][0] == 'S' && c[2][2] == 'M') ||
		(c[0][0] == 'M' && c[0][2] == 'S' && c[2][0] == 'M' && c[2][2] == 'S')
}

// Find X-MAS pattern in a crossworrd
func findXMasInCrossword(c crossword) int {
	count := 0
	rows, cols := len(c), len(c[0])

	w, h := 3, 3
	// Pre-initialise a 3x3 crossword pattern
	var pattern crossword
	for i := 0; i < 3; i++ {
		pattern = append(pattern, make([]byte, 3)) // Create a 3-byte slice for each row
	}

	for row := 0; row <= rows-h; row++ {
		for col := 0; col <= cols-w; col++ {
			// Fill with values from crossword
			for i := 0; i < w; i++ {
				copy(pattern[i], c[row+i][col:col+3])
			}

			// Match pattern
			if pattern.match() {
				count++
			}
		}
	}
	return count
}
func main() {
	c, err := read_file("input.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Part I")
	fmt.Println("Count: ", findInCrossword(c, "XMAS"))

	fmt.Println("Part II")
	fmt.Println("Count: ", findXMasInCrossword(c))
}
