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
	var partNumbers []PartNumber
	y := 0
	reNumber := regexp.MustCompile(`(\d+)`)
	reSymbol := regexp.MustCompile(`[^.\d]`)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		for _, symbolMatch := range reSymbol.FindAllStringIndex(line, -1) {
			symbols = append(symbols, Coordinates{x: symbolMatch[0], y: y})
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

	fmt.Println("Part 1", part1)
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
