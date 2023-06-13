package cubecounterimage

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/wcharczuk/go-chart/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func imageMerge(imgs []image.Image, ccr cubeCounterRequest) image.Image {

	titleString := fmt.Sprintf("From %v to %v (UTC)", ccr.startDate.Format(dataParse), ccr.endDate.Format(dataParse))

	tm, ct, rd, ha := imgs[0], imgs[1], imgs[2], imgs[3]

	r := image.Rectangle{image.Point{0, 0}, image.Point{tm.Bounds().Dx() * 2, tm.Bounds().Dy() * 2}.Add(image.Point{0, 100})}
	titleRect := image.Rectangle{image.Point{0, 0}, image.Point{tm.Bounds().Dx() * 2, 100}}
	rct := image.Rectangle{image.Point{tm.Bounds().Dx(), 0}, image.Point{tm.Bounds().Dx(), 0}.Add(ct.Bounds().Size()).Add(image.Point{0, 100})}
	rrd := image.Rectangle{image.Point{0, tm.Bounds().Dy()}, image.Point{0, tm.Bounds().Dy()}.Add(ct.Bounds().Size()).Add(image.Point{0, 100})}
	rha := image.Rectangle{image.Point{tm.Bounds().Dx(), tm.Bounds().Dy()}, image.Point{tm.Bounds().Dx(), tm.Bounds().Dy()}.Add(ct.Bounds().Size()).Add(image.Point{0, 100})}
	rgba := image.NewRGBA(r)
	draw.Draw(rgba, titleRect, &image.Uniform{chart.ColorAlternateGray}, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, r, tm, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, rct, ct, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, rrd, rd, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, rha, ha, image.Point{0, 0}, draw.Src)

	addLabel(rgba, tm.Bounds().Dx(), 50, titleString)
	return rgba
}

func addLabel(img *image.RGBA, x, y int, label string) {
    col := chart.ColorBlack
    point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

    d := &font.Drawer{
        Dst:  img,
        Src:  image.NewUniform(col),
        Face: basicfont.Face7x13,
        Dot:  point,
    }
    d.DrawString(label)
}
