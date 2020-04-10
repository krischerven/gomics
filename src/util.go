package main

import (
	"bytes"
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/krischerven/gomics/archive"
	"runtime"
	"strings"
)

var (
	firstAutomaticGC = true
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
	lastHeapSize := uint64(0)
	heap := func(i uint8) {
		var x runtime.MemStats
		runtime.ReadMemStats(&x)
		lastHeapSize = x.HeapAlloc
		fmt.Println(fmt.Sprintf("%d: %d", i, x.HeapAlloc))
	}
	linebreak := func() {
		fmt.Println(strings.Repeat("-", 10))
	}
	// on the first gc() call, it actually saves
	// about a kilobyte of memory to omit the second
	// runtime.GC() call (+2072 -> +1320 bytes)
	measureHeapChange := func(fn func()) {
		a := lastHeapSize
		fn()
		fmt.Println(lastHeapSize - a)
	}
	_ = heap
	_ = linebreak
	_ = measureHeapChange
	runtime.GC()
	if !firstAutomaticGC {
		runtime.GC()
	} else {
		firstAutomaticGC = false
	}
	/* DEBUGGING CODE
	heap(1)
	runtime.GC()
	heap(2)
	measureHeapChange(func() {
		runtime.GC()
		heap(3)
	})
	go func() {
		time.Sleep(time.Second*5)
		heap(0)
	} ()
	*/
}

func mustLoadPixbuf(data []byte) *gdk.Pixbuf {
	pixbuf, err := archive.LoadPixbuf(bytes.NewBuffer(data), true)
	tryPanic(err)
	return pixbuf
}
