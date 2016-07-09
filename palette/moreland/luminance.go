// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package moreland

import (
	"fmt"
	"image/color"
	"math"

	"github.com/gonum/plot/palette"
)

// luminance is a color palette that interpolates
// between control colors in a way that ensures a linear relationship
// between the luminance of a color and the value it represents.
type luminance struct {
	// colors are the control colors to be interpolated among.
	// The colors must be monotonically increasing in luminance.
	colors []cieLAB

	// scalars are the scalar control points associated with
	// each item in colors (above). They are monotonically
	// increasing values between zero and one that correspond
	// to the luminance of a given control color in relation
	// to the minimum and maximum luminance among all control
	// colors.
	scalars []float64

	// alpha represents the opacity of the returned
	// colors in the range (0,1). It is set to 1 by default.
	alpha float64

	// min and max are the minimum and maximum values of the range of scalars
	// that can be mapped to colors using this ColorMap.
	min, max float64
}

// NewLuminance creates a new Luminance ColorMap from the given controlColors.
// luminance is a color palette that interpolates
// between control colors in a way that ensures a linear relationship
// between the luminance of a color and the value it represents.
// If the luminance of the controls is not monotonically increasing, an
// error will be returned.
func NewLuminance(controls []color.Color) (palette.ColorMap, error) {
	l := luminance{
		colors:  make([]cieLAB, len(controls)),
		scalars: make([]float64, len(controls)),
		alpha:   1,
	}
	max := math.Inf(-1)
	min := math.Inf(1)
	for i, c := range controls {
		lab := colorTosRGBA(c).cieLAB()
		l.colors[i] = lab
		max = math.Max(max, lab.L)
		min = math.Min(min, lab.L)
		if i > 0 && lab.L <= l.colors[i-1].L {
			return nil, fmt.Errorf("moreland: luminance of color %d (%g) is not "+
				"greater than that of color %d (%g)", i, lab.L, i-1, l.colors[i-1].L)
		}
	}
	// Normalize scalar values to the range (0,1).
	rnge := max - min
	for i, c := range l.colors {
		l.scalars[i] = (c.L - min) / rnge
	}
	// Sometimes the first and last scalars do not end up
	// being exactly zero and one owing to the imperfect
	// precision of floating point operations.
	// Here we set them to exactly zero and one to avoid
	// the possibility of the At() function returning
	// an out-of-range error for values that actually
	// should be in the range.
	l.scalars[0] = 0
	l.scalars[len(l.scalars)-1] = 1
	return &l, nil
}

// At implements the palette.ColorMap interface for a luminance value.
func (l *luminance) At(v float64) (color.Color, error) {
	if err := checkRange(l.min, l.max, v); err != nil {
		return nil, err
	}
	scalar := (v - l.min) / l.max
	if !inUnitRange(scalar) {
		return nil, fmt.Errorf("moreland: interpolation value (%g) out of range [%g,%g]", scalar, l.min, l.max)
	}
	i := searchFloat64s(l.scalars, scalar)
	if i == 0 {
		return l.colors[i].cieXYZ().rgb().sRGBA(l.alpha), nil
	}
	c1 := l.colors[i-1]
	c2 := l.colors[i]
	frac := (scalar - l.scalars[i-1]) / (l.scalars[i] - l.scalars[i-1])
	o := cieLAB{
		L: frac*(c2.L-c1.L) + c1.L,
		A: frac*(c2.A-c1.A) + c1.A,
		B: frac*(c2.B-c1.B) + c1.B,
	}.cieXYZ().rgb().sRGBA(l.alpha)
	o.clamp()
	return o, nil
}

func checkRange(min, max, val float64) error {
	if max == min {
		return fmt.Errorf("moreland: color map max == min == %g", max)
	}
	if min > max {
		return fmt.Errorf("moreland: color map max (%g) < min (%g)", max, min)
	}
	if val < min {
		return palette.ErrUnderflow
	}
	if val > max {
		return palette.ErrOverflow
	}
	if math.IsNaN(val) {
		return palette.ErrNaN
	}
	return nil
}

