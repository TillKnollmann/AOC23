package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const DAY = "19"

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

	return string(content[0])
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
