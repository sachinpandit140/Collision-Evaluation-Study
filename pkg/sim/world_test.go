package sim

import (
	"testing"

	"github.com/sachinpandit/Collision-Evaluation-Study/pkg/geom"
)

func TestNewWorld_ZeroEntities(t *testing.T) {
	w := NewWorld(100, 100, 0, 42)
	if len(w.Entities) != 0 {
		t.Errorf("expected 0 entities, got %d", len(w.Entities))
	}
}

func TestNewWorld_SingleEntity_InBounds(t *testing.T) {
	w := NewWorld(100, 100, 1, 42)
	e := w.Entities[0]
	if e.Bounds.Min.X < 0 || e.Bounds.Max.X > 100 ||
		e.Bounds.Min.Y < 0 || e.Bounds.Max.Y > 100 {
		t.Errorf("entity out of world bounds: %+v", e.Bounds)
	}
}

func TestNewWorld_Determinism(t *testing.T) {
	w1 := NewWorld(100, 100, 50, 99)
	w2 := NewWorld(100, 100, 50, 99)
	for i := range w1.Entities {
		a, b := w1.Entities[i], w2.Entities[i]
		if a.Bounds != b.Bounds || a.Velocity != b.Velocity {
			t.Fatalf("entity %d differs between runs with same seed", i)
		}
	}
}

func TestTick_MovesEntities(t *testing.T) {
	w := NewWorld(1000, 1000, 1, 42)
	before := w.Entities[0].Bounds
	w.Tick(1.0)
	after := w.Entities[0].Bounds
	if before.Min.X == after.Min.X && before.Min.Y == after.Min.Y {
		t.Error("entity did not move after tick")
	}
}

func TestTick_ReflectsAtBoundary(t *testing.T) {
	w := &World{Width: 10, Height: 10, Entities: []Entity{
		{
			ID:       0,
			Bounds:   geom.AABB2{Min: geom.Vec2{X: 9, Y: 5}, Max: geom.Vec2{X: 11, Y: 7}},
			Velocity: geom.Vec2{X: 1, Y: 0},
		},
	}}
	w.Tick(1.0) // pushes Max.X to 12, exceeds Width=10
	if w.Entities[0].Velocity.X > 0 {
		t.Error("velocity should have reflected (negative) on X boundary")
	}
}