// searchFloat64s acts the same as sort.SearchFloat64s, except
// it uses a simple search algorithm instead of binary search.
func searchFloat64s(vals []float64, val float64) int {
	for j, v := range vals {
		if val <= v {
			return j
		}
	}
	return len(vals)
}

// SetMax implements the palette.ColorMap interface for a luminance value.
func (l *luminance) SetMax(v float64) {
	l.max = v
}

// SetMin implements the palette.ColorMap interface for a luminance value.
func (l *luminance) SetMin(v float64) {
	l.min = v
}

// Max implements the palette.ColorMap interface for a luminance value.
func (l *luminance) Max() float64 {
	return l.max
}

// Min implements the palette.ColorMap interface for a luminance value.
func (l *luminance) Min() float64 {
	return l.min
}

// SetAlpha sets the opacity value of this color map. Zero is transparent
// and one is completely opaque.
// The function will panic is alpha is not between zero and one.
func (l *luminance) SetAlpha(alpha float64) {
	if !inUnitRange(alpha) {
		panic(fmt.Errorf("moreland: invalid alpha: %g", alpha))
	}
	l.alpha = alpha
}

// Alpha returns the opacity value of this color map.
func (l *luminance) Alpha() float64 {
	return l.alpha
}

// Palette returns a value that fulfills the palette.Palette interface,
// where n is the number of desired colors.
func (l luminance) Palette(n int) palette.Palette {
	if l.Max() == 0 && l.Min() == 0 {
		l.SetMin(0)
		l.SetMax(1)
	}
	delta := (l.max - l.min) / float64(n-1)
	var v float64
	c := make([]color.Color, n)
	for i := 0; i < n; i++ {
		v = l.min + delta*float64(i)
		var err error
		c[i], err = l.At(v)
		if err != nil {
			panic(err)
		}
	}
	return plte(c)
}

// plte fulfils the palette.Palette interface.
type plte []color.Color

// Colors fulfils the palette.Palette interface.
func (p plte) Colors() []color.Color {
	return p
}

// BlackBody is a Luminance-class ColorMap based on the colors of black body radiation.
// Although the colors are inspired by the wavelengths of light from
// black body radiation, the actual colors used are designed to be
// perceptually uniform. Colors of the desired brightness and hue are chosen,
// and then the colors are adjusted such that the luminance is perceptually
// linear (according to the CIE LAB color space).
func BlackBody() palette.ColorMap {
	return &luminance{
		colors: []cieLAB{
			cieLAB{L: 0, A: 0, B: 0},
			cieLAB{L: 39.112572747719774, A: 55.92470934659227, B: 37.65159714510402},
			cieLAB{L: 58.45705480680232, A: 43.34389690857626, B: 65.95409116544081},
			cieLAB{L: 84.13253643355525, A: -6.459770854468639, B: 82.41994470228775},
			cieLAB{L: 100, A: 0, B: 0}},
		scalars: []float64{0, 0.39112572747719776, 0.5845705480680232, 0.8413253643355525, 1},
		alpha:   1,
	}
}

// ExtendedBlackBody is a Luminance-class ColorMap based on the colors of black body radiation
// with some blue and purple hues thrown in at the lower end to add some "color."
// The color map is similar to the default colors used in gnuplot. Colors of
// the desired brightness and hue are chosen, and then the colors are adjusted
// such that the luminance is perceptually linear (according to the CIE LAB
// color space).
func ExtendedBlackBody() palette.ColorMap {
	return &luminance{
		colors: []cieLAB{
			cieLAB{L: 0, A: 0, B: 0},
			cieLAB{L: 21.873483862751876, A: 50.19882295659109, B: -74.66982659778306},
			cieLAB{L: 34.506542513775905, A: 75.41302687474061, B: -88.73807072507786},
			cieLAB{L: 47.02980511087303, A: 70.93217189227919, B: 33.59880053746508},
			cieLAB{L: 65.17482203230537, A: 49.14591409658836, B: 56.86480950937553},
			cieLAB{L: 84.13253643355525, A: -6.459770854468639, B: 82.41994470228775},
			cieLAB{L: 100, A: 0, B: 0},
		},
		scalars: []float64{0, 0.21873483862751875, 0.34506542513775906, 0.4702980511087303,
			0.6517482203230537, 0.8413253643355525, 1},
		alpha: 1,
	}
}

