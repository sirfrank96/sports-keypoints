package draw

import (
	"context"
	"fmt"
	//"image/color"

	//"github.com/llgcode/draw2d/draw2dimg"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
)

func DrawGolfSkeleton(ctx context.Context, inputImg []byte, bodyDatapoints *skp.Body25PoseDatapoints, golfDatapoints *skp.GolfSpecificDatapoints) ([]byte, error) {
	// draw body skeleton
	d2dimg, err := createDraw2dImage(inputImg)
	if err != nil {
		return nil, fmt.Errorf("could not create a draw2dimg graphic context: %v", err)
	}
	outputImg, err := drawBodySkeleton(ctx, d2dimg, bodyDatapoints) // TODO: remove outputimg and continue passing around d2dimg
	if err != nil {
		return nil, fmt.Errorf("could not draw body skeleton: %v", err)
	}
	// TODO: draw golf equipment stuff
	return outputImg, nil
}
