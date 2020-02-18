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
		nFrames := 0
		str := os.Getenv("F")
		if str != "" {
			v, err := strconv.Atoi(str)
			if err != nil {
				return err
			}
			if v >= 0 {
				nFrames = v
			}
		}
		fps := 0.0
		str = os.Getenv("FPS")
		if str != "" {
			v, err := strconv.ParseFloat(str, 64)
			if err != nil {
				return err
			}
			if v >= 0.0 {
				fps = v
			}
		}
		fpsn := 0.0
		if fps > 0.0 {
			fpsn = float64(1e9) / fps
		}
		f := 0
		var (
			dtf time.Time
			dtt time.Time
		)
		for {
			if fps > 0.0 {
				dtf = time.Now()
			}
			fn, err := sshotFunc()
			if err != nil {
				return err
			}
			sss = append(sss, fn)
			f++
			if nFrames > 0 && f >= nFrames {
				break
			}
			if fps > 0.0 {
				dtt = time.Now()
				us := float64(dtt.Sub(dtf).Nanoseconds())
				if us < fpsn {
					nano := fpsn - us
					// fmt.Printf("Should wait %.3fms\n", nano/float64(1e6))
					time.Sleep(time.Duration(nano) * time.Nanosecond)
				}
			}
		}
		n := len(sss)
		// FIXME: handle encode pngs as mp4 and remove them
		dtEnd := time.Now()
		secs := float64(dtEnd.Sub(dtStart).Nanoseconds()) / float64(1e9)
		afps := float64(n) / float64(secs)
		fmt.Printf("Took: %.3fs, actual FPS: %f\n", secs, afps)
		// fmt.Printf("SV saved to %+v\n", sss)
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
