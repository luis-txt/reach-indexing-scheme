package main

import (
	"fmt"
	"math"
)

func getEntries(node *ListNode[Chain]) *LinkedList[int] {
	return &node.val.entries
}

func getFirst(node *ListNode[Chain]) *ListNode[int] {
	return getEntries(node).first
}

func getLast(node *ListNode[Chain]) *ListNode[int] {
	return getEntries(node).last
}

func (decomp *Decomposition) PrintChains() {
	fmt.Println("Total ", decomp.chains.n, " chains:")
	for cNode := decomp.chains.first; cNode != nil; cNode = cNode.next {
		c := cNode.val
		fmt.Print("C", c.id, ": ")
		c.entries.Print()
	}
}

func (decomp *Decomposition) PrintVtoC() {
	for _, cm := range decomp.vToChain {
		fmt.Print("{", cm.chain.val.id, ", ", cm.pos, "} ")
	}
	fmt.Println()
}

func createChain(id int) Chain {
	entries := *createLinkedList[int]()
	chain := Chain{id, entries}
	return chain
}

// Creates a new chain and adds the given vertex v to it.
// The created chain is then added to the given decomposition struct.
func addToNewChain(id, v int, decomp *Decomposition) {
	cNode := createListNode(createChain(id))
	getEntries(cNode).Add(createListNode(v))

	decomp.vToChain[v] = ChainMapping{cNode, 0}
	decomp.chains.Add(cNode)
}

// Traverses G starting from t in reversed direction (by traversing the incoming edges).
// Reuses the visited array to close in search area.
// Multiple runs are in O(|E| + l * (k_p - k_c)).
func reversedDFS(t int, g *Graph, decomp *Decomposition, visited []bool) int {
	stack := CreateStack[int](g.n)
	head := CreateStack[int](g.n)

	tChainNr := -1
	if decomp.vToChain[t].chain != nil {
		tChainNr = decomp.vToChain[t].chain.val.id
	}

	decompNodesProcessed++
	stack.Push(t)
	for !stack.IsEmpty() {
		v := stack.Peek()

		if !visited[v] {
			// New vertex discovered
			head.Push(v)
			visited[v] = true
			// Consider all incoming edges
			for e := g.nodes[v].in; e != nil; e = e.next {
				decompEdgesProcessed++
				sChain := decomp.vToChain[e.source].chain

				if sChain != nil && sChain.val.id != tChainNr && getLast(sChain).val == e.source {
					// s (e.source) is last in a chain which is different from t's chain.
					// Found chain with last vertex having a path to t.
					for !head.IsEmpty() {
						decompNodesProcessed++
						w := head.Pop()
						visited[w] = false
					}
					return e.source
				}
				if !visited[e.source] {
					decompNodesProcessed++
					stack.Push(e.source)
				}
			}
		} else if visited[v] && !head.IsEmpty() && head.Peek() == v {
			// Backtracking
			decompNodesProcessed++
			head.Pop()
			stack.Pop()
		} else {
			decompNodesProcessed++
			stack.Pop()
		}
	}
	return -1
}

// Finds a predecessor vertex that is the last one of its chain and has minimum outdegree
// of all predecessors that are last in their respective chain.
// Runs in O(N^-(v)).
func findLastOfChainMinOutdegPre(v int, g *Graph, visited []bool, decomp *Decomposition) int {
	w := -1
	minDeg := math.MaxInt

	for e := g.nodes[v].in; e != nil; e = e.next {
		decompEdgesProcessed++
		chain := decomp.vToChain[e.source].chain
		deg := g.nodes[e.source].outDeg

		if chain == nil || getLast(chain) == nil {
			continue
		}

		isLast := getLast(chain).val == e.source
		if !visited[e.source] && isLast && deg <= minDeg {
			minDeg = deg
			w = e.source
		}
	}
	return w
}

// Finds a sccessor vertex that has v as single source.
// Runs in O(N^+(v)).
func findSingleSourceSucc(v int, visited []bool, g *Graph) int {
	w := -1

	for e := g.nodes[v].out; e != nil; e = e.next {
		decompEdgesProcessed++
		if !visited[e.target] && g.nodes[e.target].inDeg == 1 {
			w = e.target
			break
		}
	}
	return w
}

// Appends the s-chain to the t-chain.
// Updates all fields in the decomp struct.
// Runs in O(l).
func combineTwoChains(sLinkNode, tLinkNode *ListNode[Chain], decomp *Decomposition) {
	offset := 0
	// Update chain mappings
	for vEntry := getFirst(sLinkNode); vEntry != nil; vEntry = vEntry.next {
		v := vEntry.val
		decomp.vToChain[v] = ChainMapping{
			tLinkNode, getEntries(tLinkNode).n + offset,
		}
		offset++
	}
	// Append s-chain to t-chain
	getEntries(tLinkNode).last.next = getFirst(sLinkNode)
	getEntries(sLinkNode).first.prev = getLast(tLinkNode)
	// Update t-chain n and last field
	getEntries(tLinkNode).n += offset
	getEntries(tLinkNode).last = getLast(sLinkNode)
	// Remove s-chain
	decomp.chains.Unlink(sLinkNode)
	// Update indices of chains
	i := 0
	for c := decomp.chains.first; c != nil; c = c.next {
		c.val.id = i
		i++
	}

}

