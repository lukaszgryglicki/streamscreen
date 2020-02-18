package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/kbinani/screenshot"
)

func streamScreen() error {
	nDisplays := screenshot.NumActiveDisplays()
	display := 0
	displayStr := os.Getenv("DISPLAY")
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
