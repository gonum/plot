// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"log"
	"math"
	"math/rand"
	"testing"
	"os"
	"fmt"
	"crypto/sha1"
	"io"
	"encoding/base64"
	"io/ioutil"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/internal/cmpimg"
	"gonum.org/v1/plot/palette/moreland"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"

)

// ExampleScatter_color draws a colored scatter plot.
// Each point is plotted with a different color depending on
// external criteria.
func ExampleScatter_color() {
	rnd := rand.New(rand.NewSource(1))

	// randomTriples returns some random but correlated x, y, z triples
	randomTriples := func(n int) XYZs {
		data := make(XYZs, n)
		for i := range data {
			if i == 0 {
				data[i].X = rnd.Float64()
			} else {
				data[i].X = data[i-1].X + 2*rnd.Float64()
			}
			data[i].Y = data[i].X + 10*rnd.Float64()
			data[i].Z = data[i].X
		}
		return data
	}

	n := 15
	scatterData := randomTriples(n)

	// Calculate the range of Z values.
	minZ, maxZ := math.Inf(1), math.Inf(-1)
	for _, xyz := range scatterData {
		if xyz.Z > maxZ {
			maxZ = xyz.Z
		}
		if xyz.Z < minZ {
			minZ = xyz.Z
		}
	}

	colors := moreland.Kindlmann() // Initialize a color map.
	colors.SetMax(maxZ)
	colors.SetMin(minZ)

	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Colored Points Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(NewGrid())

	l := &ColorBar{ColorMap: colors}
	p.Add(l)

	sc, err := NewScatter(scatterData)
	if err != nil {
		log.Panic(err)
	}

	// Specify style and color for individual points.
	sc.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		_, _, z := scatterData.XYZ(i)
		d := (z - minZ) / (maxZ - minZ)
		rng := maxZ - minZ
		k := d*rng + minZ
		c, err := colors.At(k)
		if err != nil {
			log.Panic(err)
		}
		return draw.GlyphStyle{Color: c, Radius: vg.Points(3), Shape: draw.CircleGlyph{}}
	}
	p.Add(sc)

	err = p.Save(300, 230, "testdata/scatterColor.png")
	if err != nil {
		log.Panic(err)
	}
}

func TestScatterColor(t *testing.T) {
	cmpimg.CheckPlot(ExampleScatter_color, t, "scatterColor.png")

	statistics, err := os.Stat("testdata/scatterColor.png")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("The size of scatterColor.png is %d\n", statistics.Size())

	statistics_gold, err := os.Stat("testdata/scatterColor_golden.png")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("The size of scatterColor_golden.png is %d\n", statistics_gold.Size())

	f, err := os.Open("testdata/scatterColor.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("SHA1 checksum for scatterColor.png is %x\n", h.Sum(nil))

	fGold, err := os.Open("testdata/scatterColor_golden.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fGold.Close()
	hGgold := sha1.New()
	if _, err := io.Copy(hGgold, fGold); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("SHA1 checksum for scatterColor_golden.png is %x\n", hGgold.Sum(nil))

	fStr, err := ioutil.ReadFile("testdata/scatterColor.png")
	if err != nil {
		fmt.Print(err)
	}
	b64fStr := base64.StdEncoding.EncodeToString([]byte(fStr))
	fmt.Println ("ENCODING OF SCATTERCOLOR STARTED\n", b64fStr,"\nENCODING OF SCATTERCOLOR ENDED")

	fStrGold, err := ioutil.ReadFile("testdata/scatterColor_golden.png")
	if err != nil {
		fmt.Print(err)
	}
	b64fStrGold := base64.StdEncoding.EncodeToString([]byte(fStrGold))
	fmt.Println ("ENCODING OF SCATTERCOLOR_GOLDEN STARTED\n", b64fStrGold,"\nENCODING OF SCATTERCOLOR_GOLDEN ENDED")
}


