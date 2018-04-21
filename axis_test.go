// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/gonum/plot/vg"
)

var axisSmallTickTests = []struct {
	min, max   float64
	wantValues []float64
	wantLabels []string
}{
	{
		min:        -1.9846500878911073,
		max:        0.4370974820125605,
		wantValues: []float64{-1.75, -0.75, 0.25},
		wantLabels: []string{"-1.75", "-0.75", "0.25"},
	},
	{
		min:        -1.985e15,
		max:        0.4371e15,
		wantValues: []float64{-1.75e15, -7.5e14, 2.5e14},
		wantLabels: []string{"-1.75e+15", "-7.5e+14", "2.5e+14"},
	},
	{
		min:        -1.985e-15,
		max:        0.4371e-15,
		wantValues: []float64{-1.985e-15, -7.739500000000001e-16, 4.3709999999999994e-16},
		wantLabels: []string{"-1.985e-15", "-7.7395e-16", "4.371e-16"},
	},
	{
		min:        math.MaxFloat64 / 4,
		max:        math.MaxFloat64 / 3,
		wantValues: []float64{4.4942328371557893e+307, 5.243271643348421e+307, 5.992310449541053e+307},
		wantLabels: []string{"4e+307", "5e+307", "6e+307"},
	},
	{
		min:        0.00010,
		max:        0.00015,
		wantValues: []float64{0.0001, 0.00012, 0.00014000000000000001},
		wantLabels: []string{"0.0001", "0.00012", "0.00014"},
	},
	{
		min:        555.6545,
		max:        21800.9875,
		wantValues: []float64{4000, 12000, 20000},
		wantLabels: []string{"4000", "12000", "20000"},
	},
	{
		min:        555.6545,
		max:        27800.9875,
		wantValues: []float64{5000, 15000, 25000},
		wantLabels: []string{"5000", "15000", "25000"},
	},
	{
		min:        55.6545,
		max:        1555.9875,
		wantValues: []float64{300, 900, 1500},
		wantLabels: []string{"300", "900", "1500"},
	},
	{
		min:        3.096916 - 0.125,
		max:        3.096916 + 0.125,
		wantValues: []float64{3, 3.1, 3.2},
		wantLabels: []string{"3.0", "3.1", "3.2"},
	},
}

func TestAxisSmallTick(t *testing.T) {
	d := DefaultTicks{}
	for i, test := range axisSmallTickTests {
		ticks := d.Ticks(test.min, test.max)
		gotLabels := labelsOf(ticks)
		gotValues := valuesOf(ticks)
		if !reflect.DeepEqual(gotValues, test.wantValues) {
			t.Errorf("tick values mismatch %d:\ngot: %v\nwant:%v", i, gotValues, test.wantValues)
		}
		if !reflect.DeepEqual(gotLabels, test.wantLabels) {
			t.Errorf("tick labels mismatch %d:\ngot: %q\nwant:%q", i, gotLabels, test.wantLabels)
		}
	}
}

func valuesOf(ticks []Tick) []float64 {
	var values []float64
	for _, t := range ticks {
		if t.Label != "" {
			values = append(values, t.Value)
		}
	}
	return values
}

func labelsOf(ticks []Tick) []string {
	var labels []string
	for _, t := range ticks {
		if t.Label != "" {
			labels = append(labels, t.Label)
		}
	}
	return labels
}

func allLabelsOf(ticks []Tick) []string {
	var labels []string
	for _, t := range ticks {
		labels = append(labels, t.Label)
	}
	return labels
}

