package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	input := getInput() //we read from stdin
	fmt.Println(removeInterferences( //and remove interferences
		input,
		func(message string) (cleaned string) { //this filter removes the dots
			cleaned = strings.Replace(message, ".", "", -1)
			return
		},
		func(message string) (cleaned string) { //this one removes characters between two '*'
			toKeep := make([]string, 0)

			var strSplit []string
			strSplit = strings.Split(message, "*")

			for i, s := range strSplit { //we keep only the even indexed strings
				if i%2 == 0 {
					toKeep = append(toKeep, s)
				}
			}
			return strings.Join(toKeep, "") //and join them
		}))
}

func getInput() (inputStr string) {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan() //we scan twice to discard the first line, as the range operator will determine
	scanner.Scan() 	//the length automatically

	inputStr = scanner.Text()
	inputStr = strings.Replace(inputStr, "\n", "", -1) //we replace the \n at the end of the input

	return inputStr
}

//This functions takes as input a string and any number of functions
//the string will be passed through each function which should remove all the different kind of interferences
func removeInterferences(message string, filters ...func(message string) (cleaned string)) string {
	for _, f := range filters { //for each filter
		message = f(message) //we pass the message though it
	}
	return message
}
