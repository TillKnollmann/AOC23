package main

import (
	"fmt"
	"testing"
)

var P1_IN_TEST = []string{fmt.Sprintf("test/%s/in01.txt", DAY)}
var P1_OUT_TEST = []string{"136"}

var P2_IN_TEST = []string{fmt.Sprintf("test/%s/in01.txt", DAY)}
var P2_OUT_TEST = []string{"64"}

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
