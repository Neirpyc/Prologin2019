package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const maxOfInt = int((^uint(0)) >> 1) //we must do it in a weird way because ints can be negatives so
// we shift it by 1

func main() {
	input := getAndParseInput()
	fmt.Println(solve(input))
}

func getAndParseInput() (inp input) {
	scanner := bufio.NewScanner(os.Stdin)

	inp = input{
		Map:     Map{},
		mission: mission{},
	}

	scanner.Scan()
	_, err := fmt.Sscanf(scanner.Text(), "%d", &inp.Map.typeCount)
	if err != nil {
		panic(err)
	}

	scanner.Scan()
	_, err = fmt.Sscanf(scanner.Text(), "%d", &inp.Map.planetCount)
	if err != nil {
		panic(err)
	}

	inp.Map.planets = make([][]position, inp.Map.typeCount)
	for i := 0; i < inp.Map.planetCount; i++ {
		scanner.Scan()
		str := scanner.Text()

		var x, y, k int

		_, err = fmt.Sscanf(str, "%d %d %d", &x, &y, &k)

		inp.Map.planets[k] = append(inp.Map.planets[k], position{X: x, Y: y})

		if err != nil {
			panic(err)
		}
	}

	scanner.Scan()
	_, err = fmt.Sscanf(scanner.Text(), "%d", &inp.mission.duration)
	if err != nil {
		panic(err)
	}

	inp.mission.stages = make([]uint16, inp.mission.duration)
	scanner.Scan()
	planetTypes := strings.Split(scanner.Text(), " ")
	for i := 0; i < inp.mission.duration; i++ {
		_, err = fmt.Sscanf(planetTypes[i], "%d",
			&inp.mission.stages[i])
		if err != nil {
			panic(err)
		}
	}

	return inp
}

//The algorithm used here is slightly modified djikstra
//we have a list of possible planets to start from,and a list of planets we could end to for each move
//from each beginNode, we compute the distance to each endNode and store in each endNode the shortest distance found
//leading to it
//when we have done it for every step, we return the lowest of the values of the endNodes
func solve(input input) int {
	beginNodes := make([]*node, 0)
	endNodes := input.Map.getNodesOfType(input.mission.stages[0], 0) //initialise nodes with distance 0
	for i := 1; i < input.mission.duration; i++ {
		beginNodes = endNodes                                                  //we use the endNodes as beginNodes
		endNodes = input.Map.getNodesOfType(input.mission.stages[i], maxOfInt) //we set the distance of each enNode
		// to the maximum possible
		for _, beginNode := range beginNodes { //for each
			for _, endNode := range endNodes { //combination of nodes
				if summedDistance := beginNode.bestDistanceYet + beginNode.position.distanceTo(endNode.position);
					summedDistance < endNode.bestDistanceYet {
					endNode.bestDistanceYet = summedDistance
				}
			}
		}
	}

	//we are done, we just have to find the lowest of the endNodes' distances
	lowest := maxOfInt
	for _, node := range endNodes {
		if node.bestDistanceYet < lowest {
			lowest = node.bestDistanceYet
		}
	}
	return lowest
}

type input struct {
	Map     Map
	mission mission
}

type mission struct {
	duration int
	stages   []uint16 //uint16 here to spare 1kB of memory -> this is ridiculous but fun
}

type Map struct {
	typeCount   int
	planetCount int
	planets     [][]position //is a 2D-list working as a map:
	// planets[k] returns the positions of all the planets with type k
}

//returns the list of all the nodes of type k and initialises their distances to defaultDistance
func (m Map) getNodesOfType(k uint16, defaultDistance int) []*node {
	planets := m.planets[k]
	nodes := make([]*node, 0)
	for i := 0; i < len(planets); i++ {
		nodes = append(nodes, &node{
			position:        planets[i],
			bestDistanceYet: defaultDistance,
		})
	}
	return nodes
}

type position struct {
	X int
	Y int
}

//return the Manhattan's distance between two positions
func (p position) distanceTo(p0 position) int {
	abs := func(x int) int {
		if x > 0 {
			return x
		}
		return -x
	}
	return abs(p.X-p0.X) + abs(p.Y-p0.Y)
}

type node struct {
	position        position
	bestDistanceYet int
}
