package plotutil

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/vg"
	"fmt"
)

// AddBoxPlots adds box plot plotters to a plot and
// sets the X axis of the plot to be nominal.
// The variadic arguments must be either strings
// or plotter.Valuers.  Each valuer adds a box plot
// to the plot at the X location corresponding to
// the number of box plots added before it.  If a
// plotter.Valuer is immediately preceeded by a
// string then the string value is used to label the
// tick mark for the box plot's X location.
func AddBoxPlots(plt *plot.Plot, width vg.Length, vs ...interface{}) {
	var names []string
	name := ""
	for _, v := range vs {
		switch t := v.(type) {
		case string:
			name = t

		case plotter.Valuer:
			plt.Add(plotter.NewBoxPlot(width, float64(len(names)), t))
			names = append(names, name)
			name = ""

		default:
			panic(fmt.Sprintf("AddScatters handles strings and plotter.XYers, got %T", t))
		}
	}
	plt.NominalX(names...)
}

// AddScatters adds Scatter plotters to a plot.
// The variadic arguments must be either strings
// or plotter.XYers.  Each plotter.XYer is added to
// the plot using the next color, and glyph shape
// via the Color and Shape functions. If a
// plotter.XYer is immediately preceeded by
// a string then a legend entry is added to the plot
// using the string as the name.
func AddScatters(plt *plot.Plot, vs ...interface{}) {
	name := ""
	var i int
	for _, v := range vs {
		switch t := v.(type) {
		case string:
			name = t

		case plotter.XYer:
			s := plotter.NewScatter(t)
			s.Color = Color(i)
			s.Shape = Shape(i)
			i++
			plt.Add(s)
			if name != "" {
				plt.Legend.Add(name, s)
				name = ""
			}

		default:
			panic(fmt.Sprintf("AddScatters handles strings and plotter.XYers, got %T", t))
		}
	}
}

// AddLines adds Line plotters to a plot.
// The variadic arguments must be either strings
// or plotter.XYers.  Each plotter.XYer is added to
// the plot using the next color and dashes
// shape via the Color and Dashes functions.
// If a plotter.XYer is immediately preceeded by
// a string then a legend entry is added to the plot
// using the string as the name.
func AddLines(plt *plot.Plot, vs ...interface{}) {
	name := ""
	var i int
	for _, v := range vs {
		switch t := v.(type) {
		case string:
			name = t

		case plotter.XYer:
			l := plotter.NewLine(t)
			l.Color = Color(i)
			l.Dashes = Dashes(i)
			i++
			plt.Add(l)
			if name != "" {
				plt.Legend.Add(name, l)
				name = ""
			}

		default:
			panic(fmt.Sprintf("AddLines handles strings and plotter.XYers, got %T", t))
		}
	}
}

// AddLinePoints adds Line and Scatter plotters to a
// plot.  The variadic arguments must be either strings
// or plotter.XYers.  Each plotter.XYer is added to
// the plot using the next color, dashes, and glyph
// shape via the Color, Dashes, and Shape functions.
// If a plotter.XYer is immediately preceeded by
// a string then a legend entry is added to the plot
// using the string as the name.
func AddLinePoints(plt *plot.Plot, vs ...interface{}) {
	name := ""
	var i int
	for _, v := range vs {
		switch t := v.(type) {
		case string:
			name = t

		case plotter.XYer:
			l, s := plotter.NewLinePoints(t)
			l.Color = Color(i)
			l.Dashes = Dashes(i)
			s.Color = Color(i)
			s.Shape = Shape(i)
			i++
			plt.Add(l, s)
			if name != "" {
				plt.Legend.Add(name, l, s)
				name = ""
			}

		default:
			panic(fmt.Sprintf("AddLinePoints handles strings and plotter.XYers, got %T", t))
		}
	}
}
