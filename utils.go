package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func ConvertFileToIntLists(filename string) ([][]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result [][]int
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		var initialInts []int
		for _, part := range parts {
			value, err := strconv.Atoi(part)
			if err != nil {
				return nil, err
			}
			initialInts = append(initialInts, value)
		}
		result = append(result, initialInts)
	}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) == 3 {
			index, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, err
			}

			var messageType int
			switch parts[1] {
			case "SEND":
				messageType = 0
			case "RECEIVE":
				messageType = 1
			case "SNAPSHOT":
				messageType = 3
			case "WAIT":
				messageType = 2
			default:
				continue
			}

			value, _ := strconv.Atoi(parts[2])

			// Agregar la lista de enteros al resultado
			result = append(result, []int{index, messageType, value})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
