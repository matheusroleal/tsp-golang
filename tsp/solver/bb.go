package solver

import (
	"fmt"
	"runtime"
	"sync"
)

// Element implements a branch in the TSP-Tree
type Element struct {
	FBPath   []int8 // fwd + bck
	LstVtx   int8
	Count    int8
	Boundary uint
}

// NewElement returns a new Element for the Status heap
func NewElement(fbpath []int8, lv, cnt int8) *Element {
	return &Element{fbpath, lv, cnt, 0}
}

// Status a queue of status elements
type Status struct {
	// TSP Stuff
	adjMtrx  [][]uint
	solution *Element
	solved   bool
	vtxCount int8
	// Heap Stuff
	arr     []*Element
	segSize int
	curSize int
	// Sync Stuff
	lckr sync.Mutex
	wg   sync.WaitGroup
}

// NewStatus returns a new Status heap of segment size N
func NewStatus(adjMtrx [][]uint, segSize int) *Status {
	return &Status{
		// TSP Stuff
		adjMtrx, nil, false, int8(len(adjMtrx)),
		// Heap Stuff
		make([]*Element, segSize), segSize, 0,
		// Sync Stuff
		sync.Mutex{}, sync.WaitGroup{},
	}
}

// Put inserts an element into the priority queue
func (stat *Status) Put(e *Element) {
	stat.lckr.Lock()
	stat.curSize++
	if stat.curSize%stat.segSize == 0 {
		fmt.Println("[INF] Resizing heap to", stat.curSize+stat.segSize, "elements")
		arr := make([]*Element, stat.curSize+stat.segSize)
		copy(arr, stat.arr)
		stat.arr = arr
	}
	stat.arr[stat.curSize-1] = e
	stat.up(stat.curSize - 1)
	stat.lckr.Unlock()
}

// Get returns the first element of the priority queue
func (stat *Status) Get() *Element {
	stat.lckr.Lock()
	if stat.curSize == 0 {
		return nil
	}
	v := stat.arr[0]
	stat.curSize--
	stat.arr[0] = stat.arr[stat.curSize]
	go stat.down(0)
	return v
}

// less compares an Element to another Element
func (e *Element) less(other *Element) bool {
	if e.Boundary == other.Boundary {
		return e.Count < other.Count
	}
	return e.Boundary > other.Boundary
}

// greater compares an Element to another Element
func (e *Element) greater(other *Element) bool {
	if e.Boundary == other.Boundary {
		return e.Count > other.Count
	}
	return e.Boundary < other.Boundary
}

func (stat *Status) down(i int) {
	v := stat.arr[i]
	child := (i << 1) + 1
	for l := stat.curSize; child+1 < l; child = (i << 1) + 1 {
		if stat.arr[child].less(stat.arr[child+1]) {
			child++
		}
		cv := stat.arr[child]
		b := v.less(cv)
		if b {
			stat.arr[i] = cv
			i = child
		} else {
			break
		}
	}
	if child < stat.curSize {
		cv := stat.arr[child]
		if v.less(cv) {
			stat.arr[i] = cv
			i = child
		}
	}
	stat.arr[i] = v
	stat.lckr.Unlock()
}

func (stat *Status) up(i int) {
	v := stat.arr[i]
	parent := (i - 1) >> 1
	for ; i > 0; parent = (i - 1) >> 1 {
		pv := stat.arr[parent]
		if v.greater(pv) {
			stat.arr[i] = pv
			i = parent
		} else {
			break
		}
	}
	stat.arr[i] = v
}

// TSPBB calculates the Traveling Salesman Problem on a given
// edge matrix and returns the best value and the best path while
// utilizing goroutines
func TSPBB(mtrx [][]uint, maxProcs, segSize int, grCnt int8) (uint, []int8) {
	runtime.GOMAXPROCS(maxProcs)
	status := NewStatus(mtrx, segSize)
	var i int8
	rootFBPath := make([]int8, (status.vtxCount<<1)+2)
	for i = 0; i < status.vtxCount; i++ {
		rootFBPath[i] = -1
		rootFBPath[status.vtxCount+i] = -1
	}
	status.Put(NewElement(rootFBPath, 0, 1))
	for i = 0; i < grCnt; i++ {
		status.wg.Add(1)
		go extend(status)
	}
	status.wg.Wait()
	if status.solved {
		return status.solution.Boundary, elemToPath(status)
	}
	return 2147483647, make([]int8, 0)
}

func extend(status *Status) {
	var candidate *Element
	var i int8
	for status.curSize > 0 {
		if status.solved {
			break
		}
		candidate = status.Get()
		if candidate.Count == status.vtxCount+1 {
			status.solution = candidate
			status.solved = true
		} else {
			if candidate.Count == status.vtxCount {
				i = 0
			} else {
				i = 1
			}
			for ; i < status.vtxCount; i++ {
				if candidate.FBPath[status.vtxCount+i] == -1 &&
					status.adjMtrx[candidate.LstVtx][i] != 0 {
					status.Put(getNewElement(status, candidate, i))
				}
			}
		}
	}
	status.wg.Done()
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
		} else {
			// Else we have to cycle through the matrix to find the lowest value
			min = ^uint(0)
			for j = 0; j < status.vtxCount; j++ {
				if v = status.adjMtrx[i][j]; v != 0 && v < min {
					min = v
				}
			}
			out += min
		}
	}
	// Incoming edges
	var in uint
	for i = 0; i < status.vtxCount; i++ {
		if e.FBPath[status.vtxCount+i] != -1 {
			// If there is a path we can add it's value immediately
			in += status.adjMtrx[e.FBPath[status.vtxCount+i]][i]
		} else {
			// Else we have to cycle through the matrix to find the lowest value
			min = ^uint(0)
			for j = 0; j < status.vtxCount; j++ {
				if v = status.adjMtrx[j][i]; v != 0 && v < min {
					min = v
				}
			}
			in += min
		}
	}
	//if (in == 38408 || out == 38408) && e.Count == 20 {
	//	fmt.Println("38408")
	//}
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
	path := make([]int8, status.vtxCount)
	var i, next int8 // starts with 0 anyways
	for i = 0; i < status.vtxCount; i++ {
		path[i] = next
		next = status.solution.FBPath[next]
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
