package geom

import "math"

// Vec2 is a 2D vector.
type Vec2 struct{ X, Y float64 }

// AABB2 is an axis-aligned bounding box in 2D.
type AABB2 struct {
	Min, Max Vec2
}

// Overlaps returns true if a and b intersect (inclusive bounds).
func (a AABB2) Overlaps(b AABB2) bool {
	return a.Min.X <= b.Max.X && a.Max.X >= b.Min.X &&
		a.Min.Y <= b.Max.Y && a.Max.Y >= b.Min.Y
}

// Merge returns the smallest AABB containing both a and b.
func (a AABB2) Merge(b AABB2) AABB2 {
	return AABB2{
		Min: Vec2{math.Min(a.Min.X, b.Min.X), math.Min(a.Min.Y, b.Min.Y)},
		Max: Vec2{math.Max(a.Max.X, b.Max.X), math.Max(a.Max.Y, b.Max.Y)},
	}
}

// Area returns the area of the bounding box.
func (a AABB2) Area() float64 {
	return (a.Max.X - a.Min.X) * (a.Max.Y - a.Min.Y)
}

// Contains returns true if the point p is inside the AABB (inclusive).
func (a AABB2) Contains(p Vec2) bool {
	return p.X >= a.Min.X && p.X <= a.Max.X &&
		p.Y >= a.Min.Y && p.Y <= a.Max.Y
}
