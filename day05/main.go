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
	seedNumbers := make([]int, len(seeds))
	for i, seed := range seeds {
		seedNumber, _ := strconv.Atoi(seed)
		seedNumbers[i] = seedNumber
	}

	seedRanges := make([]NumbersRange, len(seedNumbers)/2)
	for i := 0; i < len(seedNumbers); i += 2 {
		seedRanges[i/2] = buildRange(seedNumbers[i], seedNumbers[i]+seedNumbers[i+1]-1)
	}

	// Skip empty line after seeds
	scanner.Scan()

	// A map where resources are the key (seed, soil, fertilizer) and the values
	// are maps to know how to convert from that resource to the next one in the
	// chain.
	almanac := make(map[string]ResourceMap)
	headerRe := regexp.MustCompile(`(.+)-to-(.+) map:`)

	for scanner.Scan() {
		headerDefinition := scanner.Text()
		match := headerRe.FindStringSubmatch(headerDefinition)
		sourceCategory := match[1]
		destinationCategory := match[2]

		ranges := make([]NumbersRangeMapping, 0)

		for scanner.Scan() {
			rangeLine := scanner.Text()
			if rangeLine == "" {
				break
			}

			rangeDefinition := regexp.MustCompile(`\s+`).Split(strings.Trim(rangeLine, " "), -1)
			fromDestination, _ := strconv.Atoi(rangeDefinition[0])
			fromSource, _ := strconv.Atoi(rangeDefinition[1])
			length, _ := strconv.Atoi(rangeDefinition[2])
			sourceRange := buildRange(fromSource, fromSource+length-1)
			destinationRange := buildRange(fromDestination, fromDestination+length-1)
			ranges = append(ranges, NumbersRangeMapping{sourceRange: sourceRange, destinationRange: destinationRange})
		}

		almanac[sourceCategory] = ResourceMap{ranges, destinationCategory}
	}

	seedLocations := make([]int, len(seedNumbers))
	path := []string{"seed", "soil", "fertilizer", "water", "light", "temperature", "humidity"}

	for i, seedNumber := range seedNumbers {
		result := seedNumber

		for _, resourceName := range path {
			result = almanac[resourceName].valueFor(result)
		}

		seedLocations[i] = result
	}

	part1 := slices.Min(seedLocations)
	fmt.Println("Part 1", part1)

	// The ranges in which a given resource is in. Initially the ranges of seeds, but as it gets transformed
	// based on the Almanac it changes what resource it contains.
	resourceRanges := seedRanges

	for _, sourceResourceName := range path {
		// This array will contain ranges that were originally in one position, but based on the almanac were
		// transformed into a different position
		mappedDestinationRanges := make([]NumbersRange, 0)

		for _, destinationResourceRange := range almanac[sourceResourceName].ranges {
			destinationRanges := make([]NumbersRange, 0)

			for _, sourceResourceRange := range resourceRanges {
				beforeRange, rangeOverlap, afterRange := destinationResourceRange.splitAndMap(sourceResourceRange)

				if !beforeRange.isEmpty() {
					destinationRanges = append(destinationRanges, beforeRange)
				}

				if !rangeOverlap.isEmpty() {
					mappedDestinationRanges = append(mappedDestinationRanges, rangeOverlap)
				}

				if !afterRange.isEmpty() {
					destinationRanges = append(destinationRanges, afterRange)
				}
			}

			resourceRanges = destinationRanges
		}

		resourceRanges = append(resourceRanges, mappedDestinationRanges...)
	}

	locationStarts := make([]int, len(resourceRanges))
	for i, resourceRange := range resourceRanges {
		locationStarts[i] = resourceRange.from
	}

	fmt.Println("Part 2", slices.Min(locationStarts))
}

type ResourceMap struct {
	ranges              []NumbersRangeMapping
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

// Represents a range of numbers from `from` to `to`, both inclusive
type NumbersRange struct {
	from, to int
}

func buildRange(from int, to int) NumbersRange {
	return NumbersRange{from, to}
}

func emptyRange() NumbersRange {
	// This is where I'd use a `None` value... IF I HAD ONE
	return buildRange(0, -1)
}

func (numbersRange NumbersRange) isEmpty() bool {
	return numbersRange.to == -1
}

func (numbersRange NumbersRange) moveTo(startingPosition int) NumbersRange {
	shift := numbersRange.from - startingPosition
	numbersRange.from -= shift
	numbersRange.to -= shift
	return numbersRange
}

func (numbersRange NumbersRange) split(otherRange NumbersRange) (NumbersRange, NumbersRange, NumbersRange) {
	var beforeRange, rangeOverlap, afterRange NumbersRange

	if otherRange.to > numbersRange.to {
		afterRange = buildRange(max(numbersRange.to+1, otherRange.from), otherRange.to)
	} else {
		afterRange = emptyRange()
	}

	if otherRange.from < numbersRange.from {
		beforeRange = buildRange(otherRange.from, min(numbersRange.from-1, otherRange.to))
	} else {
		beforeRange = emptyRange()
	}

	if otherRange.from <= numbersRange.to && otherRange.to >= numbersRange.from {
		from := max(otherRange.from, numbersRange.from)
		to := min(otherRange.to, numbersRange.to)
		rangeOverlap = buildRange(from, to)
	} else {
		rangeOverlap = emptyRange()
	}

	return beforeRange, rangeOverlap, afterRange
}

// A mapping from a source range to a destination range.
// For example, for a source range 10-15 and a destination range 20-25,
// it means that number 10 maps to 20, 11 to 21, 12 to 22 and so on.
type NumbersRangeMapping struct {
	sourceRange, destinationRange NumbersRange
}

// It takes a range and splits it into up to three parts. The first one is all the values in the given range
// that don't overlap with the current mapping because they come before it. Those values are returned as they were.
// The third one is the same, but for values that come after the mapping.
// The one in the middle is where they overlap. If that overlapping range is not empty, the range is converted based
// on the mapping to the destination range.
// For example, for a source range 10-15 and a destination range 20-25, if the range 8-17 is given, three ranges are returned:
// 8-9, for the first range
// 20-25, for the second (mapped) range
// 16-17, for the third range
func (numbersRangeMapping NumbersRangeMapping) splitAndMap(otherRange NumbersRange) (NumbersRange, NumbersRange, NumbersRange) {
	beforeRange, rangeOverlap, afterRange := numbersRangeMapping.sourceRange.split(otherRange)

	if !rangeOverlap.isEmpty() {
		rangeOverlap = rangeOverlap.moveTo(numbersRangeMapping.destinationRange.from + rangeOverlap.from - numbersRangeMapping.sourceRange.from)
	}

	return beforeRange, rangeOverlap, afterRange
}

// Given a value, it transforms it based on the current mapping. If the value does not match the mapping, -1 is returned instead.
func (numbersRangeMapping NumbersRangeMapping) mapValue(value int) int {
	if value < numbersRangeMapping.sourceRange.from || value >= numbersRangeMapping.sourceRange.to {
		return -1
	}

	return numbersRangeMapping.destinationRange.from + value - numbersRangeMapping.sourceRange.from
}
