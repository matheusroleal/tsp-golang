package solver

import (
	"runtime"
)

// TSPBB calculates the Traveling Salesman Problem on a given
// edge matrix and returns the best value and the best path while
// utilizing goroutines
func TSPBB(mtrx [][]uint, maxProcs, segSize int, grCnt int8) (uint, []int8) {
	runtime.GOMAXPROCS(maxProcs)
	status := NewStatus(mtrx, segSize)

	rootFBPath := make([]int8, (status.vtxCount<<1)+2)
	for i := int8(0); i < status.vtxCount; i++ {
		rootFBPath[i] = -1
		rootFBPath[status.vtxCount+i] = -1
	}

	status.Put(NewElement(rootFBPath, 0, 1))
	status.wg.Add(int(grCnt))

	for i := int8(0); i < grCnt; i++ {
		go func(st *Status) {
			defer st.wg.Done()

			var c *Element
			var i int8
			for st.curSize > 0 {
				if st.solved {
					break
				}
				c = st.Get()
				if c.Count == st.vtxCount+1 {
					st.solution = c
					st.solved = true
					continue
				}
				i = 0
				if c.Count != st.vtxCount {
					i = 1
				}
				for ; i < st.vtxCount; i++ {
					if c.FBPath[st.vtxCount+i] == -1 &&
						st.adjMtrx[c.LstVtx][i] != 0 {
						st.Put(getNewElement(st, c, i))
					}
				}
			}
		}(status)
	}

	status.wg.Wait()

	if status.solved {
		return status.solution.Boundary, elemToPath(status)
	}

	return 2147483647, make([]int8, 0)
}

// UpdateBoundary updates the boundary of the Status Element
// TODO: Use more PQs to manage the edges to update more quickly
func UpdateBoundary(status *Status, e *Element) {
	// Declaring variables so we don't need to allocate space multiple times
	var min, v uint
	var j, i int8
	// Outgoing edges
	var out uint
	for i = 0; i < status.vtxCount; i++ {
		if e.FBPath[i] != -1 {
			// If there is a path we can add it's value immediately
			out += status.adjMtrx[i][e.FBPath[i]]

			continue
		}
		// Else we have to cycle through the matrix to find the lowest value
		min = ^uint(0)
		for j = 0; j < status.vtxCount; j++ {
			if v = status.adjMtrx[i][j]; v != 0 && v < min {
				min = v
			}
		}
		out += min
	}
	// Incoming edges
	var in uint
	for i = 0; i < status.vtxCount; i++ {
		if e.FBPath[status.vtxCount+i] != -1 {
			// If there is a path we can add it's value immediately
			in += status.adjMtrx[e.FBPath[status.vtxCount+i]][i]
			continue
		}
		// Else we have to cycle through the matrix to find the lowest value
		min = ^uint(0)
		for j = 0; j < status.vtxCount; j++ {
			if v = status.adjMtrx[j][i]; v != 0 && v < min {
				min = v
			}
		}
		in += min
	}
	if in > out {
		e.Boundary = in
	} else {
		e.Boundary = out
	}
}

// Adds a vertex to the paths of a candidate
func getNewElement(status *Status, candidate *Element, i int8) *Element {
	fbPath := make([]int8, (status.vtxCount<<1)+2)
	copy(fbPath, candidate.FBPath)
	fbPath[candidate.LstVtx] = i
	fbPath[status.vtxCount+i] = candidate.LstVtx
	e := &Element{fbPath, i, candidate.Count + 1, 0}
	UpdateBoundary(status, e)
	return e
}

func elemToPath(status *Status) []int8 {
	path, next := make([]int8, status.vtxCount), int8(0)
	for i := int8(0); i < status.vtxCount; i++ {
		path[i], next = next, status.solution.FBPath[next]
	}

	return path
}

func ActualCost(path []int8, adjMatrix [][]uint) uint {
	j := path[len(path)-1]
	var sum uint
	for i := 0; i < len(path); i++ {
		sum += adjMatrix[j][path[i]]
		j = path[i]
	}
	return sum
}
