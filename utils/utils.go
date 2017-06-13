package utils

import (
	"fmt"
	"sort"
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

func PrintObjectsWithRating(ratingsMap map[string]int) {
	pl := make(PairList, 0, len(ratingsMap))
	for k, v := range ratingsMap {
		pl = append(pl, Pair{Key: k, Value: v})
	}
	sort.Sort(sort.Reverse(pl))

	fmt.Printf("=====\nОбъекты и их рейтинг\n")
	for i, p := range pl {
		fmt.Printf("%d: %s = %d\n", i+1, p.Key, p.Value)
	}
	fmt.Printf("=====\n")
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
