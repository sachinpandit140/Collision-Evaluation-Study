package sim

import (
	"math/rand"

	"github.com/sachinpandit/Collision-Evaluation-Study/pkg/geom"
)

const entitySize = 2.0 // width and height of each entity AABB

// World holds the simulation state.
type World struct {
	Width, Height float64
	Entities      []Entity
}

// NewWorld creates a world and spawns n entities with deterministic random
// positions and velocities using the given seed.
func NewWorld(width, height float64, n int, seed int64) *World {
	rng := rand.New(rand.NewSource(seed))
	entities := make([]Entity, n)
	for i := range entities {
		x := rng.Float64() * (width - entitySize)
		y := rng.Float64() * (height - entitySize)
		vx := (rng.Float64() - 0.5) * 2 // velocity in [-1, 1]
		vy := (rng.Float64() - 0.5) * 2
		entities[i] = Entity{
			ID: uint32(i),
			Bounds: geom.AABB2{
				Min: geom.Vec2{X: x, Y: y},
				Max: geom.Vec2{X: x + entitySize, Y: y + entitySize},
			},
			Velocity: geom.Vec2{X: vx, Y: vy},
		}
	}
	return &World{Width: width, Height: height, Entities: entities}
}

// Tick advances every entity by its velocity and reflects off world bounds.
func (w *World) Tick(dt float64) {
	for i := range w.Entities {
		e := &w.Entities[i]
		e.Bounds.Min.X += e.Velocity.X * dt
		e.Bounds.Max.X += e.Velocity.X * dt
		e.Bounds.Min.Y += e.Velocity.Y * dt
		e.Bounds.Max.Y += e.Velocity.Y * dt

		if e.Bounds.Min.X < 0 || e.Bounds.Max.X > w.Width {
			e.Velocity.X = -e.Velocity.X
		}
		if e.Bounds.Min.Y < 0 || e.Bounds.Max.Y > w.Height {
			e.Velocity.Y = -e.Velocity.Y
		}
	}
}
