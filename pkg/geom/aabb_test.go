package geom

import "testing"

func TestOverlaps_Intersecting(t *testing.T) {
	a := AABB2{Min: Vec2{0, 0}, Max: Vec2{2, 2}}
	b := AABB2{Min: Vec2{1, 1}, Max: Vec2{3, 3}}
	if !a.Overlaps(b) {
		t.Error("expected overlap")
	}
}

func TestOverlaps_Separated(t *testing.T) {
	a := AABB2{Min: Vec2{0, 0}, Max: Vec2{1, 1}}
	b := AABB2{Min: Vec2{5, 5}, Max: Vec2{6, 6}}
	if a.Overlaps(b) {
		t.Error("expected no overlap")
	}
}

func TestOverlaps_EdgeTouching(t *testing.T) {
	a := AABB2{Min: Vec2{0, 0}, Max: Vec2{1, 1}}
	b := AABB2{Min: Vec2{1, 0}, Max: Vec2{2, 1}}
	if !a.Overlaps(b) {
		t.Error("edge-touching AABBs should overlap (inclusive bounds)")
	}
}

func TestOverlaps_ZeroSizePoint(t *testing.T) {
	a := AABB2{Min: Vec2{1, 1}, Max: Vec2{1, 1}}
	b := AABB2{Min: Vec2{0, 0}, Max: Vec2{2, 2}}
	if !a.Overlaps(b) {
		t.Error("zero-size AABB inside another should overlap")
	}
}

func TestMerge(t *testing.T) {
	a := AABB2{Min: Vec2{0, 0}, Max: Vec2{1, 1}}
	b := AABB2{Min: Vec2{3, 3}, Max: Vec2{5, 5}}
	m := a.Merge(b)
	if m.Min.X != 0 || m.Min.Y != 0 || m.Max.X != 5 || m.Max.Y != 5 {
		t.Errorf("unexpected merge result: %+v", m)
	}
}

func TestArea(t *testing.T) {
	a := AABB2{Min: Vec2{0, 0}, Max: Vec2{3, 4}}
	if a.Area() != 12 {
		t.Errorf("expected area 12, got %f", a.Area())
	}
}

func TestArea_ZeroSize(t *testing.T) {
	a := AABB2{Min: Vec2{1, 1}, Max: Vec2{1, 1}}
	if a.Area() != 0 {
		t.Errorf("expected area 0, got %f", a.Area())
	}
}

func TestContains(t *testing.T) {
	a := AABB2{Min: Vec2{0, 0}, Max: Vec2{10, 10}}
	if !a.Contains(Vec2{5, 5}) {
		t.Error("point inside should be contained")
	}
	if a.Contains(Vec2{-1, 5}) {
		t.Error("point outside should not be contained")
	}
	if !a.Contains(Vec2{0, 0}) {
		t.Error("point on edge should be contained (inclusive)")
	}
}
