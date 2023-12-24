package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const DAY = "22"

type Range struct {
	lower int
	upper int
}

type Brick struct {
	x Range
	y Range
	z Range
}

func stringify(brick Brick) string {

	return fmt.Sprintf("(%#v,%#v,%#v)", brick.x, brick.y, brick.z)
}

func (r Range) String() string {

	return fmt.Sprintf("%d-%d", r.lower, r.upper)
}

type Bricks []Brick

func parseBricks(input string) Bricks {

	var bricks Bricks

	for _, line := range strings.Split(strings.ReplaceAll(input, "\r", ""), "\n") {

		var brick Brick
		parts := strings.Split(line, "~")

		for i := 0; i < 3; i++ {

			left := stringToNumber(strings.Split(parts[0], ",")[i])
			right := stringToNumber(strings.Split(parts[1], ",")[i])

			theRange := Range{
				lower: min(left, right),
				upper: max(left, right),
			}
			switch i {
			case 0:
				brick.x = theRange
			case 1:
				brick.y = theRange
			case 2:
				brick.z = theRange
			}

		}

		bricks = append(bricks, brick)
	}

	return bricks
}

func bricksEqual(b1 Brick, b2 Brick) bool {

	return rangesEqual(b1.x, b2.x) && rangesEqual(b1.y, b2.y) && rangesEqual(b1.z, b2.z)
}

func rangesEqual(r1 Range, r2 Range) bool {

	return r1.lower == r2.lower && r1.upper == r2.upper
}

func (b Bricks) Len() int {

	return len(b)
}

func (b Bricks) Swap(i, j int) {

	b[i], b[j] = b[j], b[i]
}

func (b Bricks) Less(i, j int) bool {

	return b[i].z.lower <= b[j].z.lower
}

func sortBricksByZ(bricks Bricks) Bricks {

	sortedBricks := make(Bricks, len(bricks))
	copy(sortedBricks, bricks)
	sort.Sort(sortedBricks)
	return sortedBricks
}

func getOverlappingBricks(brick Brick, bricks Bricks) Bricks {

	var result Bricks

	for _, otherBrick := range bricks {

		if overlapIgnoringZ(brick, otherBrick) && !bricksEqual(brick, otherBrick) {

			result = append(result, otherBrick)
		}
	}

	return result
}

func overlapIgnoringZ(b1 Brick, b2 Brick) bool {

	return overlap1D(b1.x, b2.x) && overlap1D(b1.y, b2.y)
}

func overlap1D(r1 Range, r2 Range) bool {

	return !(r1.lower > r2.upper || r1.upper < r2.lower)
}

func dropBricks(bricks Bricks) (Bricks, int) {

	sortedBricks := sortBricksByZ(bricks)

	numberOfDroppedBricks := 0

	var droppedBricks Bricks

	for _, brick := range sortedBricks {

		droppedBrick, dropped := dropBrick(brick, droppedBricks)
		droppedBricks = append(droppedBricks, droppedBrick)

		if dropped {

			numberOfDroppedBricks++
		}
	}

	return droppedBricks, numberOfDroppedBricks
}

func dropBrick(brick Brick, droppedBricks Bricks) (Brick, bool) {

	height := brick.z.upper - brick.z.lower

	newBrick := Brick{
		x: Range{
			lower: brick.x.lower,
			upper: brick.x.upper,
		},
		y: Range{
			lower: brick.y.lower,
			upper: brick.y.upper,
		},
		z: Range{
			lower: 1,
			upper: 1 + height,
		},
	}

	for _, droppedBrick := range droppedBricks {

		if overlapIgnoringZ(droppedBrick, newBrick) {

			newBrick.z.lower = max(newBrick.z.lower, droppedBrick.z.upper+1)
			newBrick.z.upper = max(newBrick.z.upper, droppedBrick.z.upper+1+height)
		}
	}

	return newBrick, newBrick.z.lower != brick.z.lower
}

func getNumberOfSafeBricks(bricks Bricks) (int, Bricks) {

	bricks = sortBricksByZ(bricks)

	count := 0

	var unsafeBricks Bricks

	overlappingBricksAbove := make(map[string]Bricks)
	overlappingBricksBelow := make(map[string]Bricks)

	for index, brick := range bricks {

		brickString := stringify(brick)

		overlappingBricksAbove[brickString] = getOverlappingBricks(brick, bricks[index+1:])
		overlappingBricksBelow[brickString] = getOverlappingBricks(brick, bricks[:index])
	}

	for _, brick := range bricks {

		brickString := stringify(brick)

		if above, _ := overlappingBricksAbove[brickString]; len(above) == 0 {

			count++

		} else {

			isSafe := true

			for _, aboveBrick := range overlappingBricksAbove[brickString] {

				if overlapIgnoringZ(brick, aboveBrick) && brick.z.upper == aboveBrick.z.lower-1 {

					bricksBelowTheAbove, _ := overlappingBricksBelow[stringify(aboveBrick)]

					if len(bricksBelowTheAbove) == 1 {

						isSafe = false
						unsafeBricks = append(unsafeBricks, brick)
						break
					}

					isAnotherSupporting := false

					for _, otherBrickBelow := range bricksBelowTheAbove {

						if !bricksEqual(brick, otherBrickBelow) && overlapIgnoringZ(otherBrickBelow, aboveBrick) && otherBrickBelow.z.upper == aboveBrick.z.lower-1 {

							isAnotherSupporting = true
							break
						}
					}

					if !isAnotherSupporting {

						isSafe = false
						unsafeBricks = append(unsafeBricks, brick)
						break
					}
				}
			}

			if isSafe {

				count++

			}
		}

	}

	return count, unsafeBricks
}

func Part1(input string) string {

	content := GetContent(input)

	bricks := parseBricks(content)

	bricks, _ = dropBricks(bricks)

	numberOfSafeBricks, _ := getNumberOfSafeBricks(bricks)

	return strconv.Itoa(numberOfSafeBricks)
}

func Part2(input string) string {

	content := GetContent(input)

	bricks := parseBricks(content)

	bricks, _ = dropBricks(bricks)

	_, unsafeBricks := getNumberOfSafeBricks(bricks)

	sum := 0

	for _, brick := range unsafeBricks {

		reducedBrickSet := make(Bricks, len(bricks))
		copy(reducedBrickSet, bricks)

		reducedBrickSet = slices.DeleteFunc(reducedBrickSet, func(currentBrick Brick) bool {

			return bricksEqual(brick, currentBrick)
		})

		_, droppedBricks := dropBricks(reducedBrickSet)

		sum += droppedBricks
	}

	return strconv.Itoa(sum)
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
