package memory

import (
	"day9/input"
	"slices"
)

// FillEmpty moves all files in memory to the left to fill empty spaces
func FillEmpty(memory []input.Status) {
	// Initialise to first empty memory, and last data
	ptr_empty := slices.Index(memory, input.EMPTY)
	ptr_data := len(memory) - 1

	for ptr_empty < ptr_data {
		memory[ptr_empty] = memory[ptr_data]
		memory[ptr_data] = input.EMPTY
		// Update pointers
		ptr_data--
		for {
			if memory[ptr_empty] == input.EMPTY || ptr_empty >= ptr_data {
				break
			}
			ptr_empty++
		}
	}
}
