package main

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"log"
	"math"
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
	position Vector
	velocity Vector
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

	storm.position = parseVector(parts[0])
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
			cut := findIntersectionIgnoreZ(stormA.position, stormB.position, stormA.velocity, stormB.velocity)

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

func getInitialPositions(stormA, stormB, stormC Hailstorm) Vector {

	// Let x, y, z, dx, dy, dz, be position and velocity of the rock.
	// A collision of the rock requires for storm = stormA, stormB, stormC for time t = t1, t2, t3:
	// [x, y, z] + t [dx, dy, dz] = [storm.position.x, storm.position.y, storm.position.z] + t [storm.velocity.x, storm.velocity.y, storm.velocity.z]
	// Re-arranging for t yields for one storm:
	// (1) t = (x-storm.position.x)/(storm.velocity.x - dx)
	// (2) t = (y-storm.position.y)/(storm.velocity.y - dy)
	// (3) t = (z-storm.position.z)/(storm.velocity.z - dz)
	// All these can be set equal for three equations that can be re-arranged to
	// (4) y*dx - x*dy = y*storm.velocity.x - storm.position.y*storm.velocity.x + storm.position.y*dx - x*storm.velocity.y + storm.position.x*storm.velocity.y - storm.position.x*dy
	// (5) y*dx - x*dy = z*storm.velocity.y - storm.position.z*storm.velocity.y + storm.position.z*dy - y*storm.velocity.z + storm.position.y*storm.velocity.z - storm.position.y*dz
	// (6) y*dx - x*dy = x*storm.velocity.z - storm.position.x*storm.velocity.z + storm.position.x*dz - z*storm.velocity.x + storm.position.z*storm.velocity.x - storm.position.z*dx
	// Each left hand side is independent of the storm.
	// Thus, combining left hand sides for stormA + stormB and equalizing them yields
	// (7) (stormB.velocity.y - stormA.velocity.y) x + (stormA.velocity.x - stormB.velocity.x) y + (stormA.position.y - stormB.position.y) dx + (stormB.position.x - stormA.position.x) dy = stormB.position.x stormB.velocity.y - stormB.position.y stormB.velocity.x - stormA.position.x stormA.velocity.y + stormA.position.y stormA.velocity.x
	// (8) (stormB.velocity.z - stormA.velocity.z) x + (stormA.velocity.x - stormB.velocity.x) z + (stormA.position.z - stormB.position.z) dx + (stormB.position.x - stormA.position.x) dz = stormB.position.x stormB.velocity.z - stormB.position.z stormB.velocity.x - stormA.position.x stormA.velocity.z + stormA.position.z stormA.velocity.x
	// (9) (stormA.velocity.z - stormB.velocity.z) y + (stormB.velocity.y - stormA.velocity.y) z + (stormB.position.z - stormA.position.z) dy + (stormA.position.y - stormB.position.y) dz = - stormB.position.y stormB.velocity.z + stormB.position.z stormB.velocity.y + stormA.position.y stormA.velocity.z - stormA.position.z stormA.velocity.y
	// Combining left hand sides for stormA and stormC and equalizing them the same way yields an equation system with 6 unknown (x,y,z,dx,dy,dz) and 6 equations.

	equationSystem := [][]float64{
		getEquationTypeOne(stormA, stormB),
		getEquationTypeOne(stormA, stormC),
		getEquationTypeTwo(stormA, stormB),
		getEquationTypeTwo(stormA, stormC),
		getEquationTypeThree(stormA, stormB),
		getEquationTypeThree(stormA, stormC)}

	solution := solveLinearEquationSystem(equationSystem)

	var vector Vector

	vector.x = solution[0]
	vector.y = solution[1]
	vector.z = solution[2]

	return vector
}

func solveLinearEquationSystem(A [][]float64) []float64 {

	// Create a new matrix from the 2D slice.
	matA := mat.NewDense(len(A), len(A[0])-1, nil)
	for i := 0; i < len(A); i++ {
		for j := 0; j < len(A[0])-1; j++ {
			matA.Set(i, j, A[i][j])
		}
	}

	// Create a new vector from the last column of A.
	b := mat.NewVecDense(len(A), nil)
	for i := 0; i < len(A); i++ {
		b.SetVec(i, A[i][len(A[0])-1])
	}

	// Solve the linear equation system.
	var x mat.VecDense
	err := x.SolveVec(matA, b)
	if err != nil {
		panic(err)
	}

	// Return the solution as a slice.
	return x.RawVector().Data
}

func getEquationTypeOne(a, b Hailstorm) []float64 {

	// (stormB.velocity.y - stormA.velocity.y) x
	// + (stormA.velocity.x - stormB.velocity.x) y
	// + (stormA.position.y - stormB.position.y) dx
	// + (stormB.position.x - stormA.position.x) dy
	// = stormB.position.x stormB.velocity.y - stormB.position.y stormB.velocity.x - stormA.position.x stormA.velocity.y + stormA.position.y stormA.velocity.x
	arr := make([]float64, 7)
	arr[0] = b.velocity.y - a.velocity.y
	arr[1] = a.velocity.x - b.velocity.x
	arr[3] = a.position.y - b.position.y
	arr[4] = b.position.x - a.position.x
	arr[6] = b.position.x*b.velocity.y - b.position.y*b.velocity.x - a.position.x*a.velocity.y + a.position.y*a.velocity.x
	return arr
}

func getEquationTypeTwo(a, b Hailstorm) []float64 {

	// (stormB.velocity.z - stormA.velocity.z) x
	// + (stormA.velocity.x - stormB.velocity.x) z
	// + (stormA.position.z - stormB.position.z) dx
	// + (stormB.position.x - stormA.position.x) dz
	// = stormB.position.x stormB.velocity.z - stormB.position.z stormB.velocity.x - stormA.position.x stormA.velocity.z + stormA.position.z stormA.velocity.x
	arr := make([]float64, 7)
	arr[0] = b.velocity.z - a.velocity.z
	arr[2] = a.velocity.x - b.velocity.x
	arr[3] = a.position.z - b.position.z
	arr[5] = b.position.x - a.position.x
	arr[6] = b.position.x*b.velocity.z - b.position.z*b.velocity.x - a.position.x*a.velocity.z + a.position.z*a.velocity.x
	return arr
}

func getEquationTypeThree(a, b Hailstorm) []float64 {

	// (stormA.velocity.z - stormB.velocity.z) y
	// + (stormB.velocity.y - stormA.velocity.y) z
	// + (stormB.position.z - stormA.position.z) dy
	// + (stormA.position.y - stormB.position.y) dz
	// = - stormB.position.y stormB.velocity.z + stormB.position.z stormB.velocity.y + stormA.position.y stormA.velocity.z - stormA.position.z stormA.velocity.y
	arr := make([]float64, 7)
	arr[1] = a.velocity.z - b.velocity.z
	arr[2] = b.velocity.y - a.velocity.y
	arr[4] = b.position.z - a.position.z
	arr[5] = a.position.y - b.position.y
	arr[6] = -b.position.y*b.velocity.z + b.position.z*b.velocity.y + a.position.y*a.velocity.z - a.position.z*a.velocity.y
	return arr
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

	storms := parseStorms(content)

	result := getInitialPositions(storms[0], storms[1], storms[2])

	return strconv.Itoa(int(math.Ceil(result.x + result.y + result.z)))
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
