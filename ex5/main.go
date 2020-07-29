package main

import (
	"bufio"
	"fmt"
	"os"
)

const NIL = 0          //used in Hopcroft-karp algorithm if order for it to be easier to read
const inf = ^uint16(0) //maximum value of uint16

/*
The method use for solving this problem is the following:
-we create a map of all the points in the rectangles (we can do this as there is less than 10 000 points so we won't
run out of memory; storing all the 1 000 000 rectangles could!)
-we convert it to a list of all the possible vertical and horizontal lines and the points they each go through
-we then turn it into a graphs (one for each contiguous block) in which there is an edge from a line to a point only
if this lines goes through this point
-we convert this graph to another one in which two lines are linked only if they share a common point
-we convert it to a bipartite graph G(u, v) as the previous graph not being bipartite would imply a diagonal line. As
this graph is bipartite, Kőnig's theorem tells us that the number of edges in the maximum cardinality matching is the
same as the number of vertices in the minimum vertex cover. And the minimum vertex cover in our previous graph is the
minimal number of nodes we have to keep for each point in our problem to be part of a vertical/horizontal line! The
vertex cover is a NP-Hard problem so it's hard to solve, but the maximum cardinality matching can be obtained with
the Hopcroft–Karp algorithm which is pretty fast.
-we then sum the lengths of the maximum cardinality matching of each graph to get the minimum number of lines required
 */


func main() {
	pI := getAndParseInput()
	fmt.Println(solve(pI))
}

