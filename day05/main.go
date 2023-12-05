package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
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

	scanner.Scan()
	seedsLine := scanner.Text()
	_, seedsDefinition, _ := strings.Cut(seedsLine, ": ")
	seeds := regexp.MustCompile(`\s+`).Split(strings.Trim(seedsDefinition, " "), -1)

	// Skip empty line after seeds
	scanner.Scan()

	resourcesMap := make(map[string]ResourceMap)
	headerRe := regexp.MustCompile(`(.+)-to-(.+) map:`)

	for scanner.Scan() {
		headerDefinition := scanner.Text()
		match := headerRe.FindStringSubmatch(headerDefinition)
		sourceCategory := match[1]
		destinationCategory := match[2]

		ranges := make([]Range, 0)

		for scanner.Scan() {
			rangeLine := scanner.Text()
			if rangeLine == "" {
				break
			}

			rangeDefinition := regexp.MustCompile(`\s+`).Split(strings.Trim(rangeLine, " "), -1)
			fromDestination, _ := strconv.Atoi(rangeDefinition[0])
			fromSource, _ := strconv.Atoi(rangeDefinition[1])
			length, _ := strconv.Atoi(rangeDefinition[2])
			ranges = append(ranges, Range{fromSource, fromDestination, length})
		}

		resourcesMap[sourceCategory] = ResourceMap{ranges, destinationCategory}
	}

	seedLocations := make([]int, 0)

	for _, seed := range seeds {
		seedNumber, _ := strconv.Atoi(seed)
		soil := resourcesMap["seed"].valueFor(seedNumber)
		fertilizer := resourcesMap["soil"].valueFor(soil)
		water := resourcesMap["fertilizer"].valueFor(fertilizer)
		light := resourcesMap["water"].valueFor(water)
		temperature := resourcesMap["light"].valueFor(light)
		humidity := resourcesMap["temperature"].valueFor(temperature)
		location := resourcesMap["humidity"].valueFor(humidity)

		seedLocations = append(seedLocations, location)
	}

	part1 := slices.Min(seedLocations)
	fmt.Println("Part 1", part1)
}

type ResourceMap struct {
	ranges              []Range
	destinationCategory string
}

func (resourceMap ResourceMap) valueFor(value int) int {
	for _, mapRange := range resourceMap.ranges {
		destinationValue := mapRange.mapValue(value)
		if destinationValue != -1 {
			return destinationValue
		}
	}

	return value
}

type Range struct {
	fromSource, fromDestination, length int
}

func (mapRange Range) mapValue(value int) int {
	if value < mapRange.fromSource || value >= mapRange.fromSource+mapRange.length {
		return -1
	}

	return mapRange.fromDestination + value - mapRange.fromSource
}
