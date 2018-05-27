package main

import (
	"fmt"
	"sort"

	"vendored"
)

func main() {
	people := []vendored.Person{
		{"Bob", 31},
		{"John", 42},
		{"Michael", 17},
		{"Jenny", 26},
	}
	sort.Sort(vendored.ByAge(people))
	fmt.Println(people)
}
