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

type Graph struct {
	nodes []string
	edges []Edge
}

func DFS(graph Graph, node string, visited map[string]bool, component []string) []string {

	visited[node] = true
	component = append(component, node)

	for _, edge := range graph.edges {
		if edge.in == node && !visited[edge.out] {
			component = DFS(graph, edge.out, visited, component)
		} else if edge.out == node && !visited[edge.in] {
			component = DFS(graph, edge.in, visited, component)
		}
	}

	return component
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

func getConnectedComponents(graph Graph) [][]string {

	visited := make(map[string]bool)
	var components [][]string

	for _, node := range graph.nodes {

		if !visited[node] {

			var component []string
			component = DFS(graph, node, visited, component)
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

func getEdgesToNumberOfOccurrences(graph Graph) map[Edge]int {

	shortest := make(map[Edge]int)
	for _, e := range graph.edges {

		shortest[e] = 1
	}

	for _, n := range graph.nodes {

		dist := make(map[string]int)
		prev := make(map[string]string)
		for _, m := range graph.nodes {

			dist[m] = -1
		}
		dist[n] = 0

		q := []string{n}
		for len(q) > 0 {

			u := q[0]
			q = q[1:]

			for _, e := range graph.edges {

				if e.in == u {

					v := e.out
					if dist[v] == -1 {

						dist[v] = dist[u] + 1
						prev[v] = u
						q = append(q, v)

					} else if dist[v] == dist[u]+1 {

						shortest[e]++

					}

				} else if e.out == u {

					v := e.in
					if dist[v] == -1 {

						dist[v] = dist[u] + 1
						prev[v] = u
						q = append(q, v)

					} else if dist[v] == dist[u]+1 {

						shortest[e]++

					}
				}
			}
		}
	}

	return shortest
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

type Pair struct {
	Key   Edge
	Value int
}

func Part1(input string) string {

	content := GetContent(input)

	graph := parseGraph(content)

	edges := getEdgesToNumberOfOccurrences(graph)

	sortedEdges := sortMapByValue(edges)

	chunkSize := 10

	var calculatedCombinations map[string]bool = make(map[string]bool)

	for c := 0; c < len(sortedEdges)/chunkSize; c++ {

		for i := 0; i < min(chunkSize*c, len(sortedEdges))-2; i++ {

			edgeA := sortedEdges[i].Key

			for j := i + 1; j < min(chunkSize*c, len(sortedEdges))-1; j++ {

				edgeB := sortedEdges[j].Key

				for k := j + 1; k < min(chunkSize*c, len(sortedEdges)); k++ {

					_, skip := calculatedCombinations[fmt.Sprintf("%d,%d,%d", i, j, k)]

					if skip {

						continue
					}

					edgeC := sortedEdges[k].Key

					gPrime := Graph{nodes: graph.nodes, edges: removeEdges(graph.edges, []Edge{edgeA, edgeB, edgeC})}
					connectedComponentSizes := getConnectedComponentSizes(gPrime)

					if len(connectedComponentSizes) == 2 {

						return strconv.Itoa(connectedComponentSizes[0] * connectedComponentSizes[1])
					}

					calculatedCombinations[fmt.Sprintf("%d,%d,%d", i, j, k)] = true
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
