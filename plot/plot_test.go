// Copyright 2012 The Plotinum Authors. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package plot

import (
	"testing"
)

//func TestDrawImage(t *testing.T) {
//	if err := Example().Save(4, 4, "test.png"); err != nil {
//		t.Error(err)
//	}
//}

func TestDrawEps(t *testing.T) {
	if err := Example().Save(4, 4, "test.eps"); err != nil {
		t.Error(err)
	}
}

func TestDrawSvg(t *testing.T) {
	if err := Example().Save(4, 4, "test.svg"); err != nil {
		t.Error(err)
	}
}
