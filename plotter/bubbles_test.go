package plotter

import (
	"code.google.com/p/plotinum/vg"
	"testing"
)

func TestBubblesRadius(t *testing.T) {
	b := &Bubbles{
		MinRadius: vg.Length(0),
		MaxRadius: vg.Length(1),
	}

	tests := []struct {
		minz, maxz, z float64
		r             vg.Length
	}{
		{0, 0, 0, vg.Length(0.5)},
		{1, 1, 1, vg.Length(0.5)},
		{0, 1, 0, vg.Length(0)},
		{0, 1, 1, vg.Length(1)},
		{0, 1, 0.5, vg.Length(0.5)},
		{0, 2, 1, vg.Length(0.5)},
		{0, 4, 0, vg.Length(0)},
		{0, 4, 1, vg.Length(0.25)},
		{0, 4, 2, vg.Length(0.5)},
		{0, 4, 3, vg.Length(0.75)},
		{0, 4, 4, vg.Length(1)},
	}

	for _, test := range tests {
		b.MinZ, b.MaxZ = test.minz, test.maxz
		if r := b.radius(test.z); r != test.r {
			t.Errorf("Got incorrect radius (%g) on %v", r, test)
		}
	}
}
