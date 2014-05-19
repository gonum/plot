package plotutil

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/vg"
	"errors"
	"fmt"
)

// AddStackedAreaPlots adds stacked area plot plotters to
// a plot.
// The variadic arguments must be either strings
// or plotter.Valuers.  Each valuer adds a stacked area
// plot to the plot above the stacked area plots added
// before it.  If a plotter.Valuer is immediately
// preceeded by a string then the string value is used to
// label the legend.
//
// If an error occurs then none of the plotters are added
// to the plot, and the error is returned.
func AddStackedAreaPlots(plt *plot.Plot, xs plotter.Valuer, vs ...interface{}) error {
	var names []string
	name := ""

	xys := make([]plotter.XYs, 0)

	for _, v := range vs {
		switch t := v.(type) {
		case string:
			name = t

		case plotter.Valuer:
			// Stack the data by adding the incoming data to the previous values
			stackedData, err := plotter.CreateXYs(xs, t)
			if err != nil {
				return err
			}
			if 0 != len(xys) {
				for i := range stackedData {
					stackedData[i].Y += xys[len(xys)-1][i].Y
				}
			}

			xys = append(xys, stackedData)

			names = append(names, name)
			name = ""

		default:
			panic(fmt.Sprintf("AddStackedAreaPlots handles strings and plotter.Valuers, got %T", t))
		}
	}

	numPlots := len(xys)
	if numPlots == 0 {
		return errors.New("No data has been added")
	}

	for i := numPlots - 1; i >= 0; i-- {
		// Make a line plotter and set its style.
		l, err := plotter.NewLine(xys[i])
		if err != nil {
			return err
		}

		l.LineStyle.Width = vg.Points(0)
		color := Color(i)
		l.ShadeColor = &color

		plt.Add(l)
		plt.Legend.Add(names[i], l)
	}

	plt.Legend.Top = true
	plt.Legend.Left = true

	return nil
}

// AddBoxPlots adds box plot plotters to a plot and
// sets the X axis of the plot to be nominal.
// The variadic arguments must be either strings
// or plotter.Valuers.  Each valuer adds a box plot
// to the plot at the X location corresponding to
// the number of box plots added before it.  If a
// plotter.Valuer is immediately preceeded by a
// string then the string value is used to label the
// tick mark for the box plot's X location.
//
// If an error occurs then none of the plotters are added
// to the plot, and the error is returned.
func AddBoxPlots(plt *plot.Plot, width vg.Length, vs ...interface{}) error {
	var ps []plot.Plotter
	var names []string
	name := ""
	for _, v := range vs {
		switch t := v.(type) {
		case string:
			name = t

		case plotter.Valuer:
			b, err := plotter.NewBoxPlot(width, float64(len(names)), t)
			if err != nil {
				return err
			}
			ps = append(ps, b)
			names = append(names, name)
			name = ""

		default:
			panic(fmt.Sprintf("AddBoxPlots handles strings and plotter.Valuers, got %T", t))
		}
	}
	plt.Add(ps...)
	plt.NominalX(names...)
	return nil
}

// AddScatters adds Scatter plotters to a plot.
// The variadic arguments must be either strings
// or plotter.XYers.  Each plotter.XYer is added to
// the plot using the next color, and glyph shape
// via the Color and Shape functions. If a
// plotter.XYer is immediately preceeded by
// a string then a legend entry is added to the plot
// using the string as the name.
//
// If an error occurs then none of the plotters are added
// to the plot, and the error is returned.
func AddScatters(plt *plot.Plot, vs ...interface{}) error {
	var ps []plot.Plotter
	names := make(map[*plotter.Scatter]string)
	name := ""
	var i int
	for _, v := range vs {
		switch t := v.(type) {
		case string:
			name = t

		case plotter.XYer:
			s, err := plotter.NewScatter(t)
			if err != nil {
				return err
			}
			s.Color = Color(i)
			s.Shape = Shape(i)
			i++
			ps = append(ps, s)
			if name != "" {
				names[s] = name
				name = ""
			}

		default:
			panic(fmt.Sprintf("AddScatters handles strings and plotter.XYers, got %T", t))
		}
	}
	plt.Add(ps...)
	for p, n := range names {
		plt.Legend.Add(n, p)
	}
	return nil
}

