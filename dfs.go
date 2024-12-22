package main

// Runs DFS only from the given start vertex s and check whether it can reach t.
// Runs in O(|E|).
func (g *Graph) runDFS(s, t int, visited []bool, stack, head *Stack[int]) []int {
	stack.Push(s)
	for !stack.IsEmpty() {
		u := stack.Peek()

		if !visited[u] {
			visited[u] = true
			head.Push(u)
			if u == t {
				return head.data
			}
			for e := g.nodes[u].out; e != nil; e = e.next {
				if !visited[e.target] {
					stack.Push(e.target)
				}
			}
		} else if head.Peek() == u {
			// Backtracking
			head.Pop()
			stack.Pop()
		} else {
			stack.Pop()
		}
	}
	return nil
}

// Wrapper for DFS
func (g *Graph) DFS(s, t int) []int {
	visited := make([]bool, g.n)
	stack := CreateStack[int](g.n)
	head := CreateStack[int](g.n)

	return g.runDFS(s, t, visited, stack, head)
}

// Creates a matrix by running DFS for each pair
// of vertices.
// Runs in O(|V|^2 * (|V| + |E|)).
func (g *Graph) dfsCreateMatrix() [][]bool {
	visited := make([]bool, g.n)
	stack := CreateStack[int](g.n)
	head := CreateStack[int](g.n)

	// Create n*n matrix
	matrix := make([][]bool, g.n)
	for i := 0; i < g.n; i++ {
		matrix[i] = make([]bool, g.n)
	}
	// Fill matrix
	for v := 0; v < g.n; v++ {
		for w := 0; w < g.n; w++ {
			matrix[v][w] = g.runDFS(v, w, visited, stack, head) != nil
			// Clean up DFS structures
			head.ClearStack()
			stack.ClearStack()
			for j := 0; j < g.n; j++ {
				visited[j] = false
			}
		}
	}
	return matrix
}
