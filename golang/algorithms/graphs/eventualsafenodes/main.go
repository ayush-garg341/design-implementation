package main

import (
	"fmt"
	"slices"
)

func main() {
	graph := [][]int{
		{1, 2},
		{2, 3},
		{5},
		{0},
		{5},
		{},
		{},
	}
	nodes := safeNodes(graph)
	fmt.Println(nodes)
}

func safeNodes(graph [][]int) []int {

	pathVisited := make([]int, len(graph))
	visited := make([]int, len(graph))
	safeNodes := []int{}
	for i := 0; i < len(graph); i++ {
		if visited[i] != 1 {
			dfs(i, visited, pathVisited, &safeNodes, graph)
		}
	}

	slices.Sort(safeNodes)

	return safeNodes
}

func dfs(v int, visited, pathVisited []int, safeNodes *[]int, graph [][]int) bool {
	visited[v] = 1
	pathVisited[v] = 1
	for _, advV := range graph[v]{
		if visited[advV] != 1 {
			visited[advV] = 1
			pathVisited[advV] = 1
			hasCycle := dfs(advV, visited, pathVisited, safeNodes, graph)
			if hasCycle{
				return true
			}
		} else {
			if pathVisited[advV] == 1{
				return true
			}
		}
	}

	pathVisited[v] = 0
	*safeNodes = append(*safeNodes, v)
	return false
}
