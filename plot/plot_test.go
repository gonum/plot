package plot

import (
	"testing"
)

func TestDrawImage(t *testing.T) {
	if err := Example().Save(4, 4, "test.png"); err != nil {
		t.Error(err)
	}
}

func TestDrawEps(t *testing.T) {
	if err := Example().Save(4, 4, "test.eps"); err != nil {
		t.Error(err)
	}
}
