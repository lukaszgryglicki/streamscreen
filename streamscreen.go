package main

import (
	"fmt"
	"image/png"
	"os"
	"strconv"
	"time"

	"github.com/kbinani/screenshot"
)

func streamScreen() error {
	nDisplays := screenshot.NumActiveDisplays()
	display := 0
	displayStr := os.Getenv("D")
	if displayStr != "" {
		disp, err := strconv.Atoi(displayStr)
		if err != nil {
			return err
		}
		if disp > 0 {
			display = disp
		}
	}
	if nDisplays < display+1 {
		return fmt.Errorf("need at least %d active display(s), found %d", display+1, nDisplays)
	}
	bounds := screenshot.GetDisplayBounds(display)
	fmt.Printf("Display #%d bounds: %+v\n", display, bounds)
	x := bounds.Min.X
	y := bounds.Min.Y
	w := bounds.Max.X
	h := bounds.Max.Y
	str := os.Getenv("X")
	if str != "" {
		v, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		if v >= 0 && v <= w {
			x = v
		}
	}
	str = os.Getenv("W")
	if str != "" {
		v, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		if v >= 0 && v <= w-x {
			w = v
		}
	}
	if w > bounds.Max.X-x {
		w = bounds.Max.X - x
	}
	str = os.Getenv("Y")
	if str != "" {
		v, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		if v >= 0 && v <= h {
			y = v
		}
	}
	str = os.Getenv("H")
	if str != "" {
		v, err := strconv.Atoi(str)
		if err != nil {
			return err
		}
		if v >= 0 && v <= h-y {
			h = v
		}
	}
	if h > bounds.Max.Y-y {
		h = bounds.Max.Y - y
	}
	bounds.Min.X = x
	bounds.Min.Y = y
	bounds.Max.X = w
	bounds.Max.Y = h
	fmt.Printf("Bounds: %+v\n", bounds)
	// PNG Quality
	pngqStr := os.Getenv("PQ")
	pngq := png.BestSpeed
	if pngqStr != "" {
		v, err := strconv.Atoi(pngqStr)
		if err != nil {
			return err
		}
		if v < 0 || v > 3 {
			return fmt.Errorf("PQ must be from 0-3 range")
		}
		pngq = png.CompressionLevel(-v)
	}
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return err
	}
	sshotFunc := func() (string, error) {
		fn := fmt.Sprintf("%d.png", time.Now().UnixNano())
		file, err := os.Create(fn)
		if err != nil {
			return "", err
		}
		enc := png.Encoder{CompressionLevel: pngq}
		ierr := enc.Encode(file, img)
		if ierr != nil {
			_ = file.Close()
			return "", ierr
		}
		err = file.Close()
		if err != nil {
			return "", err
		}
		return fn, nil
	}
	ss := os.Getenv("SS")
	if ss != "" {
		fn, err := sshotFunc()
		if err != nil {
			return err
		}
		fmt.Printf("SS saved to %s\n", fn)
	}
	sv := os.Getenv("SV")
	if sv != "" {
		// FIXME: handle CTRL+c
		sss := []string{}
		dtStart := time.Now()
		for i := 0; i < 100; i++ {
			fn, err := sshotFunc()
			if err != nil {
				return err
			}
			sss = append(sss, fn)
			// FIXME: handle FPS (or unlimited)
		}
		// FIXME: handle calculate actual FPS
		// FIXME: handle encode pngs as mp4 and remove them
		dtEnd := time.Now()
		us := dtEnd.Sub(dtStart).Nanoseconds()
		fmt.Printf("Took: %d us\n", us)
		fmt.Printf("SV saved to %+v\n", sss)
	}
	return nil
}

func main() {
	dtStart := time.Now()
	err := streamScreen()
	if err != nil {
		fmt.Printf("error: %+v\n", err)
	}
	dtEnd := time.Now()
	fmt.Printf("Took: %v\n", dtEnd.Sub(dtStart))
}
