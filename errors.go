// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import "fmt"

// Errors is a collection of plotting errors. An Errors value
// must have a length of at least one.
type Errors []error

func (e Errors) Error() string {
	switch len(e) {
	case 0:
		panic("plot: invalid error")
	case 1:
		return e[0].Error()
	case 2:
		return fmt.Sprintf("plot: %v and 1 more error", e[0])
	default:
		return fmt.Sprintf("plot: %v and %d more errors", e[0], len(e)-1)
	}
}
