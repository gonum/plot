// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmpimg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const wantDiffEncoded = `iVBORw0KGgoAAAANSUhEUgAAAZAAAAEzEAIAAADAxR6YAAAHBklEQVR4nOzYjW3kRABA4RWiC0QXUMami2xRSRfrMu76gDKQZQ3+2907xEMg8X3SXZKxdzwZ29JTfrgAAJASWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAMYEFABATWAAAsR//7QX8X12v8//T9E/NvPje+R+tZjvP9882PjV/naZvf+J8lb96xb+/l+uar9d5zfPX8+rnY/PocvTZ9ebj4+j5TmyPjauNM8f8YyXrGh7NM2ZaPjXmHj8dZx4j+zUsax2zr0fHjGOG9frrOfs9WM/bjmx3bnxmP368/2MHzmtfj61j49/2+tv1bO/Ho1UdzzyO7898NsNx/NU+bHfkeMXZzz+Nnz4/3m/rPXi/He/B6r55Tt+evk3btZ73ffsGfH6MPfn8+PL1cvnydVnB779dLrfb/c/fa+zT9imaXhw9P3fH53H9fbdvwH6/lvdv/zSvT/B4e/dzL2e/7d7bZ98vOzqPvG3mWPZ42d375g1f5tzfgePI8Uk8P3GP1nAc+fWXeWT+f5rmOzH2dn5K9nt0/+b8j97ddYXnt+31vq2zHD+/3onz0XGH3m/P3tBnb8szr2Z5vdOvZzyv6Di+fDdNt9t5Dn/BAgCICSwAgJjAAgCICSwAgJjAAgCICSwAgJjAAgCICSwAgJjAAgCICSwAgJjAAgAAAOC/zV+wAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAICawAABiAgsAIPZHAAAA///KAlB2mCuyaAAAAABJRU5ErkJggg==`

func TestDiff(t *testing.T) {
	got, err := os.ReadFile(filepath.FromSlash("./testdata/failed_input.png"))
	if err != nil {
		t.Fatalf("failed to read failed file: %v", err)
	}
	want, err := os.ReadFile(filepath.FromSlash("./testdata/good_golden.png"))
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}

	v1, _, err := image.Decode(bytes.NewReader(got))
	if err != nil {
		t.Fatalf("unexpected error decoding failed file: %v", err)
	}
	v2, _, err := image.Decode(bytes.NewReader(want))
	if err != nil {
		t.Fatalf("unexpected error decoding golden file: %v", err)
	}

	dst := image.NewRGBA64(v1.Bounds().Union(v2.Bounds()))
	rect := Diff(dst, v1, v2)
	if rect != dst.Bounds() {
		t.Errorf("unexpected bound for diff: got:%+v want:%+v", rect, dst.Bounds())
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, dst)
	if err != nil {
		t.Fatalf("failed to encode difference png: %v", err)
	}
	gotDiff := base64.StdEncoding.EncodeToString(buf.Bytes())
	if gotDiff != wantDiffEncoded {
		t.Errorf("unexpected encoded diff value:\ngot:%s\nwant:%s", gotDiff, wantDiffEncoded)
	}
}

func TestEqual(t *testing.T) {
	got, err := os.ReadFile("testdata/approx_got_golden.png")
	if err != nil {
		t.Fatal(err)
	}

	ok, err := Equal("png", got, got)
	if err != nil {
		t.Fatalf("could not compare images: %+v", err)
	}
	if !ok {
		t.Fatalf("same image does not compare equal")
	}
}

func TestEqualApproxPNG(t *testing.T) {
	got, err := os.ReadFile("testdata/approx_got_golden.png")
	if err != nil {
		t.Fatal(err)
	}

	want, err := os.ReadFile("testdata/approx_want_golden.png")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		delta float64
		ok    bool
	}{
		{0, false},
		{0.01, false},
		{0.02, false},
		{0.05, true},
		{0.1, true},
		{1, true},
	} {
		t.Run(fmt.Sprintf("delta=%g", tc.delta), func(t *testing.T) {
			ok, err := EqualApprox("png", got, want, tc.delta)
			if err != nil {
				t.Fatalf("could not compare images: %+v", err)
			}
			if ok != tc.ok {
				t.Fatalf("got=%v, want=%v", ok, tc.ok)
			}
		})
	}
}

