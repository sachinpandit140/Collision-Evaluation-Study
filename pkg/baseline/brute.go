package baseline

import "github.com/sachinpandit/Collision-Evaluation-Study/pkg/sim"

// FindCollisions tests all N*(N-1)/2 entity pairs for AABB overlap.
// Returns a slice of colliding entity ID pairs. This is the O(N²)
// correctness baseline.
func FindCollisions(entities []sim.Entity) [][2]uint32 {
	var pairs [][2]uint32
	for i := 0; i < len(entities); i++ {
		for j := i + 1; j < len(entities); j++ {
			if entities[i].Bounds.Overlaps(entities[j].Bounds) {
				pairs = append(pairs, [2]uint32{entities[i].ID, entities[j].ID})
			}
		}
	}
	return pairs
}
