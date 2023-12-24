package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const DAY = "23"

type Game struct {
	fields [][]rune
	start  Point
	end    Point
}

type Point struct {
	x int
	y int
}

type Direction int

const (
	North Direction = iota + 1
	East
	South
	West
)

type Agent struct {
	position   Point
	direction  Direction
	pathLength int
}

type Graph struct {
	nodes map[string]Node
}

type Node struct {
	position Point
	adjacent map[string]int
}

func parseGame(input string) Game {

	lines := strings.Split(strings.ReplaceAll(input, "\r", ""), "\n")

	var game Game

	for y := 0; y < len(lines); y++ {

		var row []rune

		for x := 0; x < len(lines[y]); x++ {

			row = append(row, rune(lines[y][x]))

			if y == 0 && rune(lines[y][x]) == '.' {

				game.start = Point{x: x, y: y}

			} else if y == len(lines)-1 && rune(lines[y][x]) == '.' {

				game.end = Point{x: x, y: y}
			}
		}

		game.fields = append(game.fields, row)
	}

	return game
}

func stringify(point Point) string {

	return fmt.Sprintf("(%d,%d)", point.x, point.y)
}

func parse(pointString string) Point {

	var point Point

	numberRe := regexp.MustCompile(`\d+`)

	matches := numberRe.FindAllString(pointString, 2)

	point.x = stringToNumber(matches[0])
	point.y = stringToNumber(matches[1])

	return point
}

func getGraph(game Game, climbSteeps bool) Graph {

	var graph Graph

	graph.nodes = make(map[string]Node)

	addNode(&graph, game.start)
	addNode(&graph, game.end)

	getGraphFrom(game, &graph, game.start, Agent{
		position:   game.start,
		direction:  South,
		pathLength: 0,
	}, climbSteeps)

	return graph
}

func getGraphFrom(game Game, graph *Graph, lastNode Point, agent Agent, climbSteeps bool) {

	// are we at a known node?
	if _, nodeExists := graph.nodes[stringify(agent.position)]; nodeExists && stringify(agent.position) != stringify(game.start) {

		addEdge(graph, lastNode, agent.position, agent.pathLength, !climbSteeps)
		return
	}

	// get possible follow-up directions
	directions := []Direction{North, East, West, South}

	// we cannot go back
	directions = slices.DeleteFunc(directions, func(direction Direction) bool {

		return direction == getOpposite(agent.direction)
	})

	if !climbSteeps {

		// remove invalid all other directions if the point is not .
		directions = slices.DeleteFunc(directions, func(direction Direction) bool {

			switch game.fields[agent.position.y][agent.position.x] {
			case '>':
				return direction != East
			case '<':
				return direction != West
			case 'v':
				return direction != South
			case '^':
				return direction != North

			}
			return false
		})
	}

	// remove directions that do not lead to a valid point
	directions = slices.DeleteFunc(directions, func(direction Direction) bool {

		nextPoint := getNextPoint(agent.position, direction)

		if !isPointInRange(nextPoint, game) {

			return true
		}

		if game.fields[nextPoint.y][nextPoint.x] == '#' {

			return true
		}

		return false
	})

	// identify new node
	if len(directions) > 1 {

		addNode(graph, agent.position)
		addEdge(graph, lastNode, agent.position, agent.pathLength, !climbSteeps)
		agent.pathLength = 0
		lastNode = agent.position
	}

	// recurse to future points
	var nextAgents []Agent

	for _, direction := range directions {

		nextAgent := Agent{
			position:   getNextPoint(agent.position, direction),
			direction:  direction,
			pathLength: agent.pathLength + 1,
		}

		nextAgents = append(nextAgents, nextAgent)
	}

	for _, nextAgent := range nextAgents {

		getGraphFrom(game, graph, lastNode, nextAgent, climbSteeps)
	}
}

func addNode(graph *Graph, point Point) {

	if _, nodeExists := graph.nodes[stringify(point)]; !nodeExists {

		graph.nodes[stringify(point)] = Node{
			position: Point{x: point.x, y: point.y},
			adjacent: make(map[string]int),
		}
	}
}

func addEdge(graph *Graph, source Point, target Point, weight int, directed bool) {

	addNode(graph, source)
	addNode(graph, target)

	sourceString := stringify(source)
	targetString := stringify(target)

	graph.nodes[sourceString].adjacent[targetString] = weight

	if !directed {

		graph.nodes[targetString].adjacent[sourceString] = weight
	}
}

func getLongestPathLength(game Game, graph Graph) int {

	return getLongestPathLengthNode(game, graph, game.start, []string{})
}

func getLongestPathLengthNode(game Game, graph Graph, currentNode Point, lastVisitedNodes []string) int {

	currentNodeString := stringify(currentNode)

	if slices.Contains(lastVisitedNodes, currentNodeString) {

		return math.MinInt
	}

	if game.end.x == currentNode.x && game.end.y == currentNode.y {

		return 0
	}

	var pathLengths []int

	node, _ := graph.nodes[currentNodeString]

	nextVisited := append(lastVisitedNodes, currentNodeString)

	for nextNode, edgeWeight := range node.adjacent {

		pathLengths = append(pathLengths, edgeWeight+getLongestPathLengthNode(game, graph, parse(nextNode), nextVisited))
	}

	return slices.Max(pathLengths)
}

func getNextPoint(point Point, direction Direction) Point {

	deltaX := 0
	deltaY := 0

	switch direction {
	case North:
		deltaY--
	case East:
		deltaX++
	case South:
		deltaY++
	case West:
		deltaX--
	}

	return Point{x: point.x + deltaX, y: point.y + deltaY}
}

func isPointInRange(point Point, game Game) bool {

	return point.x >= 0 && point.x < len(game.fields[0]) && point.y >= 0 && point.y < len(game.fields)
}

func getOpposite(direction Direction) Direction {

	switch direction {
	case North:
		return South
	case East:
		return West
	case South:
		return North
	case West:
		return East
	}

	return North
}

func Part1(input string) string {

	content := GetContent(input)

	game := parseGame(content)

	graph := getGraph(game, false)

	maxPathLength := getLongestPathLength(game, graph)

	return strconv.Itoa(maxPathLength)
}

func Part2(input string) string {

	content := GetContent(input)

	game := parseGame(content)

	graph := getGraph(game, true)

	maxPathLength := getLongestPathLength(game, graph)

	return strconv.Itoa(maxPathLength)
}

func GetContent(filepath string) string {

	content, err := os.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func stringToNumber(s string) int {

	number, err := strconv.ParseInt(s, 10, 64)

	if err != nil {

		panic(err)
	}

	return int(number)
}

func main() {

	fmt.Println(fmt.Sprintf("Part 1: %s", Part1(fmt.Sprintf("input/%s/in.txt", DAY))))
	fmt.Println(fmt.Sprintf("Part 2: %s", Part2(fmt.Sprintf("input/%s/in.txt", DAY))))
}
