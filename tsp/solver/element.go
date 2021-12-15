package solver

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
