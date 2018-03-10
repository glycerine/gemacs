package tulib

import (
	"testing"
)

func TestIntersection(t *testing.T) {
	data := []struct {
		a        Rect
		b        Rect
		expected Rect
	}{
		{
			Rect{0, 0, 0, 0},
			Rect{0, 0, 0, 0},
			Rect{0, 0, 0, 0},
		},
		{
			Rect{1, 1, 1, 1},
			Rect{3, 3, 1, 1},
			Rect{0, 0, 0, 0},
		},
		{
			Rect{1, 1, 3, 3},
			Rect{3, 3, 1, 1},
			Rect{3, 3, 1, 1},
		},
		{
			Rect{1, 1, 5, 5},
			Rect{3, 3, 1, 1},
			Rect{3, 3, 1, 1},
		},
		{
			Rect{3, 3, 1, 1},
			Rect{3, 3, 3, 3},
			Rect{3, 3, 1, 1},
		},
		{
			Rect{5, 5, 1, 1},
			Rect{3, 3, 1, 1},
			Rect{0, 0, 0, 0},
		},
	}
	for _, item := range data {
		actual := item.a.Intersection(item.b)
		if actual != item.expected {
			t.Fatalf("Failed: %v != %v", item, actual)
		}
	}
}
