package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const intMax = int(^uint(0) >> 1)

func main() {
	var input Input
	input = getAndParseInput() // we read the input and store it in  a struct
	fmt.Println(solve(input))  //we compute and print the answer
}

type Input struct {
	oreCount      int
	mineralsCosts []uint8
	targetPrice   int
}

func getAndParseInput() (input Input) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 100), 4000000)

	var n int
	scanner.Scan()
	n, _ = strconv.Atoi(scanner.Text())
	priceList := make([]uint8, n)
	scanner.Scan()
	for i, iValue := range strings.SplitN(scanner.Text(), " ", n) {
		v, _ := strconv.Atoi(iValue)
		priceList[i] = uint8(v)
	}
	var b int
	scanner.Scan()
	b, _ = strconv.Atoi(scanner.Text())

	input.oreCount = n
	input.mineralsCosts = priceList
	input.targetPrice = b

	return input
}

//this algorithm starts from the price of the first element, adds a new one if this is cheaper than the wanted price
//and removes on if it is more expensive.
func solve(input Input) int {
	if input.targetPrice == 0 { //if we want to pay 0, we just buy nothing
		return 0
	}
	firstItemIndex := 0
	lastItemIndex := 0
	currentSum := int(input.mineralsCosts[0])

	smallestSoFar := intMax

	for true {
		if currentSum == input.targetPrice && lastItemIndex-firstItemIndex < smallestSoFar { //if the sum is the wanted
		// one
			smallestSoFar = lastItemIndex - firstItemIndex //we update the current smallest number of elements
		} else if currentSum < input.targetPrice { //if the sum is cheaper
			lastItemIndex++ //we add an item
			if lastItemIndex >= input.oreCount { //prevent out of bounds read
				break
			}
			currentSum += int(input.mineralsCosts[lastItemIndex]) //and update the sum
		} else { //if it is less
			currentSum -= int(input.mineralsCosts[firstItemIndex]) //we remove an element
			firstItemIndex++
			if firstItemIndex >= input.oreCount { //and prevent out of bounds read next time
				break
			}
		}
	}

	if smallestSoFar == intMax { //return -1 if nothing is foundo
		return -1
	}

	return smallestSoFar + 1
}
