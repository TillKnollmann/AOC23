package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const DAY = "24"

type Crossing struct {
	stormA Hailstorm
	stormB Hailstorm
	cut    Vector
}

type Vector struct {
	x float64
	y float64
	z float64
}

func (v Vector) String() string {

	return fmt.Sprintf("(%f,%f,%f)", v.x, v.y, v.z)
}

type Hailstorm struct {
	positions map[int]Vector
	velocity  Vector
}

func parseStorms(input string) []Hailstorm {

	var storms []Hailstorm

	for _, line := range strings.Split(strings.ReplaceAll(input, "\r", ""), "\n") {

		storms = append(storms, parseHailstorm(line))
	}

	return storms
}

func parseHailstorm(input string) Hailstorm {

	var storm Hailstorm

	parts := strings.Split(strings.ReplaceAll(input, " ", ""), "@")

	storm.positions = make(map[int]Vector)
	storm.positions[0] = parseVector(parts[0])
	storm.velocity = parseVector(parts[1])

	return storm
}

func parseVector(input string) Vector {

	var Vector Vector

	parts := strings.Split(input, ",")

	Vector.x = float64(stringToNumber(parts[0]))
	Vector.y = float64(stringToNumber(parts[1]))
	Vector.z = float64(stringToNumber(parts[2]))

	return Vector
}

func findIntersectionIgnoreZ(positionA Vector, positionB Vector, velocityA Vector, velocityB Vector) *Vector {

	// https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection#Given_two_points_on_each_line_segment

	x1 := positionA.x
	x2 := positionA.x + velocityA.x
	x3 := positionB.x
	x4 := positionB.x + velocityB.x
	y1 := positionA.y
	y2 := positionA.y + velocityA.y
	y3 := positionB.y
	y4 := positionB.y + velocityB.y

	t := ((x1-x3)*(y3-y4) - (y1-y3)*(x3-x4)) / ((x1-x2)*(y3-y4) - (y1-y2)*(x3-x4))
	u := ((x1-x3)*(y1-y2) - (y1-y3)*(x1-x2)) / ((x1-x2)*(y3-y4) - (y1-y2)*(x3-x4))

	if t < 0 || u < 0 {

		return nil
	}

	return &Vector{
		x: x1 + t*(x2-x1),
		y: y1 + t*(y2-y1),
		z: 0,
	}
}

func getNumberOfIntersectionsIgnoringZ(storms []Hailstorm, minimumPosition Vector, maximumPosition Vector) int {

	var crossings []Crossing

	for i := 0; i < len(storms)-1; i++ {

		stormA := storms[i]
		for j := i + 1; j < len(storms); j++ {

			stormB := storms[j]
			cut := findIntersectionIgnoreZ(stormA.positions[0], stormB.positions[0], stormA.velocity, stormB.velocity)

			if cut != nil {

				if minimumPosition.x <= cut.x && cut.x <= maximumPosition.x && minimumPosition.y <= cut.y && cut.y <= maximumPosition.y {

					crossings = append(crossings, Crossing{
						stormA: stormA,
						stormB: stormB,
						cut:    *cut,
					})
				}
			}
		}
	}

	return len(crossings)
}

func Part1(input string, testStart int, testEnd int) string {

	content := GetContent(input)

	storms := parseStorms(content)

	return strconv.Itoa(getNumberOfIntersectionsIgnoringZ(storms, Vector{
		x: float64(testStart),
		y: float64(testStart),
		z: 0,
	}, Vector{
		x: float64(testEnd),
		y: float64(testEnd),
		z: 0,
	}))
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

func stringToNumber(s string) int {

	number, err := strconv.ParseInt(s, 10, 64)

	if err != nil {

		panic(err)
	}

	return int(number)
}

func main() {

	fmt.Println(fmt.Sprintf("Part 1: %s", Part1(fmt.Sprintf("input/%s/in.txt", DAY), 200000000000000, 400000000000000)))
	fmt.Println(fmt.Sprintf("Part 2: %s", Part2(fmt.Sprintf("input/%s/in.txt", DAY))))
}
