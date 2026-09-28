package bvh

import (
	"math"
	"sort"

	"github.com/sachinpandit140/Collision-Evaluation-Study/pkg/geom"
	"github.com/sachinpandit140/Collision-Evaluation-Study/pkg/sim"
)

type node struct {
	box       geom.AABB2
	left      int32
	right     int32
	entityIdx int32
}

func (n *node) isLeaf() bool {
	return n.left == -1 && n.right == -1
}

type Tree struct {
	nodes []node
	root  int32
}

func Build(entities []sim.Entity) *Tree {
	if len(entities) == 0 {
		return &Tree{nodes: nil, root: -1}
	}

	indices := make([]int32, len(entities))
	for i := range entities {
		indices[i] = int32(i)
	}

	t := &Tree{
		nodes: make([]node, 0, len(entities)*2-1),
		root:  -1,
	}

	t.root = t.buildRecursive(entities, indices)
	return t
}

func center(box geom.AABB2, axis int) float64 {
	if axis == 0 {
		return box.Min.X + (box.Max.X-box.Min.X)*0.5
	}
	return box.Min.Y + (box.Max.Y-box.Min.Y)*0.5
}

func sortByAxis(entities []sim.Entity, indices []int32, axis int) {
	sort.Slice(indices, func(i, j int) bool {
		c1 := center(entities[indices[i]].Bounds, axis)
		c2 := center(entities[indices[j]].Bounds, axis)
		return c1 < c2
	})
}

func (t *Tree) buildRecursive(entities []sim.Entity, indices []int32) int32 {
	if len(indices) == 1 {
		idx := int32(len(t.nodes))
		t.nodes = append(t.nodes, node{
			box:       entities[indices[0]].Bounds,
			left:      -1,
			right:     -1,
			entityIdx: indices[0],
		})
		return idx
	}

	// Compute parent bounds
	var parentBox geom.AABB2
	for i, idx := range indices {
		if i == 0 {
			parentBox = entities[idx].Bounds
		} else {
			parentBox = parentBox.Merge(entities[idx].Bounds)
		}
	}
	parentArea := parentBox.Area()
	if parentArea < 1e-10 {
		parentArea = 1e-10
	}

	bestCost := math.MaxFloat64
	bestAxis := -1
	bestSplit := -1

	for axis := 0; axis < 2; axis++ {
		sortByAxis(entities, indices, axis)

		leftBoxes := make([]geom.AABB2, len(indices))
		var curLeft geom.AABB2
		for i := 0; i < len(indices); i++ {
			if i == 0 {
				curLeft = entities[indices[i]].Bounds
			} else {
				curLeft = curLeft.Merge(entities[indices[i]].Bounds)
			}
			leftBoxes[i] = curLeft
		}

		rightBoxes := make([]geom.AABB2, len(indices))
		var curRight geom.AABB2
		for i := len(indices) - 1; i >= 0; i-- {
			if i == len(indices)-1 {
				curRight = entities[indices[i]].Bounds
			} else {
				curRight = curRight.Merge(entities[indices[i]].Bounds)
			}
			rightBoxes[i] = curRight
		}

		for i := 0; i < len(indices)-1; i++ {
			nLeft := float64(i + 1)
			nRight := float64(len(indices) - i - 1)
			
			saLeft := leftBoxes[i].Area()
			saRight := rightBoxes[i+1].Area()

			cost := (saLeft/parentArea)*nLeft + (saRight/parentArea)*nRight
			if cost < bestCost {
				bestCost = cost
				bestAxis = axis
				bestSplit = i + 1
			}
		}
	}

	sortByAxis(entities, indices, bestAxis)

	nodeIdx := int32(len(t.nodes))
	t.nodes = append(t.nodes, node{
		box:       parentBox,
		entityIdx: -1,
	})

	leftIndices := make([]int32, bestSplit)
	copy(leftIndices, indices[:bestSplit])
	rightIndices := make([]int32, len(indices)-bestSplit)
	copy(rightIndices, indices[bestSplit:])

	left := t.buildRecursive(entities, leftIndices)
	right := t.buildRecursive(entities, rightIndices)

	t.nodes[nodeIdx].left = left
	t.nodes[nodeIdx].right = right

	return nodeIdx
}

func (t *Tree) Refit(entities []sim.Entity) {
	if t.root == -1 {
		return
	}
	t.refitRecursive(t.root, entities)
}

func (t *Tree) refitRecursive(idx int32, entities []sim.Entity) {
	n := &t.nodes[idx]
	if n.isLeaf() {
		n.box = entities[n.entityIdx].Bounds
		return
	}

	t.refitRecursive(n.left, entities)
	t.refitRecursive(n.right, entities)

	n.box = t.nodes[n.left].box.Merge(t.nodes[n.right].box)
}

func (t *Tree) FindCollisions(entities []sim.Entity) [][2]uint32 {
	if t.root == -1 {
		return nil
	}
	var pairs [][2]uint32
	t.selfCollide(t.root, entities, &pairs)
	return pairs
}

func (t *Tree) selfCollide(idx int32, entities []sim.Entity, pairs *[][2]uint32) {
	n := &t.nodes[idx]
	if n.isLeaf() {
		return
	}
	
	t.selfCollide(n.left, entities, pairs)
	t.selfCollide(n.right, entities, pairs)
	t.crossCollide(n.left, n.right, entities, pairs)
}

func (t *Tree) crossCollide(idxA, idxB int32, entities []sim.Entity, pairs *[][2]uint32) {
	a := &t.nodes[idxA]
	b := &t.nodes[idxB]

	if !a.box.Overlaps(b.box) {
		return
	}

	if a.isLeaf() && b.isLeaf() {
		e1 := entities[a.entityIdx]
		e2 := entities[b.entityIdx]
		if e1.Bounds.Overlaps(e2.Bounds) {
			if e1.ID < e2.ID {
				*pairs = append(*pairs, [2]uint32{e1.ID, e2.ID})
			} else {
				*pairs = append(*pairs, [2]uint32{e2.ID, e1.ID})
			}
		}
		return
	}

	if a.isLeaf() {
		t.crossCollide(idxA, b.left, entities, pairs)
		t.crossCollide(idxA, b.right, entities, pairs)
	} else if b.isLeaf() {
		t.crossCollide(a.left, idxB, entities, pairs)
		t.crossCollide(a.right, idxB, entities, pairs)
	} else {
		// Both internal, descend into the larger one
		if a.box.Area() > b.box.Area() {
			t.crossCollide(a.left, idxB, entities, pairs)
			t.crossCollide(a.right, idxB, entities, pairs)
		} else {
			t.crossCollide(idxA, b.left, entities, pairs)
			t.crossCollide(idxA, b.right, entities, pairs)
		}
	}
}
