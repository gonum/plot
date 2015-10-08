// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plotter

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gonum/plot"
)

var generateTestData = flag.Bool("regen", false, "Uses the current state to regenerate the test data.")

const (
	width  = 100
	height = 100
)

func checkPlot(dir, name, ext string, p *plot.Plot) {
	filename := filepath.Join(dir, fmt.Sprintf("%s.%s", name, ext))

	c, err := p.WriterTo(width, height, ext)
	handle(err)

	var buf bytes.Buffer
	_, err = c.WriteTo(&buf)
	handle(err)

	// Recreate Golden images.
	if *generateTestData {
		handle(p.Save(width, height, filename))
	}

	f, err := os.Open(filename)
	handle(err)

	want, err := ioutil.ReadAll(f)
	handle(err)
	f.Close()
	if !bytes.Equal(buf.Bytes(), want) {
		fmt.Printf("image mismatch for %s\n", filename)
		return
	}
	fmt.Printf("Image saved in dir %s as %s.%s. "+
		"Normally, you would use plot.Save().\n", dir, name, ext)
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}
