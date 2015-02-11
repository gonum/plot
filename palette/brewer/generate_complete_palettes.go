// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// This program generates a Brewer Palette Go source file from
// a csv/tsv file exported from the xls file available from
// http://www.personal.psu.edu/cab38/ColorBrewer/ColorBrewer_updates.html
//
// Run the program:
// go run generate_complete_palettes infile.xml infile.tsv
package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"image/color"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var delim string

func init() {
	flag.StringVar(&delim, "d", "\t", "indicate field delimiter of input")
	flag.Parse()
}

type Colors struct {
	Sequential  []Scheme `xml:"sequential>scheme"`
	Diverging   []Scheme `xml:"diverging>scheme"`
	Qualitative []Scheme `xml:"categorical>scheme"`
}

type Scheme struct {
	ID    string  `xml:"id,attr"`
	Name  string  `xml:"name"`
	Class []Class `xml:"class"`
}

type Class struct {
	Laptop     *string `xml:"laptop"`
	CRT        *string `xml:"crt"`
	ColorBlind *string `xml:"eye"`
	Copy       *string `xml:"copy"`
	Projector  *string `xml:"projector"`
	Color      []Color `xml:"color"`
}

type Color struct {
	RGB  string `xml:"rgb"`
	Hex  string `xml:"hex"`
	CMYK string `xml:"cmyk"`
}

var (
	usability = map[string]string{"x": "Bad", "q": "Unsure", "": "Good"}
	class     = map[string]string{"Diverging": "Diverging", "Qualitative": "NonDiverging", "Sequential": "NonDiverging"}
)

func mustAtoi(f string, base int) byte {
	i, err := strconv.ParseUint(f, base, 8)
	if err != nil {
		panic(err)
	}
	return byte(i)
}

func getLetters(f string) map[string]map[color.RGBA]byte {
	letters := make(map[string]map[color.RGBA]byte)

	lf, err := os.Open(f)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	label := make(map[string]int)
	scanner := bufio.NewScanner(lf)
	var (
		lastType string

		last = make(map[string]string)
	)

	for scanner.Scan() {
		line := scanner.Text()
		if len(strings.TrimSpace(line)) == 0 {
			break
		}
		fields := strings.Split(line, delim)
		if fields[0] == "ColorName" {
			for i, f := range fields {
				label[f] = i
			}
			continue
		}
		if name := fields[label["ColorName"]]; len(name) != 0 {
			l, ok := letters[name]
			if !ok {
				l = make(map[color.RGBA]byte)
				letters[name] = l
			}
			if len(fields) > label["SchemeType"] {
				if typ := fields[label["SchemeType"]]; len(typ) != 0 {
					lastType = typ
				}
			}
			if name != last[lastType] {
				last[lastType] = name
			}
		}
		letters[last[lastType]][color.RGBA{
			R: mustAtoi(fields[label["R"]], 10),
			G: mustAtoi(fields[label["G"]], 10),
			B: mustAtoi(fields[label["B"]], 10),
		}] = fields[label["ColorLetter"]][0]
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}

	return letters
}

func main() {
	var cols Colors
	xf, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
	dec := xml.NewDecoder(xf)
	err = dec.Decode(&cols)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	letters := getLetters(os.Args[2])
	fmt.Println(`// Apache-Style Software License for ColorBrewer software and ColorBrewer Color Schemes
//
// Copyright ©2002 Cynthia Brewer, Mark Harrower, and The Pennsylvania State University.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.
// Go implementation Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Go port Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Go port Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Palette Copyright ©2002 Cynthia Brewer, Mark Harrower, and The Pennsylvania State University.

package brewer

import (
	"image/color"
)
`)
	cv := reflect.ValueOf(cols)
	ct := cv.Type()
	for i := 0; i < cv.NumField(); i++ {
		f := cv.Field(i)
		n := ct.Field(i).Name
		fmt.Println("var (")
		for _, schm := range f.Interface().([]Scheme) {
			fmt.Printf("\t%s = %s{\n", schm.ID, n, n)
			for _, cls := range schm.Class {
				fmt.Printf("\t\t%d: %sPalette{\n\t\t\tID: %q,\n\t\t\tName: %q,\n",
					len(cls.Color), class[n],
					schm.ID, schm.Name,
				)
				clv := reflect.ValueOf(cls)
				clt := clv.Type()
				for j := 0; j < clv.NumField()-1; j++ {
					clvf := clv.Field(j)
					if clvf.IsNil() {
						continue
					}
					fmt.Printf("\t\t\t%s: %s,\n", clt.Field(j).Name, usability[*clvf.Interface().(*string)])
				}
				fmt.Println("\t\t\tColor: []color.Color{")
				for _, col := range cls.Color {
					r, g, b := mustAtoi(col.Hex[2:4], 16), mustAtoi(col.Hex[4:6], 16), mustAtoi(col.Hex[6:8], 16)
					fmt.Printf("\t\t\t\tColor{%q, color.RGBA{0x%02x, 0x%02x, 0x%02x, 0xff}},\n",
						letters[schm.ID][color.RGBA{r, g, b, 0}], r, g, b,
					)
				}
				fmt.Println("\t\t\t},\n\t\t},")
			}
			fmt.Println("\t}")
		}
		fmt.Println(")")
	}
	fmt.Println("var (")
	for i := 0; i < cv.NumField(); i++ {
		f := cv.Field(i)
		n := ct.Field(i).Name
		nl := strings.ToLower(n)
		fmt.Printf("\t%s = map[string]%s{\n", nl, n)
		for _, schm := range f.Interface().([]Scheme) {
			fmt.Printf("\t\t%q: %s,\n", schm.ID, schm.ID)
		}
		fmt.Println("\t}")
	}
	fmt.Println("\tall = map[string]interface{}{")
	for i := 0; i < cv.NumField(); i++ {
		f := cv.Field(i)
		for _, schm := range f.Interface().([]Scheme) {
			fmt.Printf("\t\t%q: %s,\n", schm.ID, schm.ID)
		}
	}
	fmt.Println("\t}")
	fmt.Println(")")
}
