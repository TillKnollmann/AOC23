package main

import (
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const DAY = "17"

type Game struct {
	fields           [][]int
	limitX           int
	limitY           int
	minimalDistances [][]MinDistance
}

type MinDistance struct {
	mapping    map[string]int
	minArrival int
}

type Point struct {
	positionX int
	positionY int
}

type Direction int

const (
	North Direction = iota + 1
	East
	South
	West
)

type State struct {
	positionX       int
	positionY       int
	direction       Direction
	currentHeatLoss int
	currentStraight int
}

type PriorityQueue []*State

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].currentHeatLoss < pq[j].currentHeatLoss
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*State)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func stringify(direction Direction, currentStraight int) string {

	return fmt.Sprintf("%d,%d", direction, currentStraight)
}

func parseGame(input string) Game {

	var game Game

	lines := strings.Split(strings.ReplaceAll(input, "\r", ""), "\n")

	for y := 0; y < len(lines); y++ {

		var row []int
		var initializedDistances []MinDistance

		for x := 0; x < len(lines[y]); x++ {

			row = append(row, int(stringToNumber(string(lines[y][x]))))
			initializedDistances = append(initializedDistances, MinDistance{mapping: make(map[string]int), minArrival: math.MaxInt})
		}

		game.fields = append(game.fields, row)
		game.minimalDistances = append(game.minimalDistances, initializedDistances)
	}

	game.limitY = len(lines)
	game.limitX = len(lines[0])

	return game
}

func initializeCalculateDistances(game *Game, minimumStraight int, maximumStraight int) {

	game.minimalDistances[0][0] = MinDistance{mapping: make(map[string]int), minArrival: 0}

	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	heatLoss := 0
	for i := 1; i < minimumStraight; i++ {

		heatLoss += game.fields[0][0+i]
	}

	game.minimalDistances[0][0].mapping[stringify(East, minimumStraight)] = heatLoss

	heap.Push(&pq, &State{
		positionX:       minimumStraight,
		positionY:       0,
		direction:       East,
		currentHeatLoss: heatLoss,
		currentStraight: minimumStraight,
	})

	heatLoss = 0
	for i := 1; i < minimumStraight; i++ {

		heatLoss += game.fields[0+i][0]
	}

	game.minimalDistances[0][0].mapping[stringify(South, minimumStraight)] = heatLoss

	heap.Push(&pq, &State{
		positionX:       0,
		positionY:       minimumStraight,
		direction:       South,
		currentHeatLoss: heatLoss,
		currentStraight: minimumStraight,
	})

	for pq.Len() > 0 {

		state := heap.Pop(&pq).(*State)
		calculateDistances(game, Point{
			positionX: state.positionX,
			positionY: state.positionY,
		}, state.direction, state.currentHeatLoss, state.currentStraight, minimumStraight, maximumStraight, &pq)
	}
}

func cacheValueExitDirection(game *Game, point Point, exitDirection Direction, maxCurrentStraight int) int {

	minValue := math.MaxInt

	for i := 1; i <= maxCurrentStraight; i++ {

		cachedValue, ok := game.minimalDistances[point.positionY][point.positionX].mapping[stringify(exitDirection, i)]

		if ok {

			minValue = min(minValue, cachedValue)
		}
	}

	return minValue
}

func cacheNotHitAndUpdate(game *Game, point Point, exitDirection Direction, currentStraight int, currentHeatLoss int) bool {

	cachedValue := cacheValueExitDirection(game, point, exitDirection, currentStraight)

	if cachedValue <= currentHeatLoss {

		return false
	} else {

		game.minimalDistances[point.positionY][point.positionX].mapping[stringify(exitDirection, currentStraight)] = currentHeatLoss
		return true
	}
}

