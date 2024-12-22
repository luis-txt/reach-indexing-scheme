package main

import (
	"bufio"
	"fmt"
	"os"
)

func ReadGraph(path string) *Graph {
	logger.Println("Reading Graph...")

	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file: ", err)
		os.Exit(1)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	n, _ := -1, -1

	if scanner.Scan() {
		line := scanner.Text()

		fmt.Sscanf(line, "n: %d", &n)
	} else {
		fmt.Fprintln(os.Stderr, "Error reading header of file: ", err)
		os.Exit(1)
	}
	g := CreateGraph(n)
	g.idMapping = IdMapping{make(map[int]int), make(map[int]int)}
	nextIndex := 0

	sourceID, targetID := -1, -1
	for scanner.Scan() {
		line := scanner.Text()
		_, err := fmt.Sscanf(line, "%d %d", &sourceID, &targetID)
		if err == nil {
			nextIndex = g.AddVtoIdMapping(sourceID, nextIndex)
			nextIndex = g.AddVtoIdMapping(targetID, nextIndex)
			v := g.idMapping.idToV[sourceID]
			w := g.idMapping.idToV[targetID]

			e := Edge{v, w, nil, nil, nil}
			g.AddEdge(&e)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading edges of file: ", err)
		os.Exit(1)
	}
	return g
}

func (g *Graph) WriteGraph(path string, edgeValues []int) {
	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		// Create file if it does not already exist.
		// Open file in write-only.
		// If file already exists, clear it.
		0644,
		// Owner of file has read and write permissions.
		// The group and others have read permissions.
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening file: ", err)
		os.Exit(1)
	}

	defer file.Close()

	// Write header
	header := fmt.Sprintf("n: %d\n", g.n)
	_, err = file.WriteString(header)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error writing header of file: ", err)
		os.Exit(1)
	}
	// Write edges
	edgeNumber := 0
	for i := 0; i < g.n; i++ {
		for e := g.nodes[i].out; e != nil; e = e.next {
			edgeStr := fmt.Sprintf("%d %d", e.source, e.target)
			if edgeValues != nil {
				edgeStr = fmt.Sprintf("%s %d", edgeStr, edgeValues[edgeNumber])
			}
			edgeStr = fmt.Sprintf("%s\n", edgeStr)
			_, err = file.WriteString(edgeStr)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error writing edges of file: ", err)
				os.Exit(1)

			}
			edgeNumber++
		}
	}
}
