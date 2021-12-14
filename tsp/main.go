package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/matheusroleal/tsp-golang/tsp/solver"
	utils "github.com/matheusroleal/tsp-golang/tsp/utils"
)

func main() {
	tspCLI()
}

func tspCLI() {
	var filename string
	var size int
	var threads int

	if len(os.Args) > 1 && len(os.Args) < 5 {
		filename = os.Args[1]
		size, _ = strconv.Atoi(os.Args[2])
		threads, _ = strconv.Atoi(os.Args[3])
	} else {
		fmt.Println("[ERROR]: Invalid quantity of atributes")
		os.Exit(1)
	}

	g, err := utils.ReadMatrix(size, size, filename)
	if err != nil {
		fmt.Println("[ERROR] :", err)
	}
	// utils.ShowMatrix(size, g)
	start := time.Now()
	cost, _ := solver.TSPBB(g, 8, 1<<28, int8(threads))
	fmt.Println(time.Since(start))
	// fmt.Println("[INF] TSP-Time:", time.Since(start))
	// fmt.Println("[OUT] Path:", path)
	fmt.Println("[OUT] Predicted Cost:", cost)
	// fmt.Println("[INF] Acutal Cost:", solver.ActualCost(path, g))
}
