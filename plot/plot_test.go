package plot

import (
	"testing"
)

func TestDrawImage(t *testing.T) {
	Example().Save(4, 4, "test.png")
}

func TestDrawEps(t *testing.T) {
	Example().Save(4, 4, "test.eps")
}
