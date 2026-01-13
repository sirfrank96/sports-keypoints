package draw

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
)

type Draw2dImage struct {
	gc  *draw2dimg.GraphicContext
	img image.Image
}

func createDraw2dImage(imgBytes []byte) (*Draw2dImage, error) {
	// convert byte slice to Image
	img, err := jpeg.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, fmt.Errorf("error decoding jpeg to image: %v", err)
	}
	// create destination image and graphic context
	dest := image.NewRGBA(img.Bounds())
	gc := draw2dimg.NewGraphicContext(dest)
	gc.DrawImage(img)
	return &Draw2dImage{gc: gc, img: dest}, nil
}

func (d2dimg *Draw2dImage) getOutputByteSlice() ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := jpeg.Encode(buffer, d2dimg.img, nil)
	if err != nil {
		return nil, fmt.Errorf("error encoding as jpeg: %v", err)
	}
	return buffer.Bytes(), nil
}

func (d2dimg *Draw2dImage) drawLine(point1 *skp.Datapoint, point2 *skp.Datapoint, c color.Color) {
	d2dimg.gc.MoveTo(point1.X, point1.Y)
	d2dimg.gc.LineTo(point2.X, point2.Y)
	d2dimg.gc.SetStrokeColor(c)
	d2dimg.gc.SetLineWidth(5.0)
	d2dimg.gc.Stroke()
}

func (d2dimg *Draw2dImage) drawPoint(point *skp.Datapoint, c color.Color) {
	radius := 5.0
	d2dimg.gc.SetFillColor(c)
	d2dimg.gc.SetStrokeColor(c)
	d2dimg.gc.SetLineWidth(5.0)
	d2dimg.gc.BeginPath()
	draw2dkit.Circle(d2dimg.gc, point.X, point.Y, radius)
	d2dimg.gc.FillStroke()
}
