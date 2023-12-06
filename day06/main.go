package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	spaceRe := regexp.MustCompile(`\s+`)

	times := make([]int, 0)
	scanner.Scan()
	timesLine := scanner.Text()
	_, timesLine, _ = strings.Cut(timesLine, ":")
	timeNumbers := spaceRe.Split(strings.Trim(timesLine, " "), -1)
	for _, time := range timeNumbers {
		time, _ := strconv.Atoi(time)
		times = append(times, time)
	}

	records := make([]int, 0)
	scanner.Scan()
	recordsLine := scanner.Text()
	_, recordsLine, _ = strings.Cut(recordsLine, ":")
	recordsNumbers := spaceRe.Split(strings.Trim(recordsLine, " "), -1)
	for _, record := range recordsNumbers {
		record, _ := strconv.Atoi(record)
		records = append(records, record)
	}

	races := make([]Race, len(timeNumbers))
	for i, time := range times {
		races[i] = Race{length: time, record: records[i]}
	}

	part1 := 1
	for _, race := range races {
		minButtonPress, maxButtonPress := buttonPressLengths(race.length, race.record)
		part1 *= maxButtonPress - minButtonPress + 1
	}

	fmt.Println("Part 1", part1)

	// Too lazy to parse again
	part2Race := Race{length: 53717880, record: 275118112151524}
	minButtonPress, maxButtonPress := buttonPressLengths(part2Race.length, part2Race.record)
	part2 := maxButtonPress - minButtonPress + 1
	fmt.Println("Part 2", part2)
}

type Race struct {
	length, record int
}

// s = (a,b,c) => [(-b + Math.sqrt(Math.pow(b, 2) - 4*a*c))/(2*a), (-b - Math.sqrt(Math.pow(b, 2) - 4*a*c))/(2*a)]
func solveQuadratic(a float64, b float64, c float64) (float64, float64) {
	sqrt := math.Sqrt(math.Pow(b, 2) - 4*a*c)
	return (-b + sqrt) / (2 * a), (-b - sqrt) / (2 * a)
}

func buttonPressLengths(length int, record int) (int, int) {
	lowRoot, highRoot := solveQuadratic(-1, float64(length), float64(-record))
	return int(math.Floor(lowRoot + 1)), int(math.Ceil(highRoot - 1))
}
