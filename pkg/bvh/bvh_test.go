package bvh

import (
	"sort"
	"testing"

	"github.com/sachinpandit140/Collision-Evaluation-Study/pkg/geom"
	"github.com/sachinpandit140/Collision-Evaluation-Study/pkg/sim"
)

func makeEntity(id uint32, x, y, size float64) sim.Entity {
	half := size * 0.5
	return sim.Entity{
		ID: id,
		Bounds: geom.AABB2{
			Min: geom.Vec2{X: x - half, Y: y - half},
			Max: geom.Vec2{X: x + half, Y: y + half},
		},
	}
}

func sortPairs(pairs [][2]uint32) {
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i][0] == pairs[j][0] {
			return pairs[i][1] < pairs[j][1]
		}
		return pairs[i][0] < pairs[j][0]
	})
}

func TestBuild_Empty(t *testing.T) {
	tree := Build(nil)
	pairs := tree.FindCollisions(nil)
	if len(pairs) != 0 {
		t.Errorf("expected 0 pairs, got %d", len(pairs))
	}
}

func TestBuild_SingleEntity(t *testing.T) {
	entities := []sim.Entity{makeEntity(1, 0, 0, 1)}
	tree := Build(entities)
	pairs := tree.FindCollisions(entities)
	if len(pairs) != 0 {
		t.Errorf("expected 0 pairs, got %d", len(pairs))
	}
}

func TestFindCollisions_TwoOverlapping(t *testing.T) {
	entities := []sim.Entity{
		makeEntity(0, 0, 0, 2),
		makeEntity(1, 1, 1, 2),
	}
	tree := Build(entities)
	pairs := tree.FindCollisions(entities)
	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(pairs))
	}
	if pairs[0][0] != 0 || pairs[0][1] != 1 {
		t.Errorf("expected [0, 1], got %v", pairs[0])
	}
}

func TestFindCollisions_TwoSeparated(t *testing.T) {
	entities := []sim.Entity{
		makeEntity(0, 0, 0, 1),
		makeEntity(1, 10, 10, 1),
	}
	tree := Build(entities)
	pairs := tree.FindCollisions(entities)
	if len(pairs) != 0 {
		t.Errorf("expected 0 pairs, got %d", len(pairs))
	}
}

func TestFindCollisions_AllOverlapping(t *testing.T) {
	entities := []sim.Entity{
		makeEntity(0, 0, 0, 10),
		makeEntity(1, 1, 1, 10),
		makeEntity(2, 2, 2, 10),
	}
	tree := Build(entities)
	pairs := tree.FindCollisions(entities)
	if len(pairs) != 3 {
		t.Errorf("expected 3 pairs, got %d", len(pairs))
	}
}

func TestFindCollisions_MatchesBruteForce(t *testing.T) {
	entities := []sim.Entity{
		makeEntity(0, 0, 0, 3),
		makeEntity(1, 2, 2, 3),
		makeEntity(2, 4, 4, 3),
		makeEntity(3, 10, 10, 2),
		makeEntity(4, 11, 11, 2),
	}
	tree := Build(entities)
	pairs := tree.FindCollisions(entities)
	
	sortPairs(pairs)
	
	expected := [][2]uint32{{0, 1}, {1, 2}, {3, 4}}
	
	if len(pairs) != len(expected) {
		t.Fatalf("expected %d pairs, got %d", len(expected), len(pairs))
	}
	for i := range expected {
		if pairs[i] != expected[i] {
			t.Errorf("at index %d, expected %v, got %v", i, expected[i], pairs[i])
		}
	}
}

func TestRefit(t *testing.T) {
	entities := []sim.Entity{
		makeEntity(0, 0, 0, 1),
		makeEntity(1, 10, 10, 1),
	}
	tree := Build(entities)
	pairs := tree.FindCollisions(entities)
	if len(pairs) != 0 {
		t.Fatalf("expected 0 pairs initially, got %d", len(pairs))
	}
	
	// Move entity 1 to overlap entity 0
	entities[1] = makeEntity(1, 0.5, 0.5, 1)
	tree.Refit(entities)
	
	pairs = tree.FindCollisions(entities)
	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair after refit, got %d", len(pairs))
	}
	if pairs[0][0] != 0 || pairs[0][1] != 1 {
		t.Errorf("expected [0, 1], got %v", pairs[0])
	}
}
