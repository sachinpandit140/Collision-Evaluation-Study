package kdtree

import (
	"sort"

	"github.com/sachinpandit140/Collision-Evaluation-Study/pkg/geom"
	"github.com/sachinpandit140/Collision-Evaluation-Study/pkg/sim"
)

type Tree struct {
	nodes []node
	root  int32
}

type node struct {
	box       geom.AABB2
	left      int32
	right     int32
	entityIdx int32
}

func (n *node) isLeaf() bool {
	return n.left == -1 && n.right == -1
}

func center(box geom.AABB2, axis int) float64 {
	if axis == 0 {
		return box.Min.X + (box.Max.X-box.Min.X)/2
	}
	return box.Min.Y + (box.Max.Y-box.Min.Y)/2
}

func Build(entities []sim.Entity) *Tree {
	if len(entities) == 0 {
		return &Tree{
			nodes: nil,
			root:  -1,
		}
	}

	indices := make([]int32, len(entities))
	for i := range entities {
		indices[i] = int32(i)
	}

	t := &Tree{
		nodes: make([]node, 0, len(entities)*2),
	}

	t.root = t.buildRecursive(entities, indices, 0)
	return t
}

func (t *Tree) buildRecursive(entities []sim.Entity, indices []int32, depth int) int32 {
	if len(indices) == 1 {
		idx := indices[0]
		e := entities[idx]
		t.nodes = append(t.nodes, node{
			box:       e.Bounds,
			left:      -1,
			right:     -1,
			entityIdx: idx,
		})
		return int32(len(t.nodes) - 1)
	}

	axis := depth % 2

	sort.Slice(indices, func(i, j int) bool {
		ei := entities[indices[i]]
		ej := entities[indices[j]]
		ci := center(ei.Bounds, axis)
		cj := center(ej.Bounds, axis)
		return ci < cj
	})

	mid := len(indices) / 2
	leftIndices := indices[:mid]
	rightIndices := indices[mid:]

	nodeIdx := int32(len(t.nodes))
	t.nodes = append(t.nodes, node{
		left:      -1,
		right:     -1,
		entityIdx: -1,
	})

	leftNodeIdx := t.buildRecursive(entities, leftIndices, depth+1)
	rightNodeIdx := t.buildRecursive(entities, rightIndices, depth+1)

	box := t.nodes[leftNodeIdx].box.Merge(t.nodes[rightNodeIdx].box)

	t.nodes[nodeIdx].box = box
	t.nodes[nodeIdx].left = leftNodeIdx
	t.nodes[nodeIdx].right = rightNodeIdx

	return nodeIdx
}

func (t *Tree) FindCollisions(entities []sim.Entity) [][2]uint32 {
	var pairs [][2]uint32
	if t.root == -1 {
		return pairs
	}
	
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
	nA := &t.nodes[idxA]
	nB := &t.nodes[idxB]

	if !nA.box.Overlaps(nB.box) {
		return
	}

	if nA.isLeaf() && nB.isLeaf() {
		eA := entities[nA.entityIdx]
		eB := entities[nB.entityIdx]
		if eA.Bounds.Overlaps(eB.Bounds) {
			idA, idB := eA.ID, eB.ID
			if idA > idB {
				idA, idB = idB, idA
			}
			*pairs = append(*pairs, [2]uint32{idA, idB})
		}
		return
	}

	if nA.isLeaf() {
		t.crossCollide(idxA, nB.left, entities, pairs)
		t.crossCollide(idxA, nB.right, entities, pairs)
	} else if nB.isLeaf() {
		t.crossCollide(nA.left, idxB, entities, pairs)
		t.crossCollide(nA.right, idxB, entities, pairs)
	} else {
		// Both internal. Descend into the larger node.
		if nA.box.Area() > nB.box.Area() {
			t.crossCollide(nA.left, idxB, entities, pairs)
			t.crossCollide(nA.right, idxB, entities, pairs)
		} else {
			t.crossCollide(idxA, nB.left, entities, pairs)
			t.crossCollide(idxA, nB.right, entities, pairs)
		}
	}
}
