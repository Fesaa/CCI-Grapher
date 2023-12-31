package cubecounterimage

import (
	"fmt"
	"image"
	"image/draw"

	"github.com/wcharczuk/go-chart/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

func (ccr *cubeCounterRequest) getFinalImage(imgs []image.Image) image.Image {
    if len(imgs) == 2 {
        return ccr.userImageMerge(imgs)
    }
    return ccr.imageMerge(imgs)
}

func (ccr *cubeCounterRequest) userImageMerge(imgs []image.Image) image.Image {
	titleString := fmt.Sprintf("From %v to %v (UTC)", ccr.startDate.Format(dataParse), ccr.endDate.Format(dataParse))

    tm, ha := imgs[0], imgs[1]
    r := image.Rectangle{image.Point{0, 0}, image.Point{tm.Bounds().Dx(), tm.Bounds().Dy() * 2}.Add(image.Point{0, 50})}
    titleRect := image.Rectangle{image.Point{0, 0}, image.Point{tm.Bounds().Dx(), 50}}
    rha := image.Rectangle{image.Point{0, tm.Bounds().Dy()}, tm.Bounds().Size().Add(ha.Bounds().Size()).Add(image.Point{0, 50})}

    rgba := image.NewRGBA(r)
	draw.Draw(rgba, titleRect, &image.Uniform{chart.ColorAlternateGray}, image.Point{0, 0}, draw.Src)
    draw.Draw(rgba, r, tm, image.Pt(0, 0), draw.Src)
    draw.Draw(rgba, rha, ha, image.Pt(0, 0), draw.Src)

    addLabel(rgba, 40 * tm.Bounds().Dx() / 100, tm.Bounds().Dy() + ha.Bounds().Dy() - 50, titleString)
    return rgba
}

func (ccr *cubeCounterRequest) imageMerge(imgs []image.Image) image.Image {
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

	addLabel(rgba, 85*tm.Bounds().Dx()/100, 20, titleString)
	return rgba
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := chart.ColorBlack
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: inconsolata.Bold8x16,
		Dot:  point,
	}
	d.DrawString(label)
}
