// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package palette

import (
	"image/color"
	"log"
	"reflect"
	"testing"

	"github.com/gonum/plot"
	"github.com/gonum/plot/internal/cmpimg"
	"github.com/gonum/plot/palette/moreland"
	"github.com/gonum/plot/vg"
)

func ExampleStringMap_Legend_vertical() {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	plte := moreland.ExtendedBlackBody().Palette(5)
	cm := StringMap{
		Categories: []string{"Cat 1", "Cat 2", "Cat 3", "Cat 4", "Cat 5"},
		Colors:     plte.Colors(),
	}
	l := cm.Legend(6 * vg.Millimeter)
	p.Add(l)
	cm.SetupPlot(l, p)

	if err = p.Save(45, 100, "testdata/stringMapLegendVertical.png"); err != nil {
		log.Panic(err)
	}
}

func TestStringMap_Legend_vertical(t *testing.T) {
	cmpimg.CheckPlot(ExampleStringMap_Legend_vertical, t, "stringMapLegendVertical.png")
}

func ExampleStringMap_Legend_horizontal() {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	plte := moreland.ExtendedBlackBody().Palette(5)
	cm := StringMap{
		Categories: []string{"Cat 1", "Cat 2", "Cat 3", "Cat 4", "Cat 5"},
		Colors:     plte.Colors(),
	}
	l := cm.Legend(6 * vg.Millimeter)
	l.Horizontal = true
	p.Add(l)
	cm.SetupPlot(l, p)

	if err = p.Save(150, 25, "testdata/stringMapLegendHorizontal.png"); err != nil {
		log.Panic(err)
	}
}

func TestStringMap_Legend_horizontal(t *testing.T) {
	cmpimg.CheckPlot(ExampleStringMap_Legend_horizontal, t, "stringMapLegendHorizontal.png")
}

func ExampleIntMap_Legend_vertical() {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	plte := moreland.ExtendedBlackBody().Palette(5)
	cm := IntMap{
		Categories: []int{1, 3, 5, 7, 9},
		Colors:     plte.Colors(),
	}
	l := cm.Legend(6 * vg.Millimeter)
	p.Add(l)
	cm.SetupPlot(l, p)

	if err = p.Save(25, 100, "testdata/intMapLegendVertical.png"); err != nil {
		log.Panic(err)
	}
}

func TestIntMap_Legend_vertical(t *testing.T) {
	cmpimg.CheckPlot(ExampleIntMap_Legend_vertical, t, "intMapLegendVertical.png")
}

func ExampleIntMap_Legend_horizontal() {
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	plte := moreland.ExtendedBlackBody().Palette(5)
	cm := IntMap{
		Categories: []int{1, 3, 5, 7, 9},
		Colors:     plte.Colors(),
	}
	l := cm.Legend(6 * vg.Millimeter)
	l.Horizontal = true
	p.Add(l)
	cm.SetupPlot(l, p)

	if err = p.Save(100, 25, "testdata/intMapLegendHorizontal.png"); err != nil {
		log.Panic(err)
	}
}

func TestIntMap_Legend_horizontal(t *testing.T) {
	cmpimg.CheckPlot(ExampleIntMap_Legend_horizontal, t, "intMapLegendHorizontal.png")
}

func TestStringMap_At(t *testing.T) {
	plte := palette{color.Black, color.White}
	cm := &StringMap{
		Categories: []string{"Black", "White"},
		Colors:     plte.Colors(),
	}
	c, err := cm.At("Black")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(c, color.Black) {
		t.Errorf("color should be black but is %v", c)
	}
	c, err = cm.At("White")
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(c, color.White) {
		t.Errorf("color should be white but is %v", c)
	}
	c, err = cm.At("black")
	if err == nil {
		t.Error("this should cause an error but doesn't")
	}
}

func TestIntMap_At(t *testing.T) {
	plte := palette{color.Black, color.White}
	cm := &IntMap{
		Categories: []int{6, 8},
		Colors:     plte.Colors(),
	}
	c, err := cm.At(6)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(c, color.Black) {
		t.Errorf("color should be black but is %v", c)
	}
	c, err = cm.At(8)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(c, color.White) {
		t.Errorf("color should be white but is %v", c)
	}
	c, err = cm.At(0)
	if err == nil {
		t.Error("this should cause an error but doesn't")
	}
}
