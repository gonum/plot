// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package moreland

import (
	"fmt"
	"math/rand"
	"testing"
)

func BenchmarkLuminance_At(b *testing.B) {
	pBase := ExtendedBlackBody()
	for n := 2; n < 12; n += 2 {
		p, err := NewLuminance(pBase.Palette(n).Colors())
		if err != nil {
			b.Fatal(err)
		}
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
