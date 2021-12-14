package utils

import (
	"fmt"
)

func ReadNCities() int {
	var n_cidades int
	fmt.Printf("Numero de cidades:\n")
	fmt.Scanf("%d", &n_cidades)
	return n_cidades
}

func ReadNThreads() int {
	var n_threads int
	fmt.Printf("Numero de goroutines:\n")
	fmt.Scanf("%d", &n_threads)
	return n_threads
}

func ReadMatrixPath() string {
	var word string
	fmt.Printf("Path para arquivo com a matriz:\n")
	fmt.Scanf("%s", &word)
	return word
}
