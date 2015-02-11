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
// go run generate_palettes < infile.tsv
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	delim string
	hex   bool
)

func init() {
	flag.StringVar(&delim, "d", "\t", "indicate field delimiter of input")
	flag.BoolVar(&hex, "hex", true, "indicate color values output in hex format")
	flag.Parse()
}

func mustAtoi(f string) byte {
	i, err := strconv.Atoi(f)
	if err != nil {
		panic(err)
	}
	if i < 0 || i > 0xff {
		panic(fmt.Sprintf("byte out of range", i))
	}
	return byte(i)
}

func main() {
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
	label := make(map[string]int)
	scanner := bufio.NewScanner(os.Stdin)
	var (
		lastType string

		last = make(map[string]string)

		defBuf = map[string]*bytes.Buffer{
			"Qualitative": &bytes.Buffer{},
			"Sequential":  &bytes.Buffer{},
			"Diverging":   &bytes.Buffer{},
		}

		cls = map[string]string{"Diverging": "Diverging", "Qualitative": "NonDiverging", "Sequential": "NonDiverging"}

		lookBuf = make(map[string][]string)
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
			if len(fields) > label["SchemeType"] {
				if typ := fields[label["SchemeType"]]; len(typ) != 0 {
					lastType = typ
				}
			}
			if name != last[lastType] {
				lookBuf[lastType] = append(lookBuf[lastType], fmt.Sprintf("%q: %s", name, name))
				if last[lastType] != "" {
					fmt.Fprintf(defBuf[lastType], "\t\t\t},\n\t\t},\n\t}\n")
				}
				fmt.Fprintf(defBuf[lastType], "\t%s %s = %s{\n", fields[label["ColorName"]], lastType, lastType)
				last[lastType] = name
			} else {
				fmt.Fprintf(defBuf[lastType], "\t\t\t},\n\t\t},\n")
			}
			fmt.Fprintf(defBuf[lastType], "\t\t%d: %sPalette{\n\t\t\tID: %q,\n\t\t\tColor: []color.Color{\n",
				mustAtoi(fields[label["NumOfColors"]]), cls[lastType], last[lastType])
		}
		values := []interface{}{
			fields[label["ColorLetter"]],
			mustAtoi(fields[label["R"]]),
			mustAtoi(fields[label["G"]]),
			mustAtoi(fields[label["B"]]),
		}
		if hex {
			fmt.Fprintf(defBuf[lastType], "\t\t\t\tColor{'%s', color.RGBA{0x%02x, 0x%02x, 0x%02x, 0xff}},\n", values...)
		} else {
			fmt.Fprintf(defBuf[lastType], "\t\t\t\tColor{'%s', color.RGBA{0x%02x, 0x%02x, 0x%02x, 0xff}},\n", values...)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}
	for _, typ := range []string{"Diverging", "Qualitative", "Sequential"} {
		fmt.Printf("var (\n%s\t\t\t},\n\t\t},\n\t}\n)\n", defBuf[typ].Bytes())
	}
	fmt.Println("\nvar (")
	for _, typ := range []string{"Diverging", "Qualitative", "Sequential"} {
		fmt.Printf("\t%s = map[string]%s{\n\t\t%v,\n\t}\n", strings.ToLower(typ), typ, strings.Join(lookBuf[typ], ",\n\t\t"))
	}
	fmt.Println("\tall = map[string]interface{}{")
	for _, typ := range []string{"Diverging", "Qualitative", "Sequential"} {
		fmt.Printf("\t\t%v,\n", strings.Join(lookBuf[typ], ",\n\t\t"))
	}
	fmt.Println("\t}\n)")
}
