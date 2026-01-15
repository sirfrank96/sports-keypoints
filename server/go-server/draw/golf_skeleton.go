package draw

import (
	"context"
	"fmt"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
)

func DrawGolfSkeleton(ctx context.Context, inputImg []byte, bodyDatapoints *skp.Body25PoseDatapoints, golfDatapoints *skp.GolfSpecificDatapoints) ([]byte, error) {
	// draw body skeleton
	d2dimg, err := createDraw2dImage(inputImg)
	if err != nil {
		return nil, fmt.Errorf("could not create a draw2dimg graphic context: %v", err)
	}
	_, err = drawBodySkeleton(ctx, d2dimg, bodyDatapoints)
	if err != nil {
		return nil, fmt.Errorf("could not draw body skeleton: %v", err)
	}
	// draw golf equipment stuff
	d2dimg.drawPoint(golfDatapoints.GolfBall, getColorRed())
	d2dimg.addText(golfDatapoints.GolfBall, "golf_ball")
	d2dimg.drawPoint(golfDatapoints.ClubButt, getColorRed())
	d2dimg.addText(golfDatapoints.ClubButt, "club_butt")
	d2dimg.drawPoint(golfDatapoints.ClubHead, getColorRed())
	d2dimg.addText(golfDatapoints.ClubHead, "club_head")
	d2dimg.drawLine(golfDatapoints.ClubButt, golfDatapoints.ClubHead, getColorRed())
	return d2dimg.getOutputByteSlice()
}
