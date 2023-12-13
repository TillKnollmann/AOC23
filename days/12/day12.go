package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const DAY = "12"

type Game struct {
	path   string
	groups []int64
}

func (game Game) String() string {

	return fmt.Sprintf("(%s, %#v)", game.path, game.groups)
}

func serialize(game Game) string {

	// thanks to go, we need to stringify our game manually
	return fmt.Sprintf("%s", game)
}

type Cache struct {
	mapping map[string]int64
}

func parseGame(input string) Game {

	var game Game

	numberRe := regexp.MustCompile(`\d+`)

	split := strings.Split(strings.ReplaceAll(input, "\r", ""), " ")

	// append a "." for simpler recursion (we can always check the character after a sequence of "#")
	game.path = split[0] + "."

	for _, numberString := range numberRe.FindAllString(split[1], -1) {

		game.groups = append(game.groups, stringToNumber(numberString))
	}

	return game
}

func calculatePossibilities(game Game, cache *Cache) int64 {

	// initialize cache
	if cache == nil {

		newCache := Cache{mapping: make(map[string]int64)}
		cache = &newCache
	}

	// check the cache for an answer
	res, ok := cache.mapping[serialize(game)]

	if ok {

		return res
	}

	if len(game.groups) == 0 && !strings.Contains(game.path, "#") {

		// we have no groups, but there are also none required
		cache.mapping[serialize(game)] = 1
		return 1

	}

	if len(game.groups) == 0 && strings.Contains(game.path, "#") {

		// we ran out of groups, but the path requires them
		cache.mapping[serialize(game)] = 0
		return 0
	}

	if len(game.path) == 0 {

		// the path ended, but we still have groups
		cache.mapping[serialize(game)] = 0
		return 0
	}

	// we have groups and path remaining
	// process next character on path
	if string(game.path[0]) == "." {

		return treatDot(game, cache)
	}

	if string(game.path[0]) == "#" {

		return treatHashtag(game, cache)

	}

	// game.path[0] must be "?", assume it is "." or "#" and add the result
	possibilities := treatDot(game, cache) + treatHashtag(game, cache)
	cache.mapping[serialize(game)] = possibilities
	return possibilities

}

func treatDot(game Game, cache *Cache) int64 {

	// simply skip the dot in the path
	possibilities := calculatePossibilities(Game{path: game.path[1:], groups: game.groups}, cache)
	cache.mapping[serialize(game)] = possibilities
	return possibilities
}

func treatHashtag(game Game, cache *Cache) int64 {

	nextGroupSize := int(game.groups[0])

	// can the group fit in the next segment of the path of "#"s?
	if len(game.path) < nextGroupSize {

		// path is too short
		cache.mapping[serialize(game)] = 0
		return 0
	}

	pathSegment := game.path[0:nextGroupSize]
	// treat "?" as "#" temporarily
	pathSegment = strings.ReplaceAll(pathSegment, "?", "#")
	groupString := strings.Repeat("#", nextGroupSize)

	if pathSegment != groupString {

		// path does not contain sufficient "#" (or "?")
		cache.mapping[serialize(game)] = 0
		return 0
	}

	// check if group is too short
	pathContinueIndex := nextGroupSize

	if string(game.path[pathContinueIndex]) == "#" {

		// path requires longer group
		cache.mapping[serialize(game)] = 0
		return 0
	}

	// group fits perfectly, the path continues with "." or "?"
	// remove the group, cut the path and continue
	possibilities := calculatePossibilities(Game{path: game.path[pathContinueIndex+1:], groups: game.groups[1:]}, cache)
	cache.mapping[serialize(game)] = possibilities
	return possibilities
}

func blowUpGame(game *Game) {

	// first remove the artificial dot
	game.path = game.path[:len(game.path)-1]

	originalPath := game.path
	originalGroups := make([]int64, len(game.groups))
	copy(originalGroups, game.groups)

	for i := 0; i < 4; i++ {

		game.path = game.path + "?" + originalPath
		game.groups = append(game.groups, originalGroups...)
	}

	// add the artificial dot again
	game.path = game.path + "."
}

func Part1(input string) string {

	content := GetContent(input)

	sum := int64(0)

	for _, line := range strings.Split(strings.ReplaceAll(content, "\r", ""), "\n") {

		sum += calculatePossibilities(parseGame(line), nil)
	}

	return strconv.FormatInt(sum, 10)
}

func Part2(input string) string {

	content := GetContent(input)

	sum := int64(0)

	for _, line := range strings.Split(strings.ReplaceAll(content, "\r", ""), "\n") {

		game := parseGame(line)
		blowUpGame(&game)
		sum += calculatePossibilities(game, nil)
	}

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
