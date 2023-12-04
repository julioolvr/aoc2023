package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
)

func main() {
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	part1 := 0.0

	for scanner.Scan() {
		line := scanner.Text()
		_, game, _ := strings.Cut(line, ": ")
		winningNumbers, cardNumbers, _ := strings.Cut(game, " | ")

		winningNumbersSet := mapset.NewSet(regexp.MustCompile(`\s+`).Split(strings.Trim(winningNumbers, " "), -1)...)
		cardNumbersSet := mapset.NewSet(regexp.MustCompile(`\s+`).Split(strings.Trim(cardNumbers, " "), -1)...)

		matchingNumbers := winningNumbersSet.Intersect(cardNumbersSet).Cardinality()

		if matchingNumbers > 0 {
			part1 += math.Pow(2, float64(matchingNumbers)-1)
		}
	}

	fmt.Println("Part 1", part1)
}
