package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var symbols []Coordinates
	var possibleGears []Coordinates
	var partNumbers []PartNumber
	y := 0
	reNumber := regexp.MustCompile(`(\d+)`)
	reSymbol := regexp.MustCompile(`[^.\d]`)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		for _, symbolMatch := range reSymbol.FindAllStringIndex(line, -1) {
			coordinates := Coordinates{x: symbolMatch[0], y: y}
			symbols = append(symbols, coordinates)
			if line[symbolMatch[0]:symbolMatch[1]] == "*" {
				possibleGears = append(possibleGears, coordinates)
			}
		}

		for _, numberMatch := range reNumber.FindAllStringIndex(line, -1) {
			from := Coordinates{x: numberMatch[0], y: y}
			to := Coordinates{x: numberMatch[1] - 1, y: y}
			value, err := strconv.Atoi(line[numberMatch[0]:numberMatch[1]])
			if err != nil {
				log.Fatal(err)
			}

			partNumbers = append(partNumbers, PartNumber{from, to, value})
		}

		y++
	}

	part1 := 0
partsLoop:
	for _, partNumber := range partNumbers {
		for _, symbol := range symbols {
			if partNumber.isAdjacent(symbol) {
				part1 += partNumber.value
				continue partsLoop
			}
		}
	}

	part2 := 0
possibleGearsLoop:
	for _, possibleGear := range possibleGears {
		var adjacentNumbers []PartNumber

		for _, partNumber := range partNumbers {
			if partNumber.isAdjacent(possibleGear) {
				adjacentNumbers = append(adjacentNumbers, partNumber)
			}

			if len(adjacentNumbers) > 2 {
				continue possibleGearsLoop
			}
		}

		if len(adjacentNumbers) == 2 {
			part2 += adjacentNumbers[0].value * adjacentNumbers[1].value
		}
	}

	fmt.Println("Part 1", part1)
	fmt.Println("Part 2", part2)
}

type PartNumber struct {
	from, to Coordinates
	value    int
}

func (self PartNumber) isAdjacent(coordinates Coordinates) bool {
	for x := self.from.x; x <= self.to.x; x++ {
		if coordinates.isAdjacent(Coordinates{x, self.from.y}) {
			return true
		}
	}

	return false
}

type Coordinates struct {
	x, y int
}

func (self Coordinates) isAdjacent(other Coordinates) bool {
	return other.x >= self.x-1 && other.x <= self.x+1 && other.y >= self.y-1 && other.y <= self.y+1
}
