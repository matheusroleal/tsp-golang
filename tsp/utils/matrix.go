package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadMatrix(rows int, cols int, filename string) ([][]uint, error) {
	var a = make([][]uint, rows)
	for x := range a {
		a[x] = make([]uint, cols)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		in := scanner.Text()
		parts := strings.Split(in, "  ")
		for j := 0; j < len(parts); j++ {
			result64, err := strconv.ParseFloat(parts[j], 64)
			if err != nil {
				fmt.Printf("Type: %T \n", result64)
				fmt.Println(err)
			}
			a[i][j] = uint(result64)
		}
		i++
	}
	return a, scanner.Err()
}

func ShowMatrix(n_cities int, adj_m [][]uint) {
	fmt.Printf("\nMatriz definida\n")
	for i := 0; i < n_cities; i++ {
		for j := 0; j < n_cities; j++ {
			fmt.Printf("%-4d ", adj_m[i][j])
		}
		fmt.Println("")
	}
	fmt.Printf("\n")
}
