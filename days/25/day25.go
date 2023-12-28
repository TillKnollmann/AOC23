package main

import (
	"cmp"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
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
	components := [][]string{}

	for _, node := range graph.nodes {
		if !visited[node] {
			component := []string{}
			component = DFS(graph, node, visited, component)
			components = append(components, component)
		}
	}

	return components
}

func getConnectedComponentSizes(graph Graph) []int {
	components := getConnectedComponents(graph)
	sizes := []int{}

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

func getBridgeWires(graph Graph) int {

	for i := 0; i < len(graph.edges)-2; i++ {

		for j := 0; j < len(graph.edges)-1; j++ {

			for k := 0; k < len(graph.edges); k++ {

				gPrime := Graph{nodes: graph.nodes, edges: removeEdges(graph.edges, []Edge{graph.edges[i], graph.edges[j], graph.edges[k]})}
				connectedComponentSizes := getConnectedComponentSizes(gPrime)

				if len(connectedComponentSizes) == 2 {

					return connectedComponentSizes[0] * connectedComponentSizes[1]
				}
			}
		}
	}

	return -1
}

func getBridgeWiresParallel(graph Graph) int {
	var wg sync.WaitGroup
	resultChan := make(chan int)

	for i := 0; i < len(graph.edges)-2; i++ {
		for j := 0; j < len(graph.edges)-1; j++ {
			for k := 0; k < len(graph.edges); k++ {
				wg.Add(1)
				go func(i, j, k int) {
					defer wg.Done()
					gPrime := Graph{nodes: graph.nodes, edges: removeEdges(graph.edges, []Edge{graph.edges[i], graph.edges[j], graph.edges[k]})}
					connectedComponentSizes := getConnectedComponentSizes(gPrime)
					if len(connectedComponentSizes) == 2 {
						resultChan <- connectedComponentSizes[0] * connectedComponentSizes[1]
					}
				}(i, j, k)
			}
		}
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	maximum := -1
	for result := range resultChan {
		if result > maximum {
			maximum = result
		}
	}

	return maximum
}

func Part1(input string) string {

	content := GetContent(input)

	graph := parseGraph(content)

	return strconv.Itoa(getBridgeWiresParallel(graph))
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
