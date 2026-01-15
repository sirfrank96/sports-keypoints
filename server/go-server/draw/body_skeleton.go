package draw

import (
	"context"
	"fmt"

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
			d2dimg.drawPoint(datapoint, getColorRed())
			d2dimg.addText(datapoint, string(fd.Name()))
		}
		return true
	})
	// draw head area
	d2dimg.drawLine(bodyDatapoints.LEar, bodyDatapoints.LEye, getColorRed())
	d2dimg.drawLine(bodyDatapoints.REar, bodyDatapoints.REye, getColorRed())
	d2dimg.drawLine(bodyDatapoints.LEye, bodyDatapoints.Nose, getColorRed())
	d2dimg.drawLine(bodyDatapoints.REye, bodyDatapoints.Nose, getColorRed())
	d2dimg.drawLine(bodyDatapoints.Nose, bodyDatapoints.Neck, getColorRed())
	// draw upper body
	d2dimg.drawLine(bodyDatapoints.LShoulder, bodyDatapoints.Neck, getColorRed())
	d2dimg.drawLine(bodyDatapoints.RShoulder, bodyDatapoints.Neck, getColorRed())
	d2dimg.drawLine(bodyDatapoints.LShoulder, bodyDatapoints.LElbow, getColorRed())
	d2dimg.drawLine(bodyDatapoints.LElbow, bodyDatapoints.LWrist, getColorRed())
	d2dimg.drawLine(bodyDatapoints.RShoulder, bodyDatapoints.RElbow, getColorRed())
	d2dimg.drawLine(bodyDatapoints.RElbow, bodyDatapoints.RWrist, getColorRed())
	d2dimg.drawLine(bodyDatapoints.Neck, bodyDatapoints.Midhip, getColorRed())
	// draw lower body
	d2dimg.drawLine(bodyDatapoints.LHip, bodyDatapoints.Midhip, getColorRed())
	d2dimg.drawLine(bodyDatapoints.RHip, bodyDatapoints.Midhip, getColorRed())
	d2dimg.drawLine(bodyDatapoints.LHip, bodyDatapoints.LKnee, getColorRed())
	d2dimg.drawLine(bodyDatapoints.LKnee, bodyDatapoints.LAnkle, getColorRed())
	d2dimg.drawLine(bodyDatapoints.RHip, bodyDatapoints.RKnee, getColorRed())
	d2dimg.drawLine(bodyDatapoints.RKnee, bodyDatapoints.RAnkle, getColorRed())
	// draw feet
	d2dimg.drawLine(bodyDatapoints.LAnkle, bodyDatapoints.LHeel, getColorRed())
	d2dimg.drawLine(bodyDatapoints.LAnkle, bodyDatapoints.LSmallToe, getColorRed())
	d2dimg.drawLine(bodyDatapoints.LAnkle, bodyDatapoints.LBigToe, getColorRed())
	d2dimg.drawLine(bodyDatapoints.RAnkle, bodyDatapoints.RHeel, getColorRed())
	d2dimg.drawLine(bodyDatapoints.RAnkle, bodyDatapoints.RSmallToe, getColorRed())
	d2dimg.drawLine(bodyDatapoints.RAnkle, bodyDatapoints.RBigToe, getColorRed())
	return d2dimg.getOutputByteSlice()
}
