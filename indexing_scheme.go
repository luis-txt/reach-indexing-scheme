package main

import (
	"math"
)

// Creates the indexing scheme in O(|E_{tr}| + k_c * |E_{red}|).
func (g *Graph) CreateIndexingScheme(topo []int, decomp *Decomposition) [][]int {
	indexingScheme := make([][]int, g.n)
	// Initialize indexing scheme
	for v := 0; v < g.n; v++ {
		schemeNodesProcessed++
		vScheme := make([]int, decomp.chains.n)
		indexingScheme[v] = vScheme
		// Set all reachable indices to infinity
		for i := 0; i < len(vScheme); i++ {
			vScheme[i] = math.MaxInt
		}
	}
	// Fill indexing scheme
	for i := len(topo) - 1; i >= 0; i-- {
		v := topo[i]
		schemeNodesProcessed++

		for e := g.nodes[v].out; e != nil; e = e.next {
			schemeEdgesProcessed++
			// Assuming outgoing edges are already sorted in topologgerical order
			tChain := decomp.vToChain[e.target].chain.val
			if indexingScheme[v][tChain.id] >= indexingScheme[e.target][tChain.id] {
				// Update indices
				for j := 0; j < decomp.chains.n; j++ {
					indexingScheme[v][j] = min(indexingScheme[v][j], indexingScheme[e.target][j])
				}
				if indexingScheme[v][tChain.id] > decomp.vToChain[e.target].pos {
					indexingScheme[v][tChain.id] = decomp.vToChain[e.target].pos
				}
			}
		}
	}
	return indexingScheme
}

// Converts a vertex to its respective components to use in the algorithm.
// Returns v if there are no components.
func convertV(v int, g *Graph) int {
	if g.vToComp != nil {
		v = g.vToComp[v]
	}
	return v
}

// Queries the reachability indexing scheme whether s can reach t in O(1).
func isReachable(s, t int, indexingScheme [][]int, decomp *Decomposition, g *Graph) bool {
	s = convertV(s, g)
	t = convertV(t, g)

	if s == t {
		return true
	}
	tChain := decomp.vToChain[t].chain.val
	sIndex := indexingScheme[s][tChain.id]
	tIndex := indexingScheme[t][tChain.id]
	return sIndex < tIndex
}

// Converts the indexing scheme to a reachability matrix in O(|V|^2).
func schemeToMatrix(indexingScheme [][]int, decomp *Decomposition, g *Graph) [][]bool {
	// Create n*n matrix
	matrix := make([][]bool, len(g.vToComp))
	for i := 0; i < len(g.vToComp); i++ {
		matrix[i] = make([]bool, len(g.vToComp))
	}
	// Fill matrix
	for v := 0; v < len(g.vToComp); v++ {
		for w := 0; w < len(g.vToComp); w++ {
			if isReachable(v, w, indexingScheme, decomp, g) {
				matrix[v][w] = true
			}
		}
	}
	return matrix
}

func printScheme(indexingScheme [][]int) {
	for v := 0; v < len(indexingScheme); v++ {
		logger.Println(v, ": [ ")
		for i := 0; i < len(indexingScheme[v]); i++ {
			if indexingScheme[v][i] == math.MaxInt {
				logger.Println("inf")
			} else {
				logger.Println(indexingScheme[v][i])
			}
			if i < len(indexingScheme[v])-1 {
				logger.Println(", ")
			}
		}
		logger.Println(" ]")
	}
}
