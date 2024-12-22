package main

import (
	"math"
)

type ChainMapping struct {
	chain *ListNode[Chain]
	pos   int
}

type Chain struct {
	id      int
	entries LinkedList[int]
}

type Decomposition struct {
	vToChain []ChainMapping
	chains   *LinkedList[Chain]
}

type Edge struct {
	source  int
	target  int
	partner *Edge
	prev    *Edge
	next    *Edge
}

type Node struct {
	in     *Edge
	out    *Edge
	inDeg  int
	outDeg int
}

type IdMapping struct {
	vToId map[int]int
	idToV map[int]int
}

type Graph struct {
	n         int
	m         int
	nodes     []Node
	vToComp   []int
	idMapping IdMapping
}

type SccData struct {
	u       int
	t       int
	lowLink []int
	pre     []int
	vToComp []int
	time    []int
	compToV [][]int
	onStack []bool
}

type Stacks struct {
	recStack *Stack[int]
	head     *Stack[int]
	stack    *Stack[int]
}

type NodeCollitions struct {
	reached []*Edge
	changed []int
}

func CreateGraph(n int) *Graph {
	nodes := make([]Node, n)
	return &Graph{n, 0, nodes, nil, IdMapping{}}
}

func (g *Graph) PrintGraph() {
	logger.Println("=== G: (", g.n, ",", g.m, ")")
	for v := 0; v < g.n; v++ {
		logger.Print(v, " -> ")

		for e := g.nodes[v].out; e != nil; e = e.next {
			logger.Print(e.target, " ")
		}
		logger.Println()
	}
	if g.vToComp != nil {
		logger.Println("=== V -> Comp:")
		logger.Println(g.vToComp)
	}
	if g.idMapping.vToId != nil {
		logger.Println("=== ID -> V:")
		logger.Println(g.idMapping.idToV)
	}
}

func pushEdgeToList(e *Edge, list *Edge) *Edge {
	if list == nil {
		list = e
	} else {
		list.prev = e
		e.next = list
		list = e
	}
	return list
}

func (g *Graph) AddEdge(e *Edge) {
	g.m++
	if e.partner == nil {
		partnerEdge := Edge{e.source, e.target, e, nil, nil}
		e.partner = &partnerEdge
	}
	g.nodes[e.source].out = pushEdgeToList(e, g.nodes[e.source].out)
	g.nodes[e.target].in = pushEdgeToList(e.partner, g.nodes[e.target].in)
	// Update degrees
	g.nodes[e.source].outDeg++
	g.nodes[e.target].inDeg++
}

// Sets prev and next to nil but keeps partner.
// Unlinks partner edge from its list but keeps its partner.
// Updates degree counter of graph g.
func unlink(g *Graph, e *Edge, isOut bool) {
	if !isOut {
		e = e.partner
	}
	// Unlink e
	if e.prev == nil {
		g.nodes[e.source].out = e.next
	} else {
		e.prev.next = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	}
	e.next = nil
	e.prev = nil
	// Unlink partner
	if e.partner.prev == nil {
		g.nodes[e.target].in = e.partner.next
	} else {
		e.partner.prev.next = e.partner.next
	}
	if e.partner.next != nil {
		e.partner.next.prev = e.partner.prev
	}
	e.partner.next = nil
	e.partner.prev = nil
	// Update degrees
	g.m--
	g.nodes[e.source].outDeg--
	g.nodes[e.target].inDeg--
}

func (g *Graph) AddVtoIdMapping(id int, nextIndex int) int {
	_, contained := g.idMapping.idToV[id]
	if contained {
		return nextIndex
	} else {
		g.idMapping.idToV[id] = nextIndex
		g.idMapping.vToId[nextIndex] = id
	}
	nextIndex++
	return nextIndex
}

// Runs a modified DFS starting from s for topological sorting.
func topoDFS(s, i int, g *Graph, visited []bool, topoOrder []int, stack, head *Stack[int]) int {
	stack.Push(s)

	for !stack.IsEmpty() {
		v := stack.Peek()

		if !visited[v] {
			visited[v] = true
			head.Push(v)

			// Consider all outgoing edges
			for e := g.nodes[v].out; e != nil; e = e.next {
				if !visited[e.target] {
					stack.Push(e.target)
				}
			}
		} else if head.Peek() == v {
			// Backtracking: Next node detected
			topoOrder[i] = v
			i--
			head.Pop()
			stack.Pop()
		} else {
			stack.Pop()
		}
	}
	return i
}

