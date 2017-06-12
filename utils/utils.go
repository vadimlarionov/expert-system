package utils

import (
	"fmt"
	"strconv"
)

func CompareStrings(first, operation, second string) bool {
	switch operation {
	case "=":
		return first == second
	case "!=":
		return first != second
	default:
		fmt.Printf("Unsupported operation \"%s\" for strings \"%s\" and \"%s\"\n",
			operation, first, second)
		return false
	}
}

func CompareInts(first, operation, second string) bool {
	firstInt, err := strconv.Atoi(first)
	if err != nil {
		fmt.Printf("Can't cast \"%s\" to int: %s\n", first, err)
		return false
	}

	secondInt, err := strconv.Atoi(second)
	if err != nil {
		fmt.Printf("Can't cast \"%s\" to int: %s\n", second, err)
		return false
	}

	switch operation {
	case "=":
		return firstInt == secondInt
	case "!=":
		return firstInt != secondInt
	case "<":
		return firstInt < secondInt
	case "<=":
		return firstInt <= secondInt
	case ">":
		return firstInt > secondInt
	case ">=":
		return firstInt >= secondInt
	default:
		fmt.Printf("Unexpected operation \"%s\"", operation)
		return false
	}
}
