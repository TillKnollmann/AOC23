package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const DAY = "06"

type Race struct {
	time     int64
	distance int64
}

func parse(input string) []Race {

	var races []Race

	numberRe := regexp.MustCompile(`\d+`)
	timesString := numberRe.FindAllString(strings.Split(input, "\n")[0], -1)
	distancesString := numberRe.FindAllString(strings.Split(input, "\n")[1], -1)

	for index := range timesString {

		races = append(races, Race{time: stringToNumber(timesString[index]), distance: stringToNumber(distancesString[index])})
	}

	return races
}

func stringToNumber(s string) int64 {

	number, err := strconv.ParseInt(s, 10, 64)

	if err != nil {

		panic(err)
	}

	return number
}

func getWinningPossibilities(race Race) int64 {

	minimumDriveTime := (float64(1) / float64(2)) * (math.Sqrt(math.Pow(float64(race.time), 2)-4*float64(race.distance)) + float64(race.time))

	numberOfPossibilities := (math.Ceil(minimumDriveTime-float64(1)) - math.Floor(float64(race.time)/float64(2))) * float64(2)

	// If t is even, we calculated one possibility too few, where pressing time = minimumDriveTime
	if race.time%2 == 0 {

		numberOfPossibilities += 1
	}

	return int64(numberOfPossibilities)
}

func parsePart2(input string) Race {

	var race Race

	numberRe := regexp.MustCompile(`\d+`)
	timeString := numberRe.FindString(strings.ReplaceAll(strings.Split(input, "\n")[0], " ", ""))
	distanceString := numberRe.FindString(strings.ReplaceAll(strings.Split(input, "\n")[1], " ", ""))

	race.time = stringToNumber(timeString)
	race.distance = stringToNumber(distanceString)

	return race
}

func Part1(input string) string {

	content := GetContent(input)

	races := parse(content)

	result := int64(1)

	for _, race := range races {

		possibilities := getWinningPossibilities(race)
		result *= possibilities
	}

	return fmt.Sprintf("%d", result)
}

func Part2(input string) string {

	content := GetContent(input)

	race := parsePart2(content)

	possibilities := getWinningPossibilities(race)

	return fmt.Sprintf("%d", possibilities)
}

func GetContent(filepath string) string {

	content, err := os.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func main() {

	fmt.Println(fmt.Sprintf("Part 1: %s", Part1(fmt.Sprintf("input/%s/in.txt", DAY))))
	fmt.Println(fmt.Sprintf("Part 2: %s", Part2(fmt.Sprintf("input/%s/in.txt", DAY))))
}