// Topologically sorts the given DAG in O(|V| + |E|).
func (g *Graph) TopoSort() []int {
	visited := make([]bool, g.n)
	topoOrder := make([]int, g.n)
	stack := CreateStack[int](g.n)
	head := CreateStack[int](g.n)
	i := g.n - 1

	// Create topologgerical order
	for v := 0; v < g.n; v++ {
		if g.nodes[v].inDeg == 0 && !visited[v] {
			i = topoDFS(v, i, g, visited, topoOrder, stack, head)
			// Clean up stacks
			head.ClearStack()
			stack.ClearStack()
		}
	}
	return topoOrder
}

// Runs DFS version for strongly connected components on a part of the Graph g.
// Returns the current time value and the calculated array
// that maps components to vertices: compToV.
func sccDFS(data SccData, g *Graph, stacks Stacks) (int, [][]int) {
	stacks.recStack.Push(data.u)
	collapseNodesProcessed++
	for !stacks.recStack.IsEmpty() {
		v := stacks.recStack.Peek()

		if data.time[v] == math.MaxInt {
			// First discovery of this vertex
			data.time[v] = data.t
			data.lowLink[v] = data.t
			data.onStack[v] = true
			stacks.stack.Push(v)
			stacks.head.Push(v)
			collapseNodesProcessed++
			data.t++
			// Consider all outgoing edges
			for e := g.nodes[v].out; e != nil; e = e.next {
				collapseEdgesProcessed++

				if data.time[e.target] == math.MaxInt {
					stacks.recStack.Push(e.target)
					data.pre[e.target] = v
				} else if data.onStack[e.target] {
					data.lowLink[v] = min(data.lowLink[v], data.time[e.target])
				}
			}
		} else if stacks.head.Peek() == v {
			// Backtracking
			for e := g.nodes[v].out; e != nil; e = e.next {
				collapseEdgesProcessed++

				if data.pre[e.target] == v {
					// "recursively-called" on this vertex
					data.lowLink[v] = min(data.lowLink[v], data.lowLink[e.target])
				}
			}
			if data.lowLink[v] == data.time[v] {
				comp := make([]int, 0, g.n)
				for {
					w := stacks.stack.Pop()
					collapseNodesProcessed++
					data.onStack[w] = false
					data.vToComp[w] = len(data.compToV)
					comp = append(comp, w)
					if stacks.stack.IsEmpty() || v == w {
						break
					}
				}
				data.compToV = append(data.compToV, comp)
			}
			stacks.head.Pop()
			stacks.recStack.Pop()
			collapseNodesProcessed++
		} else {
			stacks.recStack.Pop()
			collapseNodesProcessed++
		}
	}
	return data.t, data.compToV
}

// Tarjan's strongly connected components algorithm.
// Returns an array that maps the vertices to components (vToComp)
// and another array that maps the components to vertices.
func (g *Graph) FindSCCs() ([]int, [][]int) {
	lowLink := make([]int, g.n)
	time := make([]int, g.n)
	vToComp := make([]int, g.n)
	compToV := make([][]int, 0, g.n)
	stacks := Stacks{
		CreateStack[int](g.n), CreateStack[int](g.n), CreateStack[int](g.n),
	}

	pre := make([]int, g.n)
	onStack := make([]bool, g.n)
	t := 0

	for i := 0; i < g.n; i++ {
		collapseNodesProcessed++
		lowLink[i] = math.MaxInt
		time[i] = math.MaxInt
	}
	// Consider all nodes
	for v := 0; v < g.n; v++ {
		collapseNodesProcessed++

		if time[v] == math.MaxInt {
			data := SccData{
				v, t, lowLink, pre, vToComp, time, compToV, onStack,
			}
			t, compToV = sccDFS(data, g, stacks)
			// Clean up stacks
			stacks.head.ClearStack()
			stacks.stack.ClearStack()
			stacks.recStack.ClearStack()
		}
	}
	return vToComp, compToV
}

