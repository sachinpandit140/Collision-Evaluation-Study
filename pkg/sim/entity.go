package sim

import "github.com/sachinpandit140/Collision-Evaluation-Study/pkg/geom"

// Entity represents a single moving object in the simulation.
type Entity struct {
	ID       uint32
	Bounds   geom.AABB2
	Velocity geom.Vec2
}
