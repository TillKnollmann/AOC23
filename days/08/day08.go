package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const DAY = "08"

type State struct {
	commandIndex int64
	position     string
	loopPosition int
	loopSequence []int64
	offset       int64
	steps        int64
}

func (b State) String() string {

	return fmt.Sprintf("State(commandIndex=%d, position=%s, loopPosition=%d, loopSequence=%#v, offset=%d, steps=%d)", b.commandIndex, b.position, b.loopPosition, b.loopSequence, b.offset, b.steps)
}

type Tuple struct {
	left  string
	right string
}

type Game struct {
	commands string
	tuples   map[string]Tuple
}

func parseGame(input string) Game {

	lines := strings.Split(input, "\n")

	var game Game
	game.commands = strings.Trim(lines[0], "\r")
	game.tuples = make(map[string]Tuple)

	for _, line := range lines[2:] {

		chunks := strings.Split(line, " = ")
		game.tuples[chunks[0]] = parseTuple(chunks[1])
	}

	return game
}

func parseTuple(input string) Tuple {

	input = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(input, "(", ""), ")", ""), " ", ""), "\r", "")

	var tuple Tuple

	tuple.left = strings.Split(input, ",")[0]
	tuple.right = strings.Split(input, ",")[1]

	return tuple
}

func transition(state State, game Game) State {

	if state.commandIndex >= int64(len(game.commands)) {
		state.commandIndex = 0
	}

	var newState State

	newState.commandIndex = state.commandIndex + 1

	destinationTuple := game.tuples[state.position]

	if rune(game.commands[state.commandIndex]) == 'L' {

		newState.position = destinationTuple.left
	} else {

		newState.position = destinationTuple.right
	}

	return newState
}

func getOffsetAndLoop(state State, game Game) (int64, []int64) {

	steps := int64(0)

	var finishPositions []string
	var finishCommandIndex []int64
	var finishSteps []int64

	for {

		if rune(state.position[2]) == 'Z' {

			if slices.Contains(finishPositions, state.position) {

				index := indexOfFirst(finishPositions, finishCommandIndex, state.position, state.commandIndex)

				if index != -1 {

					var loop []int64

					for finishStepIndex := index + 1; finishStepIndex < len(finishSteps); finishStepIndex++ {

						loop = append(loop, finishSteps[finishStepIndex]-finishSteps[finishStepIndex-1])
					}

					loop = append(loop, steps-finishSteps[len(finishSteps)-1])

					return finishSteps[index], loop
				}
			}

			finishPositions = append(finishPositions, state.position)
			finishCommandIndex = append(finishCommandIndex, state.commandIndex)
			finishSteps = append(finishSteps, steps)
		}
		state = transition(state, game)
		steps++
	}
}

func indexOfFirst(stringSlice []string, intSlice []int64, s string, i int64) int {

	for index, elem := range stringSlice {

		if elem == s && intSlice[index] == i {

			return index
		}
	}

	return -1
}

func areStatesAlignedAndSteps(states []State) (bool, int64) {

	steps := states[0].steps

	for _, state := range states[1:] {

		if state.steps != steps {

			return false, -1
		}
	}

	return true, steps
}

func getIndexOfSlowestState(states []State) int {

	slowestIndex := 0
	slowestSteps := states[0].steps

	for index, state := range states {

		if state.steps < slowestSteps {

			slowestIndex = index
		}
	}

	return slowestIndex
}

func Part1(input string) string {

	content := GetContent(input)

	game := parseGame(content)

	state := State{position: "AAA", commandIndex: int64(0)}

	steps := int64(0)

	for {

		if state.position == "ZZZ" {

			break
		}
		state = transition(state, game)
		steps++
	}

	return strconv.FormatInt(steps, 10)
}

func Part2(input string) string {

	content := GetContent(input)

	game := parseGame(content)

	// Get starting positions
	var startPositions []string
	var states []State

	for key, _ := range game.tuples {

		if rune(key[2]) == 'A' && !slices.Contains(startPositions, key) {

			states = append(states, State{position: key, commandIndex: 0})
			startPositions = append(startPositions, key)
		}
	}

	// get offsets and loop sizes
	var newStates []State

	for _, state := range states {

		state.offset, state.loopSequence = getOffsetAndLoop(state, game)
		state.loopPosition = 0
		state.steps = state.offset

		newStates = append(newStates, state)
	}

	states = newStates

	// until all states are aligned, increase the one with the fewest steps by one loop
	for {

		aligned, steps := areStatesAlignedAndSteps(states)

		if aligned {
			return strconv.FormatInt(steps, 10)
		}

		slowestIndex := getIndexOfSlowestState(states)

		states[slowestIndex].loopPosition += 1
		if states[slowestIndex].loopPosition >= len(states[slowestIndex].loopSequence) {
			states[slowestIndex].loopPosition = 0
		}

		states[slowestIndex].steps += states[slowestIndex].loopSequence[states[slowestIndex].loopPosition]
	}
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
