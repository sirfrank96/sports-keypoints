package draw

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"testing"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

var bodyDatapoints = &skp.Body25PoseDatapoints{
	RShoulder: &skp.Datapoint{
		X:          100,
		Y:          100,
		Confidence: 1.0,
	},
	Neck: &skp.Datapoint{
		X:          150,
		Y:          150,
		Confidence: 1.0,
	},
	LShoulder: &skp.Datapoint{
		X:          200,
		Y:          200,
		Confidence: 1.0,
	},
}

func TestDrawBodyDatapoints(t *testing.T) {
	ctx := context.Background()
	// get current executables path
	executable, err := os.Executable()
	if err != nil {
		t.Errorf("Could not get current executable info: %v", err)
	}
	currentFileDirectory := path.Dir(executable)
	// grab file from path
	path := filepath.Join(currentFileDirectory, "..", "..", "..", "..", "static", "dtl-spineangle-normal.jpg")
	file, closeReadFile, err := util.GetFileFromPath(path)
	if err != nil {
		t.Errorf("Failed to getFileFromPath: %v", err)
	}
	defer closeReadFile()
	bytesEncodedAsJpg, err := util.DecodeAndEncodeFileAsJpg(file)
	if err != nil {
		t.Errorf("Failed to decodeAndEncodeFileAsJpg for image: %v", err)
	}
	// test DrawBodyDatapoints
	outputImg, err := DrawBodySkeleton(ctx, bytesEncodedAsJpg, bodyDatapoints)
	if err != nil {
		t.Errorf("Could not draw body skeleton on image: %v", err)
	}
	// write back to another file
	jpegBytes, err := util.DecodeAndEncodeBytesAsJpg(outputImg)
	if err != nil {
		t.Errorf("Failed to decodeAndEncodeBytesAsJpg for return image: %v", err)
	}
	closeWriteFile, err := util.WriteBytesToJpgFile(jpegBytes, filepath.Join(currentFileDirectory, "body-datapoints.jpg"))
	if err != nil {
		t.Errorf("Failed to writeBytesToJpgFile: %v", err)
	}
	defer closeWriteFile()
}
