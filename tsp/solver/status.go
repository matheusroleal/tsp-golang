package solver

import (
	"fmt"
	"sync"
)

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
