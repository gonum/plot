package plot

import "testing"

func TestAxisSmallTick(t *testing.T) {
	d := DefaultTicks{}
	for _, test := range []struct {
		Min, Max float64
		Labels   []string
	}{
		{
			Min:    -1.9846500878911073,
			Max:    0.4370974820125605,
			Labels: []string{"-1.6", "-0.8", "0"},
		},
	} {
		ticks := d.Ticks(test.Min, test.Max)
		var count int
		for _, tick := range ticks {
			if tick.Label != "" {
				if test.Labels[count] != tick.Label {
					t.Error("Ticks mismatch: Want", test.Labels[count], ", got", tick.Label)
				}
				count++
			}
		}
		if count != len(test.Labels) {
			t.Errorf("Too many tick labels")
		}
	}
}
