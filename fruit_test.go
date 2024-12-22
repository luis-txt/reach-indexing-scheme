package main

import (
	"math/rand"
	"testing"
)

func testIndexingSchemeFully(file string, t *testing.T) bool {
	t.Log("Reading graph...\n")
	g := ReadGraph(file)
	t.Log("Read graph with ", g.n, " nodes and ", g.m, " edges.")

	t.Log("Creating matrix using DFS for all vertex-pairs.")
	m1 := g.dfsCreateMatrix()
	t.Log("Created matrix successfully!")

	t.Log("Creating matrix using indexing scheme.\n")
	g, _, decomp, scheme := g.RunIndexingScheme()
	m2 := schemeToMatrix(scheme, decomp, g)
	t.Log("Created matrix successfully!")

	areEqual := compQuadraticMatrices(m1, m2)

	if areEqual {
		t.Log("Matrices are equal!")
	} else {
		t.Log("Matrices are not equal!")
	}
	return areEqual
}

func testIndexingSchemeRandomly(file string, nTests int, t *testing.T) bool {
	t.Log("Reading graph...")
	originalG := ReadGraph(file) // for DFS calculations
	g := ReadGraph(file)         // for scheme calculations

	g, _, decomp, scheme := g.RunIndexingScheme()

	t.Log("Testing scheme...")
	visited := make([]bool, originalG.n)
	stack := CreateStack[int](originalG.n)
	head := CreateStack[int](originalG.n)

	for i := 0; i < nTests; i++ {
		v := rand.Intn(originalG.n)
		w := rand.Intn(originalG.n)

		schemeAnswer := isReachable(v, w, scheme, decomp, g)
		dfsAnswer := originalG.runDFS(v, w, visited, stack, head) != nil

		if schemeAnswer != dfsAnswer {
			t.Log("False answer for: ", v, " and ", w, "!")
			t.Log("Scheme: ", schemeAnswer, ", DFS: ", dfsAnswer)
			return false
		}

		// Clean up DFS structures
		head.ClearStack()
		stack.ClearStack()
		for j := 0; j < originalG.n; j++ {
			visited[j] = false
		}
	}
	t.Log("All ", nTests, " tests successfully!\n")
	return true
}

func testDecomposition(g *Graph, decomp *Decomposition) bool {
	checked := make([]bool, g.n)

	for cNode := decomp.chains.first; cNode != nil; cNode = cNode.next {
		chain := cNode.val.entries
		i := 0
		for vNode := chain.first; vNode != nil; vNode = vNode.next {
			v := vNode.val
			if !checked[v] && decomp.vToChain[v].chain == cNode && decomp.vToChain[v].pos == i {
				checked[v] = true
			} else {
				return false
			}
			i++
		}
	}
	return true
}

// Tests

func TestDecomposition(t *testing.T) {
	files := []string{
		"./test_graphs/collapse.gr",
		"./data/real_world/twitter_combined.gr",
		"./data/gnm/gnm_1000_100000.gr",
		"./data/gnm/gnm_1000_10000.gr",
		"./data/gnm/gnm_1000_1000.gr",
		"./data/gnm/gnm_1000_100.gr",
		"./data/gnm/gnm_1000_10.gr",
		"./data/gnm/gnm_1000_1.gr",
	}

	for _, file := range files {
		t.Log("Testing ", file, "...")
		t.Run(file, func(t *testing.T) {
			g := ReadGraph(file)
			g = g.CollapseToDAG()
			topo := g.TopoSort()
			h3Decomp := g.HthreeConcat(topo)
			nodeDecomp := g.NodeOrderPathDecomp(topo)
			chainDecomp := g.ChainOrderPathDecomp(topo)

			if !testDecomposition(g, h3Decomp) {
				t.Errorf("H3-Concat decomposition test failed on graph: %s", file)
			}
			if !testDecomposition(g, nodeDecomp) {
				t.Errorf("Node-Order decomposition test failed on graph: %s", file)
			}
			if !testDecomposition(g, chainDecomp) {
				t.Errorf("Chain-Order decomposition test failed on graph: %s", file)
			}
			nodeDecomp.Concat(g)
			if !testDecomposition(g, nodeDecomp) {
				t.Errorf("Concat on Node-Order decomposition test failed on graph: %s", file)
			}
			chainDecomp.Concat(g)
			if !testDecomposition(g, chainDecomp) {
				t.Errorf("Concat on Node-Order decomposition test failed on graph: %s", file)
			}
		})
	}
}

func TestIndexingSchemeRandomly(t *testing.T) {
	files := []string{
		"./test_graphs/collapse.gr",
		"./data/real_world/twitter_combined.gr",
		"./data/gnm/gnm_1000_100000.gr",
		"./data/gnm/gnm_1000_10000.gr",
		"./data/gnm/gnm_1000_1000.gr",
		"./data/gnm/gnm_1000_100.gr",
		"./data/gnm/gnm_1000_10.gr",
		"./data/gnm/gnm_1000_1.gr",
	}

	for _, file := range files {
		t.Log("Testing ", file, "...")
		t.Run(file, func(t *testing.T) {
			if !testIndexingSchemeRandomly(file, 1000, t) {
				t.Errorf("Randomized indexing scheme failed on graph: %s", file)
			}
		})
	}
}

func TestIndexingSchemeFully(t *testing.T) {
	files := []string{
		"./test_graphs/collapse.gr",
		"./data/gnm/gnm_100_1000.gr",
		"./data/gnm/gnm_100_100.gr",
		"./data/gnm/gnm_100_10.gr",
		"./data/gnm/gnm_100_1.gr",
		"./data/gn/gn_100.gr",
	}

	for _, file := range files {
		t.Log("Testing ", file, "...")
		t.Run(file, func(t *testing.T) {
			if !testIndexingSchemeFully(file, t) {
				t.Errorf("Indexing scheme failed on graph: %s", file)
			}
		})
	}
}
