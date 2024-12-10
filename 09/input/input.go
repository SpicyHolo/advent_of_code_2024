package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// Status type, for memory status
type Status int

const (
	EMPTY Status = -1
)

func (s Status) String() string {
	if s == EMPTY {
		return "."
	}
	return strconv.Itoa(int(s))
}

// Parsing input
func ParseInput(filename string) ([]Status, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error could not read file: %w", err)
	}

	// Read data
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	data := scanner.Text()
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Parse data
	var id Status
	isFile := true
	var res []Status

	for _, char := range data {
		num := int(char - '0')
		if isFile {
			for i := 0; i < num; i++ {
				res = append(res, id)
			}
			id++
			isFile = false
		} else {
			for i := 0; i < num; i++ {
				res = append(res, EMPTY)
			}
			isFile = true
		}
	}
	return res, nil
}
