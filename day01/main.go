package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

	for scanner.Scan() {
		line := scanner.Text()
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

		calibrationValue, err := strconv.Atoi(string([]rune{firstDigit, lastDigit}))
		if err != nil {
			log.Fatal(err)
		}

		part1 += calibrationValue
	}

	fmt.Println("Part 1:", part1)
}
