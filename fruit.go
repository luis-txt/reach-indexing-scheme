package main

import (
	"flag"
	"fmt"
	"time"
)

var logger *Log
var numRemovedTransitveEdges uint
var schemeNodesProcessed uint
var schemeEdgesProcessed uint
var decompNodesProcessed uint
var decompEdgesProcessed uint
var collapseNodesProcessed uint
var collapseEdgesProcessed uint
var verboseFlag bool
var benchFlag bool
var matrixFlag bool
var nodeOrderFlag bool
var chainOrderFlag bool
var nodeConcFlag bool
var chainConcFlag bool

func init() {
	flag.BoolVar(&verboseFlag, "v", false,
		"Enable verbose output.",
	)
	flag.BoolVar(&matrixFlag, "m", false,
		"Create transitive closure matrix.",
	)
	flag.BoolVar(&benchFlag, "b", false,
		"If set, returns informations from computation and time.",
	)
	flag.BoolVar(&nodeOrderFlag, "no", false,
		"Sets the chain decomposition algorithm to the node-order heuristic.",
	)
	flag.BoolVar(&chainOrderFlag, "co", false,
		"Sets the chain decomposition algorithm to the chain-order heuristic.",
	)
	flag.BoolVar(&nodeConcFlag, "noc", false,
		"Sets the chain decomposition algorithm to the node-order heuristic with path concatenation.",
	)
	flag.BoolVar(&chainConcFlag, "coc", false,
		"Sets the chain decomposition algorithm to the chain-order heuristic with path concatenation.",
	)
	logger = &Log{false}
}

func decomposeAccordingToFlag(g *Graph, topo []int) *Decomposition {
	var decomp *Decomposition
	if nodeOrderFlag {
		decomp = g.NodeOrderPathDecomp(topo)
	} else if chainOrderFlag {
		decomp = g.ChainOrderPathDecomp(topo)
	} else if nodeConcFlag {
		decomp = g.NodeOrderPathDecomp(topo)
		decomp.Concat(g)
	} else if chainConcFlag {
		decomp = g.ChainOrderPathDecomp(topo)
		decomp.Concat(g)
	} else {
		decomp = g.HthreeConcat(topo)
	}
	return decomp
}

func (g *Graph) RunIndexingScheme() (*Graph, []int, *Decomposition, [][]int) {
	oldM := g.m

	logger.Println("Collapsing the graph to a DAG...")
	g = g.CollapseToDAG()
	logger.Print("Collapsed graph to DAG with ", g.n, " components.\n\n")

	logger.Println("Topologically sorting the DAG...")
	topo := g.TopoSort()
	logger.Print("Sorted successfully.\n\n")

	logger.Println("Decomposing the DAG into chains...")
	decomp := decomposeAccordingToFlag(g, topo)
	logger.Print("Decomposed DAG into ", decomp.chains.n, " chains.\n\n")

	logger.Println("Removing some transitive edges...")
	g.RemoveTransitiveEdges(decomp)
	logger.Print("Reduced number of edges from ", oldM, " to ", g.m, ".\n\n")

	logger.Println("Sorting adjacency lists in topologgerical order...")
	g.TopoSortOutEdges(topo)
	logger.Print("Sorted successfully.\n\n")

	logger.Println("Creating indexing scheme...")
	scheme := g.CreateIndexingScheme(topo, decomp)
	logger.Print("Successfully creating indexing scheme.\n\n")

	return g, topo, decomp, scheme
}

func (g *Graph) BenchIndexingScheme(readingTime float64, totalStart time.Time) (*Graph, []int, *Decomposition, [][]int) {
	oldN := g.n
	oldM := g.m

	compStart := time.Now()

	preprocessStart := time.Now()
	g = g.CollapseToDAG()
	collapseTime := getTimeMS(preprocessStart)

	topoStart := time.Now()
	topo := g.TopoSort()
	topoTime := getTimeMS(topoStart)

	preprocessTime := getTimeMS(preprocessStart)

	decompStart := time.Now()
	decomp := decomposeAccordingToFlag(g, topo)
	decompTime := getTimeMS(decompStart)

	preprocessStart = time.Now()
	g.RemoveTransitiveEdges(decomp)
	removeEdgesTime := getTimeMS(preprocessStart)

	topoEdgesStart := time.Now()
	g.TopoSortOutEdges(topo)
	topoEdgesTime := getTimeMS(topoEdgesStart)

	preprocessTime += getTimeMS(preprocessStart)

	schemeStart := time.Now()
	scheme := g.CreateIndexingScheme(topo, decomp)
	schemeTime := getTimeMS(schemeStart)

	totalTime := getTimeMS(totalStart)
	compTime := getTimeMS(compStart)

	fmt.Println(
		"#nodes: ", oldN, ", #edges: ", oldM,
		", #scc: ", g.n,
		", #chains: ", decomp.chains.n,
		", scheme-size: ", g.n*decomp.chains.n,
		", #removed-edges: ", numRemovedTransitveEdges,
		", #collapse-nodes: ", collapseNodesProcessed,
		", #collapse-edges: ", collapseEdgesProcessed,
		", #decomp-nodes: ", decompNodesProcessed,
		", #decomp-edges: ", decompEdgesProcessed,
		", #scheme-nodes: ", schemeNodesProcessed,
		", #scheme-edges: ", schemeEdgesProcessed,
		", time-decomp: ", fmt.Sprintf("%.4f ms", decompTime),
		", time-preprocess: ", fmt.Sprintf("%.4f ms", preprocessTime),
		", time-scheme: ", fmt.Sprintf("%.4f ms", schemeTime),
		", time-reading: ", fmt.Sprintf("%.4f ms", readingTime),
		", time-comp: ", fmt.Sprintf("%.4f ms", compTime),
		", time-total: ", fmt.Sprintf("%.4f ms", totalTime),
		", time-collapse: ", fmt.Sprintf("%.4f ms", collapseTime),
		", time-topo: ", fmt.Sprintf("%.4f ms", topoTime),
		", time-remove_edges: ", fmt.Sprintf("%.4f ms", removeEdgesTime),
		", time-topo_edges_time: ", fmt.Sprintf("%.4f ms", topoEdgesTime),
	)
	return g, topo, decomp, scheme
}

func main() {
	flag.Parse()
	args := flag.Args()

	logger.verbose = verboseFlag

	if len(args) < 1 {
		fmt.Println("Usage: go run fruit [-v or -m or -b] [-no or -noc or -co or -coc] <file_path>")
		return
	}
	file := args[0]

	totalStart := time.Now()
	start := time.Now()
	g := ReadGraph(file)

	var scheme [][]int
	var decomp *Decomposition

	if benchFlag {
		readingStart := float64(time.Since(start).Nanoseconds()) / 1e6
		g, _, decomp, scheme = g.BenchIndexingScheme(readingStart, totalStart)
	} else {
		logger.Print("Read Graph (|V|=", g.n, ", |E|=", g.m, ").\n\n")
		g, _, decomp, scheme = g.RunIndexingScheme()
	}

	if matrixFlag {
		matrix := schemeToMatrix(scheme, decomp, g)
		printMatrix(matrix)
	}
}
