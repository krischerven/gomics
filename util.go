// Copyright (c) 2013-2018 Utkan Güngördü <utkan@freeconsole.org>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/krischerven/gomics/archive"
	"runtime"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func wrap(val, low, mod int) int {
	val %= mod
	if val < low {
		val = mod + val
	}
	return val
}

func clamp(val, low, high float64) float64 {
	if val < low {
		val = low
	} else if val > high {
		val = high
	}
	return val
}

func fit(sw, sh, fw, fh int) (int, int) {
	r := float64(sw) / float64(sh)
	var nw, nh float64
	if float64(fw) >= float64(fh)*r {
		nw, nh = float64(fh)*r, float64(fh)
	} else {
		nw, nh = float64(fw), float64(fw)/r
	}
	return int(nw), int(nh)
}

func tryPanic(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func gc() {
	heap := func(i uint8) {
		var x runtime.MemStats
		runtime.ReadMemStats(&x)
		fmt.Println(fmt.Sprintf("%d: %d", i, x.HeapAlloc))
	}
	linebreak := func() {
		fmt.Println(strings.Repeat("-", 10))
	}
	_ = heap
	_ = linebreak
	runtime.GC()
	runtime.GC()
}

func mustLoadPixbuf(data []byte) *gdk.Pixbuf {
	pixbuf, err := archive.LoadPixbuf(bytes.NewBuffer(data), true)
	tryPanic(err)
	return pixbuf
}