// Creates a path decomposition using the Chain-Order heuristic.
// Runs in O(|V|+|E|).
func (g *Graph) ChainOrderPathDecomp(topo []int) *Decomposition {
	chains := createLinkedList[Chain]()
	vToChain := make([]ChainMapping, g.n)
	decomp := &Decomposition{vToChain, chains}
	used := make([]bool, g.n)
	id := 0

	for _, v := range topo {
		decompNodesProcessed++
		if !used[v] {
			// Create new chain for current vertex
			used[v] = true
			cNode := createListNode(createChain(id))

			getEntries(cNode).Add(createListNode(v))
			decompNodesProcessed++
			vToChain[v] = ChainMapping{cNode, 0}
			id++

			e := g.nodes[v].out
			for e != nil {
				decompEdgesProcessed++
				if !used[e.target] {
					// Add target of edge to current chain
					used[e.target] = true
					vToChain[e.target] = ChainMapping{cNode, getEntries(cNode).n}
					sNode := createListNode(e.target)
					getEntries(cNode).Add(sNode)
					decompNodesProcessed++
					e = g.nodes[e.target].out
				} else {
					e = e.next
				}
			}
			decomp.chains.Add(cNode)
		}
	}
	return decomp
}

// Creates a path decomposition using the Node-Order heuristic.
// Runs in O(|V|+|E|).
func (g *Graph) NodeOrderPathDecomp(topo []int) *Decomposition {
	// Init empty decomposition
	chains := createLinkedList[Chain]()
	vToChain := make([]ChainMapping, g.n)
	decomp := &Decomposition{vToChain, chains}
	id := 0
	for i := range vToChain {
		vToChain[i] = ChainMapping{nil, -1}
	}
	// Fill decomposition entries
	for _, v := range topo {
		decompNodesProcessed++
		used := false
		for e := g.nodes[v].in; e != nil; e = e.next {
			decompEdgesProcessed++
			sChainNode := vToChain[e.source].chain
			lastInSChain := getLast(sChainNode)
			if sChainNode != nil && lastInSChain != nil && lastInSChain.val == e.source {
				// Source is last of its chain
				// Insert v into chain of source
				vToChain[v] = ChainMapping{sChainNode, getEntries(sChainNode).n}
				tNode := createListNode(e.target)

				getEntries(sChainNode).Add(tNode)
				decompNodesProcessed++
				e = g.nodes[e.target].out
				used = true
				break
			}
		}
		if !used {
			// Create new chain
			addToNewChain(id, v, decomp)
			decompNodesProcessed++
			id++
		}
	}
	return decomp
}

// Tries to concatenate each chain in a given chain decomposition
// to another chain in the decomposition.
// Runs in O(|E| + l * (k_p - k_c)) by using the improved reversed DFS function
// and instant conatenation.
func (decomp *Decomposition) Concat(g *Graph) {
	visited := make([]bool, g.n)

	pNode := decomp.chains.first
	for pNode != nil {
		nextPNode := pNode.next
		// Consider each path (chain)
		path := getEntries(pNode)
		firstInTchain := path.first.val

		s := reversedDFS(firstInTchain, g, decomp, visited)

		if s != -1 {
			sChainNode := decomp.vToChain[s].chain
			combineTwoChains(pNode, sChainNode, decomp)
		}
		pNode = nextPNode
	}
}

// H3-Conc. heuristic that uses the improved reversed DFS function.
// Runs in O(|E| + l * (k_p - k_c)).
func (g *Graph) HthreeConcat(topo []int) *Decomposition {
	chains := createLinkedList[Chain]()
	vToChain := make([]ChainMapping, g.n)
	decomp := &Decomposition{vToChain, chains}
	visited := make([]bool, g.n)
	id := 0

	for i := 0; i < len(decomp.vToChain); i++ {
		decomp.vToChain[i] = ChainMapping{nil, -1}
	}

	for _, v := range topo {
		decompNodesProcessed++
		if decomp.vToChain[v].chain == nil {
			// v not assigned to a chain
			w := findLastOfChainMinOutdegPre(v, g, visited, decomp)
			if w == -1 {
				w = reversedDFS(v, g, decomp, visited)
			}
			if w != -1 {
				chainOfW := decomp.vToChain[w].chain
				vToChain[v] = ChainMapping{chainOfW, getEntries(chainOfW).n}
				vNode := createListNode(v)
				getEntries(chainOfW).Add(vNode)
				decompNodesProcessed++
			} else {
				// Create new chain
				addToNewChain(id, v, decomp)
				decompNodesProcessed++
				id++
			}
		}
		t := findSingleSourceSucc(v, visited, g)
		if t != -1 {
			// Found immediate successor with in-degree 1.
			chainOfV := decomp.vToChain[v].chain
			vToChain[t] = ChainMapping{chainOfV, getEntries(chainOfV).n}
			tNode := createListNode(t)
			getEntries(chainOfV).Add(tNode)
			decompNodesProcessed++
		}
	}
	return decomp
}