func getAndParseInput() []bipartite {
	pointsMap := make(map[Point]uint16) //this map will store every points which exists and its ID
	// first 14 bits are for the ID; bit of weight 2 is whether or not an horizontal line going through this points has
	//been found; bit of weight 1 is the same with a vertical line
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan() //ignore map width
	scanner.Scan() //ignore map height

	var shipCount int
	scanner.Scan()
	_, err := fmt.Sscanf(scanner.Text(), "%d", &shipCount)
	if err != nil {
		panic(err)
	}

	currentId := uint16(0) //this will sore the ID that ne next point will have
	for i := 0; i < shipCount; i++ {
		scanner.Scan()
		str := scanner.Text()

		var x, y, u, v int
		if _, err = fmt.Sscanf(str, "%d %d %d %d", &x, &y, &u, &v); err != nil {
			panic(err)
		}

		for x0 := x; x0 < u; x0++ { //for each point in the rectangle
			for y0 := y; y0 < v; y0++ {
				pt := Point{
					X: x0,
					Y: y0,
				}
				if _, exists := pointsMap[pt]; !exists { // if it isn't yet in the map
					pointsMap[pt] = currentId << 2 //we add it with it's ID in the 14 bits of higher weight
					currentId++                    //and we increment the currentId so the next point has another ID
				}
			}
		}
	}

	//these functions are inline so they share context with getAndParseInput
	lineVerticaly := func(p Point) []uint16 { //this function creates a vertical line and marks every point in the map
		// bit setting the second bit to 1
		currentList := make([]uint16, 0) // we store all the points we have found in this line
		pCopy := p
		exists := true
		var value uint16
		for ; ; pCopy.Y-- { //we start by decrementing Y
			value, exists = pointsMap[pCopy]
			if ! exists { //if the point doesn't exist, we reached the end of the line
				break
			}
			pointsMap[pCopy] |= 2                       //otherwise we set the bit to 1
			currentList = append(currentList, value>>2) //and add the id of the point to the list of explored ones
			//we shift by two to ignore the bits 1 and 2 which are not part of the ID
		}
		pCopy = p //then we reset the point
		exists = true
		for pCopy.Y++; ; pCopy.Y++ { //and increment Y
			value, exists = pointsMap[pCopy]
			if ! exists {
				break
			}
			pointsMap[pCopy] |= 2 //we then do the same
			currentList = append(currentList, value>>2)
		}
		return currentList //and return the list of explored points
	}

	//see lineVerticaly for explanation
	lineHorizontaly := func(p Point) []uint16 {
		currentList := make([]uint16, 0)
		pCopy := p
		exists := true
		var value uint16
		for ; ; pCopy.X-- {
			value, exists = pointsMap[pCopy]
			if ! exists {
				break
			}
			pointsMap[pCopy] |= 1 //this time we set bit 1
			currentList = append(currentList, value>>2)
		}
		pCopy = p
		exists = true
		for pCopy.X++; ; pCopy.X++ {
			value, exists = pointsMap[pCopy]
			if ! exists {
				break
			}
			pointsMap[pCopy] |= 1
			currentList = append(currentList, value>>2)
		}
		return currentList
	}

	blocks := make([][][]uint16, 0) //block[n] is the nth block of points
	//block[n0][n1] is the list of the points crossed by the n1th line of block n

	//this loop split the map of points in blocks.
	//if there is no path from a point to another, they are in a different block
	//otherwise they are in the same block
	id := 0
	for key, value := range pointsMap { //for each point
		if value&3 == 0 { //if it is unexplored
			blocks = append(blocks, make([][]uint16, 0))        //we create a new block
			blocks[id] = append(blocks[id], lineVerticaly(key)) //and draw our first lines
			blocks[id] = append(blocks[id], lineHorizontaly(key))
		} else {
			continue
		}
		found := true
		for found { //while we can find a point which has less than two lines going though it
			found = false
			for key, value := range pointsMap { //we seek unexplored points
				if value&3 == 0 || value&3 == 3 {
					continue
				}
				if value&1 == 0 { //and line from it
					found = true
					blocks[id] = append(blocks[id], lineHorizontaly(key))
				} else {
					found = true
					blocks[id] = append(blocks[id], lineVerticaly(key))
				}
			}
		}

		id++
	}

	pointsMap = nil

	//in this part, we convert each block to a graph
	// a node is linked to another one if and only if there is a single line going through both of them
	graphs := make([][][]uint16, len(blocks))
	for i0, block := range blocks { //for each block
		for range block { // we initialise the graph with empty lists
			graphs[i0] = append(graphs[i0], make([]uint16, 0))
		}

		//this is likely a complicated way of doing what I want but it works
		//this is the part converting the block to a graph
		for i1 := 0; i1 < len(block); i1++ {
			for i2 := 0; i2 < len(block[i1]); i2++ { //for each point in the lines
			searchLoop0:
				for i3 := i1 + 1; i3 < len(block); i3++ {
					for i4 := 0; i4 < len(block[i3]); i4++ {
						if block[i3][i4] == block[i1][i2] { //if it is connected to another one
							graphs[i0][i1] = append(graphs[i0][i1], uint16(i3)) //we add them to the graph
							graphs[i0][i3] = append(graphs[i0][i3], uint16(i1))
							break searchLoop0 //and break
						}
					}
				}
			}
		}
	}

	//we use a modified Breadth First Search to convert our graph to a bipartite one
	//the utility nof this is explained in the header
	makeBipartite := func(tree [][]uint16) bipartite {
		groups := make([]uint8, len(tree)) //we create three groups :
		//0 -> unexplored
		//1 -> group 1
		//2 -> group 2
		groups[0] = 1
		queue := make([]uint16, 1) //we put our first node in the queue
		queue [0] = 0
		for len(queue) > 0 { //as long as there is a point in the queue
			node := queue[0]
			queue[0] = queue[len(queue)-1]
			queue = queue[:len(queue)-1]           //we pop a node from the queue
			for _, neighbour := range tree[node] { // for each node neighbour of it
				if groups[neighbour] == 0 { // if it is unexplored
					queue = append(queue, neighbour) //we add it to the queue
					if groups[node] == 1 { //if our node is in group 1
						groups[neighbour] = 2 //it's neighbour goes to group 2
					} else {
						groups[neighbour] = 1 //otherwise it goes to group 1
					}
				}
			}
		}

		//from this list of groups, we convert the graph to the bipartite one
		adj := make([][]uint16, 1) //this list will store the points adjacent to u[n]
		//adj[0] is empty as it is the NIL point, described in the header

		//as we are building this new tree, the IDs will change
		idUToTree := make([]uint16, 0) //this list and the map below are here for the conversion
		idTreeToV := make(map[uint16]uint16)
		u := make([]uint16, 0) //we initialise empty lists
		v := make([]uint16, 0)

		for treeId, group := range groups { //for each point
			if group == 1 { //if it's group is 1 we add it to u
				index := uint16(len(u) + 1)
				u = append(u, index)
				idUToTree = append(idUToTree, uint16(treeId)) //and add it's original ID to the conversion list
			} else { //otherwise we add it to v
				index := uint16(len(v) + 1)
				v = append(v, index)
				idTreeToV[uint16(treeId)] = index //and add it's ID to the conversion map
			}
		}

		for idU, idTree := range idUToTree { //then, from the conversion lists
			adj = append(adj, make([]uint16, 0)) //we create and fill adj
			for _, neighbour := range tree[idTree] {
				adj[idU+1] = append(adj[idU+1], idTreeToV[neighbour])
			}
		}

		return bipartite{ //and we return
			u:   u,
			v:   v,
			adj: adj,
		}
	}

	//this array will store all the graphs
	bipart := make([]bipartite, len(graphs))

	//for each graph, we convert it to a bipartite one and store it
	for index, tree := range graphs {
		bipart[index] = makeBipartite(tree)
		graphs[index] = nil
	}

	return bipart
}

