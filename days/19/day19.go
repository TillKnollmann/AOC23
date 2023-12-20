package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const DAY = "19"

type Range struct {
	minimum int
	maximum int
}

type RangeSet struct {
	ranges []RangeElement
}

type RangeElement struct {
	items map[string]Range
}

type Game struct {
	workflows map[string]Workflow
	elements  []Element
}

type Rule struct {
	property    string
	compare     string
	compareTo   int
	destination string
	bypassCheck bool
}

type Workflow struct {
	label string
	rules []Rule
}

type Element struct {
	items map[string]int
}

func parseGame(input string) Game {

	var game Game

	game.workflows = make(map[string]Workflow)

	workflows := strings.Split(strings.ReplaceAll(input, "\r", ""), "\n\n")[0]

	elements := strings.Split(strings.ReplaceAll(input, "\r", ""), "\n\n")[1]

	for _, workflow := range strings.Split(workflows, "\n") {

		flow := parseWorkflow(workflow)
		game.workflows[flow.label] = flow
	}

	for _, element := range strings.Split(elements, "\n") {

		game.elements = append(game.elements, parseElement(element))
	}

	return game
}

func parseWorkflow(input string) Workflow {

	var workflow Workflow

	trimmed := strings.ReplaceAll(input, "\r", "")

	workflow.label = strings.Split(trimmed, "{")[0]

	instructions := strings.Split(strings.ReplaceAll(strings.Split(trimmed, "{")[1], "}", ""), ",")

	for _, instruction := range instructions {

		workflow.rules = append(workflow.rules, parseRule(instruction))
	}

	return workflow
}

func parseRule(instruction string) Rule {

	var rule Rule

	if strings.Contains(instruction, ":") {

		rule.bypassCheck = false
	} else {

		rule.bypassCheck = true
	}

	if rule.bypassCheck {

		rule.destination = instruction
	} else {

		chunks := strings.Split(instruction, ":")

		rule.destination = chunks[1]

		rule.property = string(chunks[0][0])
		rule.compare = string(chunks[0][1])
		rule.compareTo = int(stringToNumber(string(chunks[0][2:])))
	}

	return rule
}

func parseElement(input string) Element {

	var element Element

	element.items = make(map[string]int)

	itemStrings := strings.Split(strings.ReplaceAll(strings.ReplaceAll(input, "{", ""), "}", ""), ",")

	for _, itemString := range itemStrings {

		chunks := strings.Split(itemString, "=")
		element.items[chunks[0]] = int(stringToNumber(chunks[1]))
	}

	return element
}

func isElementAccepted(element Element, workflows map[string]Workflow) bool {

	currentLabel := "in"
	accepted := false
	rejected := false

	for !accepted && !rejected {

		nextWorkflow, _ := workflows[currentLabel]
		currentLabel = transition(element, nextWorkflow)

		if currentLabel == "R" {

			rejected = true
		} else if currentLabel == "A" {

			accepted = true
		}
	}

	return accepted
}

func transition(element Element, workflow Workflow) string {

	nextLabel := ""

	for _, rule := range workflow.rules {

		if ruleAccepts(element, rule) {

			nextLabel = rule.destination
			break
		}
	}

	return nextLabel
}

func ruleAccepts(element Element, rule Rule) bool {

	if rule.bypassCheck {

		return true
	}

	item, _ := element.items[rule.property]

	switch rule.compare {
	case "<":
		return item < rule.compareTo
	case ">":
		return item > rule.compareTo
	}

	return false
}

func getElementScore(element Element) int {

	sum := 0

	for _, value := range element.items {

		sum += value
	}

	return sum
}