func TestEqualApprox(t *testing.T) {
	read := func(name string) []byte {
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("could not read file %q: %+v", name, err)
		}
		return raw
	}

	asPNG_RGBA64 := func(raw []byte) []byte {
		src, _, err := image.Decode(bytes.NewReader(raw))
		if err != nil {
			t.Fatalf("could not decode image: %+v", err)
		}
		var (
			bnds = src.Bounds()
			dst  = image.NewRGBA64(bnds)
			out  = new(bytes.Buffer)
		)
		draw.Draw(dst, bnds, src, image.Point{}, draw.Src)
		err = png.Encode(out, dst)
		if err != nil {
			t.Fatalf("could not encode image: %+v", err)
		}
		return out.Bytes()
	}

	for _, tc := range []struct {
		name  string
		img1  []byte
		img2  []byte
		delta float64
		want  bool
	}{
		{
			name: "svg-ok",
			img1: []byte("<svg></svg>"),
			img2: []byte("<svg></svg>"),
			want: true,
		},
		{
			name: "svg-diff",
			img1: []byte("<svg></svg>"),
			img2: []byte("<svg>1</svg>"),
			want: false,
		},
		{
			name: "eps-ok",
			img1: []byte("line1\nline2\nCreationDate:now\n"),
			img2: []byte("line1\nline2\nCreationDate:later\n"),
			want: true,
		},
		{
			name: "eps-diff-1",
			img1: []byte("line1\nline2\nCreationDate:now\n"),
			img2: []byte("line1\nline2\nCreationDate:later"),
			want: false,
		},
		{
			name: "eps-diff-2",
			img1: []byte("line1\nline2\nCreationDate:now\n"),
			img2: []byte("line1\nline3\nCreationDate:later\n"),
			want: false,
		},
		{
			name: "pdf-ok",
			img1: read("../vg/vgpdf/testdata/arc_golden.pdf"),
			img2: read("../vg/vgpdf/testdata/arc_golden.pdf"),
			want: true,
		},
		{
			name: "pdf-diff",
			img1: read("../vg/vgpdf/testdata/arc_golden.pdf"),
			img2: read("../vg/vgpdf/testdata/issue540_golden.pdf"),
			want: false,
		},
		{
			name: "pdf-diff-2",
			img1: read("../vg/vgpdf/testdata/arc_golden.pdf"),
			img2: read("../vg/vgpdf/testdata/multipage_golden.pdf"),
			want: false,
		},
		{
			name:  "png-ok",
			img1:  read("testdata/approx_got_golden.png"),
			img2:  read("testdata/approx_want_golden.png"),
			delta: 0.1,
			want:  true,
		},
		{
			name:  "png-ok-rgba64",
			img1:  read("testdata/approx_got_golden.png"),
			img2:  asPNG_RGBA64(read("testdata/approx_want_golden.png")),
			delta: 0.1,
			want:  true,
		},
		{
			name:  "png-ok-rgba64-2",
			img1:  read("testdata/approx_got_golden.png"),
			img2:  asPNG_RGBA64(read("testdata/approx_got_golden.png")),
			delta: -10, // clips to 0.0
			want:  true,
		},
		{
			name:  "png-ok-rgba64-3",
			img1:  asPNG_RGBA64(read("testdata/approx_got_golden.png")),
			img2:  read("testdata/approx_got_golden.png"),
			delta: 0,
			want:  true,
		},
		{
			name:  "png-diff-1",
			img1:  read("testdata/approx_got_golden.png"),
			img2:  read("testdata/approx_want_golden.png"),
			delta: 0,
			want:  false,
		},
		{
			name:  "png-diff-2",
			img1:  read("testdata/approx_got_golden.png"),
			img2:  read("testdata/approx_want_golden.png"),
			delta: 0.01,
			want:  false,
		},
		{
			name:  "png-diff-3",
			img1:  read("testdata/approx_got_golden.png"),
			img2:  read("testdata/good_golden.png"),
			delta: 10, // clips to 1.0
			want:  false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			typ := tc.name[:strings.Index(tc.name, "-")]
			got, err := EqualApprox(typ, tc.img1, tc.img2, tc.delta)
			if err != nil {
				t.Fatalf("could not run equal-approx: %+v", err)
			}

			if got != tc.want {
				t.Fatalf("invalid equal-approx: got=%v, want=%v", got, tc.want)
			}
		})
	}
}
