package utils

import (
	"fmt"
)

func ReadNCities() int {
	var cities int
	fmt.Printf("Number of cities:\n")
	fmt.Scanf("%d", &cities)
	return cities
}

func ReadNThreads() int {
	var n_threads int
	fmt.Printf("Number of goroutines:\n")
	fmt.Scanf("%d", &n_threads)
	return n_threads
}

func ReadMatrixPath() string {
	var word string
	fmt.Printf("Path to the matrix file:\n")
	fmt.Scanf("%s", &word)
	return word
}