func calculateAcceptedRanges(currentRange RangeElement, currentWorkflow Workflow, game Game, acceptedRanges *RangeSet) {

	rangeForNextRule := copyRangeElement(currentRange)

	for _, rule := range currentWorkflow.rules {

		if rule.bypassCheck {

			if rule.destination == "A" {

				acceptedRanges.ranges = append(acceptedRanges.ranges, rangeForNextRule)

			} else if rule.destination != "R" {

				calculateAcceptedRanges(rangeForNextRule, game.workflows[rule.destination], game, acceptedRanges)

			}

			break

		} else {

			switch rule.compare {
			case "<":
				acceptedRange := copyRangeElement(rangeForNextRule)
				deniedRange := copyRangeElement(rangeForNextRule)

				if acceptedRange.items[rule.property].minimum < rule.compareTo {

					acceptedRange.items[rule.property] = Range{
						minimum: acceptedRange.items[rule.property].minimum,
						maximum: min(acceptedRange.items[rule.property].maximum, rule.compareTo-1),
					}

					if rule.destination == "A" {

						acceptedRanges.ranges = append(acceptedRanges.ranges, copyRangeElement(acceptedRange))

					} else if rule.destination != "R" {

						calculateAcceptedRanges(acceptedRange, game.workflows[rule.destination], game, acceptedRanges)

					}
				}

				if deniedRange.items[rule.property].maximum >= rule.compareTo {

					deniedRange.items[rule.property] = Range{
						minimum: max(deniedRange.items[rule.property].minimum, rule.compareTo),
						maximum: deniedRange.items[rule.property].maximum,
					}
					rangeForNextRule = deniedRange
				}
			case ">":

				acceptedRange := copyRangeElement(rangeForNextRule)
				deniedRange := copyRangeElement(rangeForNextRule)

				if acceptedRange.items[rule.property].maximum > rule.compareTo {

					acceptedRange.items[rule.property] = Range{
						minimum: max(acceptedRange.items[rule.property].minimum, rule.compareTo+1),
						maximum: acceptedRange.items[rule.property].maximum,
					}

					if rule.destination == "A" {

						acceptedRanges.ranges = append(acceptedRanges.ranges, copyRangeElement(acceptedRange))

					} else if rule.destination != "R" {

						calculateAcceptedRanges(acceptedRange, game.workflows[rule.destination], game, acceptedRanges)

					}
				}

				if deniedRange.items[rule.property].minimum <= rule.compareTo {

					deniedRange.items[rule.property] = Range{
						minimum: deniedRange.items[rule.property].minimum,
						maximum: min(deniedRange.items[rule.property].maximum, rule.compareTo),
					}
					rangeForNextRule = deniedRange
				}
			}
		}
	}
}

func getPossibilities(element RangeElement) int {

	var result = 1

	for _, value := range element.items {

		result *= value.maximum - value.minimum + 1
	}

	return result
}

func copyRangeElement(element RangeElement) RangeElement {

	var newElement RangeElement

	newElement.items = make(map[string]Range)

	for id, v := range element.items {

		newElement.items[id] = v
	}

	return newElement
}

func Part1(input string) string {

	content := GetContent(input)

	game := parseGame(content)

	score := 0

	for _, element := range game.elements {

		if isElementAccepted(element, game.workflows) {

			score += getElementScore(element)
		}
	}

	return strconv.Itoa(score)
}

func Part2(input string) string {

	content := GetContent(input)

	game := parseGame(content)

	var initialRange RangeElement

	initialRange.items = make(map[string]Range)
	initialRange.items["x"] = Range{
		minimum: 1,
		maximum: 4000,
	}
	initialRange.items["m"] = Range{
		minimum: 1,
		maximum: 4000,
	}
	initialRange.items["a"] = Range{
		minimum: 1,
		maximum: 4000,
	}
	initialRange.items["s"] = Range{
		minimum: 1,
		maximum: 4000,
	}

	var acceptedRanges RangeSet

	calculateAcceptedRanges(initialRange, game.workflows["in"], game, &acceptedRanges)

	var result = 0

	for _, acceptedRange := range acceptedRanges.ranges {

		result += getPossibilities(acceptedRange)
	}

	return strconv.Itoa(result)
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