// AddLines adds Line plotters to a plot.
// The variadic arguments must be either strings
// or plotter.XYers.  Each plotter.XYer is added to
// the plot using the next color and dashes
// shape via the Color and Dashes functions.
// If a plotter.XYer is immediately preceeded by
// a string then a legend entry is added to the plot
// using the string as the name.
//
// If an error occurs then none of the plotters are added
// to the plot, and the error is returned.
func AddLines(plt *plot.Plot, vs ...interface{}) error {
	var ps []plot.Plotter
	names := make(map[*plotter.Line]string)
	name := ""
	var i int
	for _, v := range vs {
		switch t := v.(type) {
		case string:
			name = t

		case plotter.XYer:
			l, err := plotter.NewLine(t)
			if err != nil {
				return err
			}
			l.Color = Color(i)
			l.Dashes = Dashes(i)
			i++
			ps = append(ps, l)
			if name != "" {
				names[l] = name
				name = ""
			}

		default:
			panic(fmt.Sprintf("AddLines handles strings and plotter.XYers, got %T", t))
		}
	}
	plt.Add(ps...)
	for p, n := range names {
		plt.Legend.Add(n, p)
	}
	return nil
}

// AddLinePoints adds Line and Scatter plotters to a
// plot.  The variadic arguments must be either strings
// or plotter.XYers.  Each plotter.XYer is added to
// the plot using the next color, dashes, and glyph
// shape via the Color, Dashes, and Shape functions.
// If a plotter.XYer is immediately preceeded by
// a string then a legend entry is added to the plot
// using the string as the name.
//
// If an error occurs then none of the plotters are added
// to the plot, and the error is returned.
func AddLinePoints(plt *plot.Plot, vs ...interface{}) error {
	var ps []plot.Plotter
	names := make(map[[2]plot.Thumbnailer]string)
	name := ""
	var i int
	for _, v := range vs {
		switch t := v.(type) {
		case string:
			name = t

		case plotter.XYer:
			l, s, err := plotter.NewLinePoints(t)
			if err != nil {
				return err
			}
			l.Color = Color(i)
			l.Dashes = Dashes(i)
			s.Color = Color(i)
			s.Shape = Shape(i)
			i++
			ps = append(ps, l, s)
			if name != "" {
				names[[2]plot.Thumbnailer{l, s}] = name
				name = ""
			}

		default:
			panic(fmt.Sprintf("AddLinePoints handles strings and plotter.XYers, got %T", t))
		}
	}
	plt.Add(ps...)
	for ps, n := range names {
		plt.Legend.Add(n, ps[0], ps[1])
	}
	return nil
}

// AddErrorBars adds XErrorBars and YErrorBars
// to a plot.  The variadic arguments must be
// of type plotter.XYer, and must be either a
// plotter.XErrorer, plotter.YErrorer, or both.
// Each errorer is added to the plot the color from
// the Colors function corresponding to its position
// in the argument list.
//
// If an error occurs then none of the plotters are added
// to the plot, and the error is returned.
func AddErrorBars(plt *plot.Plot, vs ...interface{}) error {
	var ps []plot.Plotter
	for i, v := range vs {
		added := false

		if xerr, ok := v.(interface {
			plotter.XYer
			plotter.XErrorer
		}); ok {
			e, err := plotter.NewXErrorBars(xerr)
			if err != nil {
				return err
			}
			e.Color = Color(i)
			ps = append(ps, e)
			added = true
		}

		if yerr, ok := v.(interface {
			plotter.XYer
			plotter.YErrorer
		}); ok {
			e, err := plotter.NewYErrorBars(yerr)
			if err != nil {
				return err
			}
			e.Color = Color(i)
			ps = append(ps, e)
			added = true
		}

		if added {
			continue
		}
		panic(fmt.Sprintf("AddErrorBars expects plotter.XErrorer or plotter.YErrorer, got %T", v))
	}
	plt.Add(ps...)
	return nil
}

// AddXErrorBars adds XErrorBars to a plot.
// The variadic arguments must be
// of type plotter.XYer, and plotter.XErrorer.
// Each errorer is added to the plot the color from
// the Colors function corresponding to its position
// in the argument list.
//
// If an error occurs then none of the plotters are added
// to the plot, and the error is returned.
func AddXErrorBars(plt *plot.Plot, es ...interface {
	plotter.XYer
	plotter.XErrorer
}) error {
	var ps []plot.Plotter
	for i, e := range es {
		bars, err := plotter.NewXErrorBars(e)
		if err != nil {
			return err
		}
		bars.Color = Color(i)
		ps = append(ps, bars)
	}
	plt.Add(ps...)
	return nil
}

// AddYErrorBars adds YErrorBars to a plot.
// The variadic arguments must be
// of type plotter.XYer, and plotter.YErrorer.
// Each errorer is added to the plot the color from
// the Colors function corresponding to its position
// in the argument list.
//
// If an error occurs then none of the plotters are added
// to the plot, and the error is returned.
func AddYErrorBars(plt *plot.Plot, es ...interface {
	plotter.XYer
	plotter.YErrorer
}) error {
	var ps []plot.Plotter
	for i, e := range es {
		bars, err := plotter.NewYErrorBars(e)
		if err != nil {
			return err
		}
		bars.Color = Color(i)
		ps = append(ps, bars)
	}
	plt.Add(ps...)
	return nil
}