func TestTimeTicks(t *testing.T) {
	d := TimeTicks{Width: 4 * vg.Inch}
	for _, test := range []struct {
		min, max string
		want     []string
	}{
		{
			min:  "2016-01-01 12:56:30",
			max:  "2016-01-01 12:56:31",
			want: []string{"12:56:30.200", ".400", ".600", ".800", "12:56:31.000"},
		},
		{
			min:  "2016-01-01 12:56:01",
			max:  "2016-01-01 12:56:59",
			want: []string{"12:56:05", ":10", ":15", ":20", ":25", ":30", ":35", ":40", ":45", ":50", ":55"},
		},
		{
			min:  "2016-01-01 12:56:30",
			max:  "2016-01-01 12:57:29",
			want: []string{"12:56:35", ":40", ":45", ":50", ":55", "12:57:00", ":05", ":10", ":15", ":20", ":25"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-01-01 12:07:00",
			want: []string{"12:02:00", "12:03:00", "12:04:00", "12:05:00", "12:06:00", "12:07:00"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-01-01 12:17:00",
			want: []string{"Jan 01, 12:02pm", "12:04pm", "12:06pm", "12:08pm", "12:10pm", "12:12pm", "12:14pm", "12:16pm"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-01-01 12:28:00",
			want: []string{"Jan 01, 12:05pm", "12:10pm", "12:15pm", "12:20pm", "12:25pm"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-01-01 12:35:00",
			want: []string{"Jan 01, 12:05pm", "12:10pm", "12:15pm", "12:20pm", "12:25pm", "12:30pm", "12:35pm"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-01-01 12:40:00",
			want: []string{"Jan 01, 12:05pm", "12:10pm", "12:15pm", "12:20pm", "12:25pm", "12:30pm", "12:35pm", "12:40pm"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-01-01 12:45:00",
			want: []string{"Jan 01, 12:05pm", "12:10pm", "12:15pm", "12:20pm", "12:25pm", "12:30pm", "12:35pm", "12:40pm", "12:45pm"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-01-01 13:05:00",
			want: []string{"Jan 01, 12:15pm", "12:30pm", "12:45pm", "1:00pm"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-01-01 13:05:00",
			want: []string{"Jan 01, 12:15pm", "12:30pm", "12:45pm", "1:00pm"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-01-01 16:05:00",
			want: []string{"Jan 1, 1pm", "2pm", "3pm", "4pm"},
		},
		{
			min:  "2016-01-01 20:01:05",
			max:  "2016-01-02 07:59:00",
			want: []string{"Jan 1, 9pm", "10pm", "11pm", "Jan 2, 12am", "1am", "2am", "3am", "4am", "5am", "6am", "7am"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-01-02 13:59:00",
			want: []string{"Jan 1, 6pm", "Jan 2, 12am", "6am", "12pm"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-01-04 13:59:00",
			want: []string{"Jan 2", "12pm", "Jan 3", "12pm", "Jan 4", "12pm"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-01-06 13:59:00",
			want: []string{"Jan 2", "3", "4", "5", "6"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-01-09 13:59:00",
			want: []string{"Jan 2", "4", "6", "8"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-01-25 13:59:00",
			want: []string{"Jan 2", "4", "6", "8", "10", "12", "14", "16", "18", "20", "22", "24"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-02-06 13:59:00",
			want: []string{"Jan 4", "11", "18", "25", "Feb 1"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-02-28 13:59:00",
			want: []string{"Jan 4", "11", "18", "25", "Feb 1", "8", "15", "22"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-04-28 13:59:00",
			want: []string{"Jan 4", "11", "18", "25", "Feb 1", "8", "15", "22", "29", "Mar 7", "14", "21", "28", "Apr 4", "11", "18", "25"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-09-28 13:59:00",
			want: []string{"Feb 2016", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2016-12-28 13:59:00",
			want: []string{"Feb 2016", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2017-02-28 13:59:00",
			want: []string{"Feb 2016", "May", "Aug", "Nov", "Feb 2017"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2017-08-28 13:59:00",
			want: []string{"Feb 2016", "May", "Aug", "Nov", "Feb 2017", "May", "Aug"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2018-08-28 13:59:00",
			want: []string{"Feb 2016", "Aug", "Feb 2017", "Aug", "Feb 2018", "Aug"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2020-08-28 13:59:00",
			want: []string{"2016", "2017", "2018", "2019", "2020"},
		},
		{
			min:  "2016-01-01 12:01:05",
			max:  "2048-08-28 13:59:00",
			want: []string{"2017", "2022", "2027", "2032", "2037", "2042", "2047"},
		},
	} {
		//fmt.Println("For dates", test.min, test.max)
		ticks := d.Ticks(dateToFloat64(test.min), dateToFloat64(test.max))
		got := allLabelsOf(ticks)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("tick labels mismatch:\ndate1: %s\ndate2: %s\ngot: %#v\nwant:%q", test.min, test.max, got, test.want)
		}
	}
}

func dateToFloat64(date string) float64 {
	t, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		panic(err)
	}

	return float64(t.UTC().UnixNano()) / float64(time.Second)
}
