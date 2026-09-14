package baseline

import (
	"testing"

	"github.com/sachinpandit140/Collision-Evaluation-Study/pkg/geom"
	"github.com/sachinpandit140/Collision-Evaluation-Study/pkg/sim"
)

func makeEntity(id uint32, x, y, size float64) sim.Entity {
	return sim.Entity{
		ID:     id,
		Bounds: geom.AABB2{Min: geom.Vec2{X: x, Y: y}, Max: geom.Vec2{X: x + size, Y: y + size}},
	}
}

func TestFindCollisions_Empty(t *testing.T) {
	pairs := FindCollisions(nil)
	if len(pairs) != 0 {
		t.Errorf("expected 0 pairs, got %d", len(pairs))
	}
}

func TestFindCollisions_SingleEntity(t *testing.T) {
	pairs := FindCollisions([]sim.Entity{makeEntity(0, 0, 0, 1)})
	if len(pairs) != 0 {
		t.Errorf("expected 0 pairs, got %d", len(pairs))
	}
}

func TestFindCollisions_TwoOverlapping(t *testing.T) {
	entities := []sim.Entity{
		makeEntity(0, 0, 0, 2),
		makeEntity(1, 1, 1, 2),
	}
	pairs := FindCollisions(entities)
	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(pairs))
	}
	if pairs[0] != [2]uint32{0, 1} {
		t.Errorf("unexpected pair: %v", pairs[0])
	}
}

func TestFindCollisions_TwoSeparated(t *testing.T) {
	entities := []sim.Entity{
		makeEntity(0, 0, 0, 1),
		makeEntity(1, 10, 10, 1),
	}
	pairs := FindCollisions(entities)
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
	pairs := FindCollisions(entities)
	// 3 entities all overlapping = 3 pairs: (0,1), (0,2), (1,2)
	if len(pairs) != 3 {
		t.Errorf("expected 3 pairs, got %d", len(pairs))
	}
}
