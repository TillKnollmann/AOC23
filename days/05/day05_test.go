package main

import (
	"fmt"
	"testing"
)

var P1_IN_TEST = []string{fmt.Sprintf("test/%s/in01.txt", DAY)}
var P1_OUT_TEST = []string{"35"}

var P2_IN_TEST = []string{fmt.Sprintf("test/%s/in01.txt", DAY)}
var P2_OUT_TEST = []string{"46"}

func TestParseMapping(t *testing.T) {

	testIn := "seed-to-soil map:\n50 98 2\n52 50 5"
	testOut := Mapping{source: "seed", target: "soil", mappingLines: []MappingLine{{destRangeStart: 50, sourceRangeStart: 98, length: 2}, {destRangeStart: 52, sourceRangeStart: 50, length: 5}}}

	result := parseMapping(testIn)

	assert("parseMapping", testIn, fmt.Sprintf("%s", testOut), fmt.Sprintf("%s", result), t)
}

func TestPart1(t *testing.T) {

	for index, element := range P1_IN_TEST {

		expected := P1_OUT_TEST[index]
		received := Part1(element)
		assert("Part1", element, expected, received, t)
	}
}

func TestPart2(t *testing.T) {

	for index, element := range P2_IN_TEST {

		expected := P2_OUT_TEST[index]
		received := Part2(element)
		assert("Part2", element, expected, received, t)
	}
}

func assert(method string, input string, expected string, received string, t *testing.T) {

	if expected != received {

		t.Errorf("%s(%s) expected '%s' but received '%s'", method, input, expected, received)
	}
}
