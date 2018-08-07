package main

import (
	"log"
)

var dict = map[string][]string{
	"A": nil,
	"B": []string{"A"},
	"C": []string{"B"},
	"D": []string{"B", "C"},
	"E": []string{"D"},
	"F": []string{"E"},
}

func dig_deep(child, seekFor string) (found bool) {
	parents := dict[child]
	if parents == nil {
		return false
	} else {
		for _, p := range parents {
			if p == seekFor {
				found = true
				break
			}
		}
		if found {
			return
		}
		for _, p := range parents {
			if found = dig_deep(p, seekFor); found {
				return
			}
		}
	}
	return
}

func task() {
	var test = [][]string{
		[]string{"B", "A"},
		[]string{"C", "B"},
		[]string{"D", "A"},
		[]string{"A", "D"},
		[]string{"F", "A"},
		[]string{"F", "C"},
		[]string{"C", "F"},
	}
	for _, scenario := range test {
		log.Printf(" %v ", dig_deep(scenario[0], scenario[1]))

	}
}
