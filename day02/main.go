package main

import (
	"bufio"
	"fmt"
	"log"
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

	re := regexp.MustCompile(`Game \d+: (.+)`)
	gameId := 0
	part1 := 0
	part2 := 0

gameLoop:
	for scanner.Scan() {
		gameId++
		line := scanner.Text()
		parsedLine := re.FindStringSubmatch(line)
		game := parseGame(parsedLine[1])

		part2 += game.minimumPower()

		for _, round := range game.rounds {
			if !round.isPossible() {
				continue gameLoop
			}
		}

		part1 += gameId
	}

	fmt.Println("Part 1", part1)
	fmt.Println("Part 2", part2)
}

type Game struct {
	rounds []Round
}

func (game Game) minimumPower() int {
	red := 0
	green := 0
	blue := 0

	for _, round := range game.rounds {
		if round.red > red {
			red = round.red
		}
		if round.green > green {
			green = round.green
		}
		if round.blue > blue {
			blue = round.blue
		}
	}

	return red * green * blue
}

type Round struct {
	red, green, blue int
}

func (round Round) isPossible() bool {
	return round.red <= 12 && round.green <= 13 && round.blue <= 14
}

func parseGame(game string) Game {
	var rounds []Round

	for _, roundDefinition := range strings.Split(game, ";") {
		rounds = append(rounds, parseRound(roundDefinition))
	}

	return Game{rounds}
}

func parseRound(round string) Round {
	redRe := regexp.MustCompile(`(\d+) red`)
	greenRe := regexp.MustCompile(`(\d+) green`)
	blueRe := regexp.MustCompile(`(\d+) blue`)

	red := parseColor(*redRe, round)
	green := parseColor(*greenRe, round)
	blue := parseColor(*blueRe, round)

	return Round{red, green, blue}
}

func parseColor(re regexp.Regexp, roundDefinition string) int {
	result := re.FindStringSubmatch(roundDefinition)

	if len(result) > 0 {
		value, err := strconv.Atoi(result[1])
		if err != nil {
			log.Fatal(err)
		}

		return value
	} else {
		return 0
	}
}