// Kindlmann is a Luminance-class ColorMap that uses the colors
// first proposed in a paper
// by Kindlmann, Reinhard, and Creem. The map is basically the rainbow
// color map with the luminance adjusted such that it monotonically
// changes, making it much more perceptually viable.
//
// Citation:
// Gordon Kindlmann, Erik Reinhard, and Sarah Creem. 2002. Face-based
// luminance matching for perceptual colormap generation. In Proceedings
// of the conference on Visualization '02 (VIS '02). IEEE Computer Society,
// Washington, DC, USA, 299-306.
func Kindlmann() palette.ColorMap {
	return &luminance{
		colors: []cieLAB{
			cieLAB{L: 0, A: 0, B: 0},
			cieLAB{L: 10.479520542426698, A: 34.05557958902206, B: -34.21934877170809},
			cieLAB{L: 21.03011379005111, A: 52.30473571100955, B: -61.852601228346536},
			cieLAB{L: 31.03098927978494, A: 23.814976212074402, B: -57.73419358300511},
			cieLAB{L: 40.21480513626115, A: -24.858012706049536, B: -7.322176588219942},
			cieLAB{L: 52.73108089333358, A: -19.064976357731634, B: -25.558178073848147},
			cieLAB{L: 60.007326812392634, A: -61.75624590074585, B: 56.43522875191319},
			cieLAB{L: 69.81578343076002, A: -58.33353084882392, B: 68.37457857626646},
			cieLAB{L: 79.55703752324776, A: -22.50477758899383, B: 78.57946686200843},
			cieLAB{L: 89.818961593653, A: 7.586705160677109, B: 15.375961528833981},
			cieLAB{L: 100, A: 0, B: 0},
		},
		scalars: []float64{0, 0.10479520542426699, 0.2103011379005111, 0.3103098927978494,
			0.4021480513626115, 0.5273108089333358, 0.6000732681239264, 0.6981578343076003,
			0.7955703752324775, 0.89818961593653, 1},
		alpha: 1,
	}
}

// ExtendedKindlmann is a Luminance-class ColorMap uses the colors from
// Kindlmann but also
// adds more hues by doing a more than 360 degree loop around the hues.
// This works because the endpoints have low saturation and very
// different brightness.
func ExtendedKindlmann() palette.ColorMap {
	return &luminance{
		colors: []cieLAB{
			cieLAB{L: 0, A: 0, B: 0},
			cieLAB{L: 13.371291966477482, A: 40.39368469479174, B: -47.73239449160565},
			cieLAB{L: 25.072421338587574, A: -18.01441053740843, B: -5.313556572210176},
			cieLAB{L: 37.411516363056116, A: -43.058336774976055, B: 39.30203907343062},
			cieLAB{L: 49.75026355291354, A: -15.774050138318895, B: 53.507917567416094},
			cieLAB{L: 61.643756252245225, A: 52.67703578954919, B: 43.82595336046358},
			cieLAB{L: 74.93187540089825, A: 50.92061741619164, B: -30.235411697966242},
			cieLAB{L: 87.64732748562544, A: 14.355163639545697, B: -17.471161313826332},
			cieLAB{L: 100, A: 0, B: 0},
		},
		scalars: []float64{0, 0.13371291966477483, 0.25072421338587575, 0.37411516363056113,
			0.4975026355291354, 0.6164375625224523, 0.7493187540089825, 0.8764732748562544, 1},
		alpha: 1,
	}
}
