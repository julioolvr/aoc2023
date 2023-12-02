package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	part1 := 0
	part2 := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineWithReplacements := part2Replacer(line)

		part1 += findCalibrationValue(line)
		part2 += findCalibrationValue(lineWithReplacements)
	}

	fmt.Println("Part 1:", part1)
	fmt.Println("Part 2:", part2)
}

var numberNames = map[string]string{
	"one":   "1",
	"two":   "2",
	"three": "3",
	"four":  "4",
	"five":  "5",
	"six":   "6",
	"seven": "7",
	"eight": "8",
	"nine":  "9",
}

func part2Replacer(input string) string {
	var result strings.Builder

	for i, c := range input {
		found := false

		for name, digit := range numberNames {
			if strings.HasPrefix(input[i:], name) {
				result.WriteString(digit)
				found = true
				break
			}
		}

		if !found {
			result.WriteRune(c)
		}
	}

	return result.String()
}

func findCalibrationValue(line string) int {
	var firstDigit rune
	var lastDigit rune

	for _, c := range line {
		if unicode.IsDigit(c) {
			if firstDigit == 0 {
				firstDigit = c
			}

			lastDigit = c
		}
	}

	if firstDigit == 0 {
		firstDigit = '0'
	}

	if lastDigit == 0 {
		lastDigit = '0'
	}

	calibrationValue, err := strconv.Atoi(string([]rune{firstDigit, lastDigit}))
	if err != nil {
		log.Fatal(err)
	}

	return calibrationValue
}
