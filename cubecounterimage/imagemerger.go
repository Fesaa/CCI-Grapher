package cubecounterimage

import (
	"image"
	"image/draw"
)

func imageMerge(imgs []image.Image) image.Image {

	tm, ct, rd, ha := imgs[0], imgs[1], imgs[2], imgs[3]

	r := image.Rectangle{image.Point{0, 0}, image.Point{tm.Bounds().Dx() * 2, tm.Bounds().Dy() * 2}}
	rct := image.Rectangle{image.Point{tm.Bounds().Dx(), 0}, image.Point{tm.Bounds().Dx(), 0}.Add(ct.Bounds().Size())}
	rrd := image.Rectangle{image.Point{0, tm.Bounds().Dy()}, image.Point{0, tm.Bounds().Dy()}.Add(ct.Bounds().Size())}
	rha := image.Rectangle{image.Point{tm.Bounds().Dx(), tm.Bounds().Dy()}, image.Point{tm.Bounds().Dx(), tm.Bounds().Dy()}.Add(ct.Bounds().Size())}
	rgba := image.NewRGBA(r)
	draw.Draw(rgba, r, tm, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, rct, ct, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, rrd, rd, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, rha, ha, image.Point{0, 0}, draw.Src)

	return rgba
}