type Point struct {
	X int
	Y int
}

//this structure represents a bipartite graph
type bipartite struct {
	u   []uint16   //list of the points in the first set of the graph
	v   []uint16   //list of the points in the second one
	adj [][]uint16 //adj[u[n]] stores every points of v to which u[n] is connected
}

//this is and implementation of the Depth First Search algorithm to explore our graphs for the Hopcroft-Karp algorithm
//this one isn't inline as it is recursive
func DFS(u uint16, bipart bipartite, dist []uint16, Pair_U []uint16, Pair_V []uint16) bool {
	if u != NIL {
		for _, v := range bipart.adj[u] {
			if dist[Pair_V[v]] == dist[u]+1 {
				if DFS(Pair_V[v], bipart, dist, Pair_U, Pair_V) {
					Pair_V[v] = u
					Pair_U[u] = v
					return true
				}
			}
		}
		dist[u] = inf
		return false
	}
	return true
}

//this function returns the length of the maximum cardinality matching of our graph
//more details can be found at https://en.wikipedia.org/wiki/Hopcroft%E2%80%93Karp_algorithm#Pseudocode
func HopcroftKarp(bipart bipartite) int {
	Pair_U := make([]uint16, len(bipart.u)+1)
	for _, u := range bipart.u {
		Pair_U[u] = NIL
	}
	Pair_V := make([]uint16, len(bipart.v)+1)
	for _, v := range bipart.v {
		Pair_V[v] = NIL
	}
	matching := 0
	dist := make([]uint16, len(bipart.u)+1)

	//this is the Breadth First Search algorithm
	BFS := func() bool {
		queue := make([]uint16, 0)
		for _, u := range bipart.u {
			if Pair_U[u] == NIL {
				dist[u] = 0
				queue = append(queue, u)
			} else {
				dist[u] = inf
			}
		}
		dist[NIL] = inf
		for len(queue) > 0 {
			u := queue[0]
			queue[0] = queue[len(queue)-1]
			queue = queue[:len(queue)-1]
			if dist[u] < dist[NIL] {
				for _, v := range bipart.adj[u] {
					if dist[Pair_V[v]] == inf {
						dist[Pair_V[v]] = dist[u] + 1
						queue = append(queue, Pair_V[v])
					}
				}
			}
		}
		return dist[NIL] != inf
	}

	for BFS() {
		for _, u := range bipart.u {
			if Pair_U[u] == NIL {
				if DFS(u, bipart, dist, Pair_U, Pair_V) {
					matching++
				}
			}
		}
	}

	return matching
}

func solve(input []bipartite) int {
	//we sum the number of lines required for each bipartite graph
	sum := 0
	for _, bipart := range input {
		curr := HopcroftKarp(bipart)
		sum += curr
	}
	return sum
}
