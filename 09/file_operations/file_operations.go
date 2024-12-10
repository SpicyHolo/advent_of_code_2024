package file_operations

import (
	"day9/input"
	"slices"
	"sort"
)

// input.Status type, for memory status

const (
	EMPTY input.Status = -1
)

// set type for unique file IDs
type set map[input.Status]struct{}

// Add an element to the set
func (s set) add(element input.Status) {
	s[element] = struct{}{}
}

// remove an element from the set
func (s set) remove(element input.Status) {
	delete(s, element)
}

// ToArray returns a sorted array of the set elements
func (s set) toArray(sorted bool) []input.Status {
	var keys []int
	for key := range s {
		keys = append(keys, int(key))
	}
	if sorted {
		sort.Ints(keys)
	}
	var res []input.Status
	for key := range keys {
		res = append(res, input.Status(key))
	}
	return res
}

// GetFileIDs returns unique file IDs in sorted order
func getFileIDs(memory []input.Status, sorted bool) []input.Status {
	setIDs := make(set)

	for _, id := range memory {
		setIDs.add(id)
	}
	setIDs.remove(EMPTY)
	return setIDs.toArray(sorted)
}

// FindFile returns the start index and length of a file in memory
func findFile(memory []input.Status, fileID input.Status) (int, int) {
	startIdx := slices.Index(memory, fileID)
	length := 1
	for i := startIdx + 1; i < len(memory); i++ {
		if memory[i] == fileID {
			length++
		} else {
			break
		}
	}
	return startIdx, length
}

// FindNextFree returns the start index and length of the next free space
func findNextFree(memory []input.Status, startIdx int) (int, int) {
	for i := startIdx; i < len(memory); i++ {
		if memory[i] == EMPTY {
			// Found the start of an empty space
			freeStart := i
			freeLength := 0
			for j := i; j < len(memory); j++ {
				if memory[j] == EMPTY {
					freeLength++
				} else {
					break
				}
			}
			return freeStart, freeLength
		}
	}
	return -1, 0 // No empty space found
}

// MoveFile moves a file in memory from one position to another
func moveFile(memory []input.Status, fromIdx, toIdx, length int) {
	copy(memory[toIdx:], memory[fromIdx:fromIdx+length])
}

// ClearFileSpace clears the space occupied by a file in memory
func clearFileSpace(memory []input.Status, startIdx, length int) {
	for i := startIdx; i < startIdx+length; i++ {
		memory[i] = EMPTY
	}
}

// FillEmpty2 attempts to compact files by moving them to available spaces
func FillEmpty2(memory []input.Status) {
	fileIDs := getFileIDs(memory, true) // Get unique file IDs in descending order

	// Loop over files starting from the highest ID
	for i := len(fileIDs) - 1; i >= 0; i-- {
		fileID := fileIDs[i]

		// Find the start and length of the current file
		startIdx, fileLength := findFile(memory, fileID)

		// Try to move the file to an available free space
		moved := false
		freeStart := 0
		freeLength := 0
		for !moved {
			// Find the next available free space starting after the last move
			freeStart, freeLength = findNextFree(memory, freeStart)
			// If there's enough space to fit the file, move it
			if fileLength <= freeLength && freeStart != -1 {
				if freeStart < startIdx {
					moveFile(memory, startIdx, freeStart, fileLength)
					clearFileSpace(memory, startIdx, fileLength)
					moved = true
				} else {
					break
				}
			}

			// If we couldn't find a suitable free space, break out of the loop
			if freeStart == -1 {
				break
			} else if freeLength < fileLength {
				freeStart++
			}
		}
	}
}

// CheckSum calculates the checksum for the memory state
func CheckSum(memory []input.Status) (check_sum int) {
	for i, data := range memory {
		if data == EMPTY {
			continue
		}
		check_sum += i * int(data)
	}
	return
}
