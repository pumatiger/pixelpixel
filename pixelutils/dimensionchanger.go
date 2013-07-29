package pixelutils

import (
	"image"
	"image/color"
	"image/draw"
)

type dimensionChanger struct {
	draw.Image
	w, h               int
	pixel              image.Rectangle
	paddingX, paddingY int
}

func DimensionChanger(img draw.Image, w, h int) draw.Image {
	b := img.Bounds().Canon()
	return &dimensionChanger{
		Image:    img,
		w:        w,
		h:        h,
		paddingX: b.Dx() % w,
		paddingY: b.Dy() % h,
		pixel:    image.Rect(0, 0, b.Dx()/w, b.Dy()/h),
	}
}

func (d *dimensionChanger) Set(x, y int, c color.Color) {
	draw.Draw(d.Image, d.pixel.Add(image.Point{x * d.pixel.Dx(), y * d.pixel.Dy()}).Add(d.Image.Bounds().Canon().Min), &image.Uniform{c}, image.Point{0, 0}, draw.Over)
}

func (d *dimensionChanger) At(x, y int) color.Color {
	return d.Image.At(x*d.pixel.Dx()+d.Image.Bounds().Canon().Min.X, y*d.pixel.Dy()+d.Image.Bounds().Canon().Min.Y)
}

func (d *dimensionChanger) Bounds() image.Rectangle {
	return image.Rect(0, 0, d.w, d.h)
}
