// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package palette

import (
	"fmt"
	"image/color"
	"math/rand"
	"testing"
)

func BenchmarkColorMap_At(b *testing.B) {
	pBase := FromPalette(New(color.White, color.Black))
	for n := 2; n < 22; n += 2 {
		p := FromPalette(pBase.Palette(n))
		p.SetMax(1)
		rand.Seed(1)
		b.Run(fmt.Sprintf("%d controls", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, err := p.At(rand.Float64()); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkIntMap_At(b *testing.B) {
	pBase := FromPalette(New(color.White, color.Black))
	for n := 2; n < 22; n += 2 {
		ints := make([]int, n)
		for i := 0; i < n; i++ {
			ints[i] = i
		}
		p := IntMap{
			Colors:     pBase.Palette(n).Colors(),
			Categories: ints,
		}
		b.ResetTimer()
		b.Run(fmt.Sprintf("%d colors", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, err := p.At(i % n); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkStringMap_At(b *testing.B) {
	pBase := FromPalette(New(color.White, color.Black))
	for n := 2; n < 22; n += 2 {
		strs := make([]string, n)
		for i := 0; i < n; i++ {
			strs[i] = fmt.Sprintf("%d", i)
		}
		p := StringMap{
			Colors:     pBase.Palette(n).Colors(),
			Categories: strs,
		}
		b.ResetTimer()
		b.Run(fmt.Sprintf("%d colors", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, err := p.At(fmt.Sprintf("%d", i%n)); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
