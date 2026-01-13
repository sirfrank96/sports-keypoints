package draw

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
)

type Draw2dImage struct {
	gc  *draw2dimg.GraphicContext
	img image.Image
}

type GoFontCache map[string]*truetype.Font

func (fc GoFontCache) Store(fd draw2d.FontData, font *truetype.Font) {
	fc[fd.Name] = font
}

func (fc GoFontCache) Load(fd draw2d.FontData) (*truetype.Font, error) {
	/*font, stored := fc[fd.Name]
	if !stored {
		return nil, fmt.Errorf("font %s is not stored in font cache.", fd.Name)
	}
	return font, nil*/
	return fc["goregular"], nil
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
	// initialize fontcache for text
	fontCache := GoFontCache{}
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("error parsing font: %v", err)
	}
	fontCache.Store(draw2d.FontData{Name: "goregular"}, font)
	gc.FontCache = fontCache
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
	if datapointIsEmpty(point1) || datapointIsEmpty(point2) {
		return
	}
	d2dimg.gc.SetStrokeColor(c)
	d2dimg.gc.SetLineWidth(5.0)
	d2dimg.gc.MoveTo(point1.X, point1.Y)
	d2dimg.gc.LineTo(point2.X, point2.Y)
	d2dimg.gc.Stroke()
}

func (d2dimg *Draw2dImage) drawPoint(point *skp.Datapoint, c color.Color) {
	if datapointIsEmpty(point) {
		return
	}
	radius := 10.0
	d2dimg.gc.SetFillColor(c)
	d2dimg.gc.SetStrokeColor(c)
	d2dimg.gc.SetLineWidth(5.0)
	d2dimg.gc.BeginPath()
	draw2dkit.Circle(d2dimg.gc, point.X, point.Y, radius)
	d2dimg.gc.FillStroke()
}

func (d2dimg *Draw2dImage) addText(point *skp.Datapoint, text string) {
	d2dimg.gc.SetFontSize(20)
	d2dimg.gc.SetFillColor(getColorBlue())
	d2dimg.gc.FillStringAt(text, point.X, point.Y)
}

func getColorRed() color.RGBA {
	return color.RGBA{R: 255, G: 0, B: 0, A: 255}
}

func getColorBlue() color.RGBA {
	return color.RGBA{R: 0, G: 0, B: 255, A: 255}
}

func datapointIsEmpty(point *skp.Datapoint) bool {
	if point == nil {
		return true
	}
	if point.X == 0 && point.Y == 0 {
		return true
	}
	return false
}
