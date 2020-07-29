package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	weights := getAndParseInput() //we read from stdin and convert it to a slice of floats each containing
	// a weight
	fmt.Println(totalFuelRequired(weights)) //then we compute the mass of fuel which is required and print the result
}

func getAndParseInput() (inputFloat []float64) {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan() //we scan twice to discard the first line, as the range operator will determine
	scanner.Scan() //the length automatically
	inputStr := scanner.Text()

	inputStr = strings.Replace(inputStr, "\n", "", -1) //we remove the '\n' at the end of the input

	inputSplit := strings.Split(inputStr, " ") //split each value into substrings

	inputFloat = []float64{}
	for _, str := range inputSplit { //and convert it all to float64
		f, err := strconv.ParseFloat(str, 64)
		if err != nil { //this should only happen with a badly formatted input
			panic(err)
		}
		inputFloat = append(inputFloat, f)
	}

	return inputFloat
}

func totalFuelRequired(weights []float64) (totalFuel float64) {
	for _, w := range weights {
		totalFuel += fuelRequired(w)
	}
	return totalFuel
}

func fuelRequired(weight float64) float64 {
	if weight > 90 {
		return 80
	}
	return 60
}