func calculateDistances(game *Game, point Point, entryDirection Direction, totalHeatLoss int, currentStraight int, minimumStraight int, maximalStraight int, pq *PriorityQueue) {

	if point.positionX < 0 || point.positionX >= game.limitX || point.positionY < 0 || point.positionY >= game.limitY {

		// out of range
		return
	}

	game.minimalDistances[point.positionY][point.positionX].minArrival = min(game.minimalDistances[point.positionY][point.positionX].minArrival, totalHeatLoss+game.fields[point.positionY][point.positionX])

	// calculate possible further steps
	// can we go south?
	if (entryDirection != North) && !(entryDirection == South && currentStraight == maximalStraight) {

		nextStraight := 0
		nextSteps := 0
		nextHeatloss := 0

		if entryDirection == South {

			nextStraight = currentStraight + 1
			nextSteps = 1
			nextHeatloss = totalHeatLoss + game.fields[point.positionY][point.positionX]

		} else {

			nextStraight = minimumStraight
			nextSteps = minimumStraight

			nextHeatloss = totalHeatLoss
			for i := 0; i < minimumStraight; i++ {

				if point.positionY+i < game.limitY {
					nextHeatloss += game.fields[point.positionY+i][point.positionX]
				}
			}
		}

		if point.positionY+nextSteps < game.limitY {

			if cacheNotHitAndUpdate(game, point, South, nextStraight, nextHeatloss) {

				heap.Push(pq, &State{
					positionX:       point.positionX,
					positionY:       point.positionY + nextSteps,
					direction:       South,
					currentHeatLoss: nextHeatloss,
					currentStraight: nextStraight,
				})
			}
		}
	}

	// can we go west?
	if (entryDirection != East) && !(entryDirection == West && currentStraight == maximalStraight) {

		nextStraight := 0
		nextSteps := 0
		nextHeatloss := 0

		if entryDirection == West {

			nextStraight = currentStraight + 1
			nextSteps = 1
			nextHeatloss = totalHeatLoss + game.fields[point.positionY][point.positionX]
		} else {

			nextStraight = minimumStraight
			nextSteps = minimumStraight

			nextHeatloss = totalHeatLoss
			for i := 0; i < minimumStraight; i++ {

				if point.positionX-i >= 0 {
					nextHeatloss += game.fields[point.positionY][point.positionX-i]
				}
			}
		}

		if point.positionX-nextSteps >= 0 {

			if cacheNotHitAndUpdate(game, point, West, nextStraight, nextHeatloss) {

				heap.Push(pq, &State{
					positionX:       point.positionX - nextSteps,
					positionY:       point.positionY,
					direction:       West,
					currentHeatLoss: nextHeatloss,
					currentStraight: nextStraight,
				})
			}
		}
	}

	// can we go north?
	if (entryDirection != South) && !(entryDirection == North && currentStraight == maximalStraight) {

		nextStraight := 0
		nextSteps := 0
		nextHeatloss := 0

		if entryDirection == North {

			nextStraight = currentStraight + 1
			nextSteps = 1
			nextHeatloss = totalHeatLoss + game.fields[point.positionY][point.positionX]
		} else {

			nextSteps = minimumStraight
			nextStraight = minimumStraight

			nextHeatloss = totalHeatLoss
			for i := 0; i < minimumStraight; i++ {

				if point.positionY-i >= 0 {
					nextHeatloss += game.fields[point.positionY-i][point.positionX]
				}
			}
		}

		if point.positionY-nextSteps >= 0 {

			if cacheNotHitAndUpdate(game, point, North, nextStraight, nextHeatloss) {

				heap.Push(pq, &State{
					positionX:       point.positionX,
					positionY:       point.positionY - nextSteps,
					direction:       North,
					currentHeatLoss: nextHeatloss,
					currentStraight: nextStraight,
				})
			}
		}
	}

	// can we go east?
	if (entryDirection != West) && !(entryDirection == East && currentStraight == maximalStraight) {

		nextStraight := 0
		nextSteps := 0
		nextHeatloss := 0

		if entryDirection == East {

			nextStraight = currentStraight + 1
			nextSteps = 1
			nextHeatloss = totalHeatLoss + game.fields[point.positionY][point.positionX]
		} else {

			nextStraight = minimumStraight
			nextSteps = minimumStraight

			nextHeatloss = totalHeatLoss
			for i := 0; i < minimumStraight; i++ {

				if point.positionX+i < game.limitX {
					nextHeatloss += game.fields[point.positionY][point.positionX+i]
				}
			}
		}

		if point.positionX+nextSteps < game.limitX {

			if cacheNotHitAndUpdate(game, point, East, nextStraight, nextHeatloss) {

				heap.Push(pq, &State{
					positionX:       point.positionX + nextSteps,
					positionY:       point.positionY,
					direction:       East,
					currentHeatLoss: nextHeatloss,
					currentStraight: nextStraight,
				})
			}
		}
	}
}

func getMinimalDistance(minDistance MinDistance) int {

	var minimum = math.MaxInt

	for _, value := range minDistance.mapping {

		minimum = min(minimum, value)
	}

	return minimum
}

func Part1(input string) string {

	content := GetContent(input)

	game := parseGame(content)
	initializeCalculateDistances(&game, 1, 3)

	return strconv.Itoa(game.minimalDistances[game.limitY-1][game.limitX-1].minArrival)
}

func Part2(input string) string {

	content := GetContent(input)

	game := parseGame(content)
	initializeCalculateDistances(&game, 4, 10)

	return strconv.Itoa(game.minimalDistances[game.limitY-1][game.limitX-1].minArrival)
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
