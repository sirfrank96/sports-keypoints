package draw

import (
	"context"
	"fmt"
	"image/color"

	"google.golang.org/protobuf/reflect/protoreflect"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
)

func DrawBodySkeleton(ctx context.Context, inputImg []byte, bodyDatapoints *skp.Body25PoseDatapoints) ([]byte, error) {
	d2dimg, err := createDraw2dImage(inputImg)
	if err != nil {
		return nil, fmt.Errorf("could not create a draw2dimg graphic context: %v", err)
	}
	return drawBodySkeleton(ctx, d2dimg, bodyDatapoints)
}

func drawBodySkeleton(ctx context.Context, d2dimg *Draw2dImage, bodyDatapoints *skp.Body25PoseDatapoints) ([]byte, error) {
	// draw body data points
	prm := bodyDatapoints.ProtoReflect()
	prm.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		datapoint, ok := v.Message().Interface().(*skp.Datapoint)
		if ok {
			fmt.Printf("Drawing a point\n")
			d2dimg.drawPoint(datapoint, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
		return true
	})
	// draw line from lshoulder to rshoulder
	d2dimg.drawLine(bodyDatapoints.LShoulder, bodyDatapoints.RShoulder, color.RGBA{R: 0, G: 0, B: 255, A: 255})
	return d2dimg.getOutputByteSlice()
}