// Collapses the given graph to its strongly connected components.
// Uses Tarjan's Strongly-connected-components algorithm.
// Runs in O(|V|+|E|).
func (g *Graph) CollapseToDAG() *Graph {
	vToComp, compToV := g.FindSCCs()
	gPrime := CreateGraph(len(compToV))
	gPrime.vToComp = vToComp
	gPrime.idMapping = g.idMapping
	collision := make([]bool, len(compToV))
	changed := make([]int, len(compToV)) // for resetting changed values for later nodes

	for compNr := 0; compNr < len(compToV); compNr++ {
		i := 0
		comp := compToV[compNr]
		for _, v := range comp {
			collapseNodesProcessed++
			for e := g.nodes[v].out; e != nil; e = e.next {
				collapseEdgesProcessed++

				tCompNr := vToComp[e.target]
				if compNr != tCompNr && !collision[tCompNr] {
					changed[i] = tCompNr
					i++
					ePrime := Edge{compNr, tCompNr, nil, nil, nil}
					gPrime.AddEdge(&ePrime)
					collision[tCompNr] = true
				}
			}
		}
		// Clear reached-indices for the next vertex
		for i > 0 {
			collision[changed[i-1]] = false
			i--
		}
	}
	return gPrime
}

// Removes incoming or outgoing edges for chains, depending on isVtoC flag in O(|V|+|E|).
func (g *Graph) removeOneSidedTransitiveEdges(collitions NodeCollitions, decomp *Decomposition, isVtoC bool) {
	for v := 0; v < g.n; v++ {
		i := 0
		e := g.nodes[v].in // Use incoming edge-list (for vertex to chain removal)
		if isVtoC {
			e = g.nodes[v].out // Use outgoing edge-list (for vertex to chain removal)
		}
		for e != nil {
			nextE := e.next
			w := e.source // Use chain of source node of edge (for vertex to chain removal)
			if isVtoC {
				w = e.target // Use chain of target node of edge (for chain to vertex removal)
			}
			wChain := decomp.vToChain[w].chain.val
			if collitions.reached[wChain.id] == nil {
				collitions.reached[wChain.id] = e
				collitions.changed[i] = wChain.id
				i++
			} else {
				// Handle already reached chain of current target
				oldEdge := collitions.reached[wChain.id]
				numRemovedTransitveEdges++

				newPosT := decomp.vToChain[e.target].pos
				oldPosT := decomp.vToChain[oldEdge.target].pos
				newPosS := decomp.vToChain[e.source].pos
				oldPosS := decomp.vToChain[oldEdge.source].pos
				if (newPosT > oldPosT && isVtoC) || (newPosS < oldPosS && !isVtoC) {
					unlink(g, e, isVtoC)
				} else {
					unlink(g, oldEdge, isVtoC)
					collitions.reached[wChain.id] = e
				}
			}
			e = nextE
		}
		// Clear reached-indices for the next vertex
		for i > 0 {
			collitions.reached[collitions.changed[i-1]] = nil
			i--
		}
	}
}

// Heuristic to remove transitive edges in O(|V| + |E|).
func (g *Graph) RemoveTransitiveEdges(decomp *Decomposition) {
	reached := make([]*Edge, decomp.chains.n)
	changed := make([]int, decomp.chains.n) // for resetting changed values for later nodes
	collitions := NodeCollitions{reached, changed}
	// Remove transitive edges from vertex to chain
	g.removeOneSidedTransitiveEdges(collitions, decomp, true)
	// Remove transitive edges from chain to vertex
	g.removeOneSidedTransitiveEdges(collitions, decomp, false)
}

// Topologically sorts the outgoing edges of the given graph in O(|V| + |E|).
func (g *Graph) TopoSortOutEdges(topoOrder []int) {
	edgeList := make([]*Edge, g.n)
	// Fill new adjacency lists in topologgerical order
	for i := len(topoOrder) - 1; i >= 0; i-- {
		v := topoOrder[i]

		for e := g.nodes[v].in; e != nil; e = e.next {
			e.partner.next = nil
			e.partner.prev = nil
			edgeList[e.source] = pushEdgeToList(e.partner, edgeList[e.source])
		}
	}
	// Replace old outgoing adjacency lists
	for v := 0; v < g.n; v++ {
		g.nodes[v].out = edgeList[v]
	}
}
