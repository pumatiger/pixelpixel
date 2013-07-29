package main

import (
	"github.com/voxelbrain/pixelpixel/pixelutils"
)

func main() {
	c := pixelutils.PixelPusher()
	pixel := pixelutils.NewPixel()

	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			c := pixelutils.HSLA{
				H: float64(x) / 256,
				S: 1.0,
				L: 0.5 * (1 - float64(y)/256),
				A: 1.0,
			}
			pixel.Set(x, y, c)
		}
	}

	c <- pixel

	// Block indefinitely
	select {}
}
