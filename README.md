# Gonum Plot  [![Build Status](https://travis-ci.org/gonum/plot.svg?branch=master)](https://travis-ci.org/gonum/plot) [![Coverage Status](https://coveralls.io/repos/gonum/plot/badge.svg?branch=master&service=github)](https://coveralls.io/github/gonum/plot?branch=master) [![GoDoc](https://godoc.org/gonum.org/v1/plot?status.svg)](https://godoc.org/gonum.org/v1/plot)

`gonum/plot` is the new, official fork of code.google.com/p/plotinum.
It provides an API for building and drawing plots in Go.
*Note* that this new API is still in flux and may change.
See the wiki for some [example plots](http://github.com/gonum/plot/wiki/Example-plots).

For additional Plotters, see the [Community Plotters](https://github.com/gonum/plot/wiki/Community-Plotters) Wiki page.

There is a discussion list on Google Groups: gonum-dev@googlegroups.com.

`gonum/plot` is split into following packages:

* `plot` provides a simple interface for laying out a plot and primitives for drawing to it.
* `plotter` provides a standard set of `Plotter`s which use the primitives provided by `plot` for drawing lines, scatter plots, box plots, error bars, etc. to a plot. You do not need to apply `plotter` for the use of `gonum/plot`. In addition, see the wiki for a tutorial on making your own custom plotters.
* `plotutil` contains routines that allow some common plot types to be created easily. However, `plotutil` is a new package which is not fully tested and is subject to change.
* `vg` provides a generic vector graphics API on top of other vector graphics back-ends such as a custom EPS back-end, draw2d, SVGo, X-Window and gopdf.

## Documentation

Documentation is available at:

  https://godoc.org/gonum.org/v1/plot

## Installation

You can get `gonum/plot` using go get:

`go get gonum.org/v1/plot/...`

If you write a cool plotter that you think others may be interested in using, please post to the list so that we can link to it in the `gonum/plot` wiki or possibly integrate it into the `plotter` package.
