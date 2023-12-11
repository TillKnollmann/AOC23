package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const DAY = "11"

type Galaxy struct {
	x int
	y int
}

type EmptySpace struct {
	size int64
}

type Universe struct {
	galaxies         []Galaxy
	horizontalSpaces map[int]EmptySpace
	verticalSpaces   map[int]EmptySpace
}

func parseUniverse(input string) Universe {

	var universe Universe

	lines := strings.Split(strings.ReplaceAll(input, "\r", ""), "\n")

	// contain rows/columns which are empty
	var horizontalSpaces map[int]bool = make(map[int]bool)
	var verticalSpaces map[int]bool = make(map[int]bool)

	for i := 0; i < len(lines); i++ {
		horizontalSpaces[i] = true
	}
	for i := 0; i < len(lines[0]); i++ {
		verticalSpaces[i] = true
	}

	for y := 0; y < len(lines); y++ {
		for x := 0; x < len(lines[y]); x++ {
			if string(lines[y][x]) == "#" {
				universe.galaxies = append(universe.galaxies, Galaxy{x: x, y: y})
				delete(horizontalSpaces, y)
				delete(verticalSpaces, x)
			}
		}
	}

	universe.horizontalSpaces = make(map[int]EmptySpace)
	universe.verticalSpaces = make(map[int]EmptySpace)

	for key, _ := range horizontalSpaces {
		universe.horizontalSpaces[key] = EmptySpace{size: 1}
	}
	for key, _ := range verticalSpaces {
		universe.verticalSpaces[key] = EmptySpace{size: 1}
	}

	return universe
}

func expandSpaces(universe *Universe, expansionFactor int64) {

	for key, _ := range universe.verticalSpaces {
		universe.verticalSpaces[key] = EmptySpace{size: expansionFactor * universe.verticalSpaces[key].size}
	}
	for key, _ := range universe.horizontalSpaces {
		universe.horizontalSpaces[key] = EmptySpace{size: expansionFactor * universe.horizontalSpaces[key].size}
	}
}

func getDistance(first Galaxy, second Galaxy, universe Universe) int64 {

	minX := min(first.x, second.x)
	maxX := max(first.x, second.x)
	minY := min(first.y, second.y)
	maxY := max(first.y, second.y)

	distance := int64(0)

	if minX != maxX {
		distance += int64(maxX - minX)
		for i := minX + 1; i < maxX; i++ {
			emptySpace, ok := universe.verticalSpaces[i]
			if ok {
				distance += emptySpace.size - 1
			}
		}
	}
	if minY != maxY {
		distance += int64(maxY - minY)
		for i := minY + 1; i < maxY; i++ {
			emptySpace, ok := universe.horizontalSpaces[i]
			if ok {
				distance += emptySpace.size - 1
			}
		}
	}

	return distance
}

func Part1(input string) string {

	content := GetContent(input)

	universe := parseUniverse(content)

	expandSpaces(&universe, 2)

	sum := int64(0)

	for _, first := range universe.galaxies {
		for _, second := range universe.galaxies {
			if !(first.x == second.x && first.y == second.y) {
				sum += getDistance(first, second, universe)
			}
		}
	}

	sum = sum / 2

	return strconv.FormatInt(sum, 10)
}

func Part2(input string) string {

	content := GetContent(input)
	universe := parseUniverse(content)

	expandSpaces(&universe, 1000000)

	sum := int64(0)

	for _, first := range universe.galaxies {
		for _, second := range universe.galaxies {
			if !(first.x == second.x && first.y == second.y) {
				sum += getDistance(first, second, universe)
			}
		}
	}

	sum = sum / 2

	return strconv.FormatInt(sum, 10)
}

func GetContent(filepath string) string {

	content, err := os.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func stringToNumber(s string) int64 {

	number, err := strconv.ParseInt(s, 10, 64)

	if err != nil {

		panic(err)
	}

	return number
}

func main() {

	fmt.Println(fmt.Sprintf("Part 1: %s", Part1(fmt.Sprintf("input/%s/in.txt", DAY))))
	fmt.Println(fmt.Sprintf("Part 2: %s", Part2(fmt.Sprintf("input/%s/in.txt", DAY))))
}
