package main

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const DAY = "25"

type Edge struct {
	in  string
	out string
}

func (e Edge) String() string {

	return fmt.Sprintf("{%s,%s}", e.in, e.out)
}

type Graph struct {
	nodes []string
	edges []Edge
}

func parseGraph(input string) Graph {

	var graph Graph

	for _, line := range strings.Split(strings.ReplaceAll(input, "\r", ""), "\n") {

		nodes := strings.Split(strings.ReplaceAll(line, ":", ""), " ")

		for _, node := range nodes {

			if !slices.Contains(graph.nodes, node) {

				graph.nodes = append(graph.nodes, node)
			}
		}

		for _, node := range nodes[1:] {

			if !slices.Contains(graph.edges, Edge{
				in:  nodes[0],
				out: node,
			}) {

				graph.edges = append(graph.edges, Edge{
					in:  nodes[0],
					out: node,
				})
			}
		}
	}

	slices.SortFunc(graph.nodes, func(a, b string) int {

		return cmp.Compare(a, b)
	})

	return graph
}

func BFSAllNodes(graph Graph) map[Edge]int {

	edgePaths := make(map[Edge]int)

	for _, node := range graph.nodes {

		_, visitedEdges := BFS(graph, node)

		for edge, count := range visitedEdges {

			edgePaths[edge] += count
		}
	}
	return edgePaths
}

func BFS(graph Graph, start string) ([]string, map[Edge]int) {

	visited := make(map[string]bool)

	for _, node := range graph.nodes {

		visited[node] = false
	}

	visited[start] = true
	visitedNodes := []string{start}
	edgePaths := make(map[Edge]int)

	for _, edge := range graph.edges {

		edgePaths[edge] = 0
	}

	queue := []string{start}

	for len(queue) > 0 {

		node := queue[0]
		queue = queue[1:]

		for _, edge := range graph.edges {

			if edge.in == node && !visited[edge.out] {

				visited[edge.out] = true
				visitedNodes = append(visitedNodes, edge.out)
				queue = append(queue, edge.out)
				edgePaths[edge] += 1

			} else if edge.out == node && !visited[edge.in] {

				visited[edge.in] = true
				visitedNodes = append(visitedNodes, edge.in)
				queue = append(queue, edge.in)
				edgePaths[edge] += 1

			}
		}
	}

	return visitedNodes, edgePaths
}

func getConnectedComponents(graph Graph) [][]string {

	visited := make(map[string]bool)
	var components [][]string

	for _, node := range graph.nodes {

		if !visited[node] {

			var component []string
			component, _ = BFS(graph, node)

			for _, visitedNode := range component {

				visited[visitedNode] = true
			}

			components = append(components, component)
		}
	}

	return components
}

func getConnectedComponentSizes(graph Graph) []int {

	components := getConnectedComponents(graph)
	var sizes []int

	for _, component := range components {
		sizes = append(sizes, len(component))
	}

	return sizes
}

func removeEdges(edges []Edge, edgesToRemove []Edge) []Edge {

	var newEdges []Edge

	for _, edge := range edges {

		if !slices.Contains(edgesToRemove, Edge{
			in:  edge.in,
			out: edge.out,
		}) && !slices.Contains(edgesToRemove, Edge{
			in:  edge.out,
			out: edge.in,
		}) {

			newEdges = append(newEdges, edge)
		}
	}

	return newEdges
}

type Pair struct {
	Key   Edge
	Value int
}

func sortMapByValue(m map[Edge]int) []Pair {

	var pairs []Pair
	for k, v := range m {

		pairs = append(pairs, Pair{k, v})
	}

	sort.Slice(pairs, func(i, j int) bool {

		return pairs[i].Value > pairs[j].Value
	})

	return pairs
}

func Part1(input string) string {

	content := GetContent(input)

	graph := parseGraph(content)

	edges := BFSAllNodes(graph)

	sortedEdges := sortMapByValue(edges)

	for i := 0; i < len(sortedEdges)-2; i++ {

		edgeA := sortedEdges[i].Key

		for j := i + 1; j < len(sortedEdges)-1; j++ {

			edgeB := sortedEdges[j].Key

			for k := j + 1; k < len(sortedEdges); k++ {

				edgeC := sortedEdges[k].Key

				gPrime := Graph{nodes: graph.nodes, edges: removeEdges(graph.edges, []Edge{edgeA, edgeB, edgeC})}
				connectedComponentSizes := getConnectedComponentSizes(gPrime)

				if len(connectedComponentSizes) == 2 {

					return strconv.Itoa(connectedComponentSizes[0] * connectedComponentSizes[1])
				}

			}
		}
	}

	return "-1"
}

func GetContent(filepath string) string {

	content, err := os.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func main() {

	fmt.Println(fmt.Sprintf("Part 1: %s", Part1(fmt.Sprintf("input/%s/in.txt", DAY))))
}
