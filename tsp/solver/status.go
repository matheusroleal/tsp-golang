package solver

import (
	"log"
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
	sync.Mutex
	wg sync.WaitGroup
}

// NewStatus returns a new Status heap of segment size N
func NewStatus(adjMtrx [][]uint, segSize int) *Status {
	if segSize > 1<<16 {
		segSize = 1 << 16
	}

	status := Status{
		// TSP Stuff
		adjMtrx:  adjMtrx,
		vtxCount: int8(len(adjMtrx)),
		// Heap Stuff
		arr:     make([]*Element, segSize),
		segSize: segSize,
	}

	return &status
}

// Put inserts an element into the priority queue
func (stat *Status) Put(e *Element) {
	stat.Lock()
	defer stat.Unlock()

	stat.curSize++
	if stat.curSize%stat.segSize == 0 {
		log.Println("[INF] Resizing heap to", stat.curSize+stat.segSize, "elements")
		stat.arr = append(stat.arr, make([]*Element, stat.segSize)...)
	}

	stat.arr[stat.curSize-1] = e
	stat.up(stat.curSize - 1)
}

// Get returns the first element of the priority queue
func (stat *Status) Get() *Element {
	stat.Lock()
	defer stat.Unlock()

	if stat.curSize == 0 {
		return nil
	}

	v := stat.arr[0]
	stat.curSize--
	stat.arr[0] = stat.arr[stat.curSize]

	stat.down(0)

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
		if v.less(cv) {
			stat.arr[i] = cv
			i = child
			continue
		}
		break
	}

	if child < stat.curSize {
		cv := stat.arr[child]
		if v.less(cv) {
			stat.arr[i] = cv
			i = child
		}
	}

	stat.arr[i] = v
}

func (stat *Status) up(i int) {
	v := stat.arr[i]
	for parent := (i - 1) >> 1; i > 0; parent = (i - 1) >> 1 {
		pv := stat.arr[parent]
		if v.greater(pv) {
			stat.arr[i] = pv
			i = parent
			continue
		}
		break
	}

	stat.arr[i] = v
}
