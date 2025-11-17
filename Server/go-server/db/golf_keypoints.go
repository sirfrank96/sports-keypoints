package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
)

type GolfKeypoints struct {
	Id                    primitive.ObjectID       `bson:"_id,omitempty"`
	UserId                string                   `bson:"user_id,omitempty"`
	InputImageId          string                   `bson:"input_image_id,omitempty"`
	OutputImg             []byte                   `bson:"output_img,omitempty"`
	OutputKeypoints       cv.Body25PoseKeypoints   `bson:"output_keypoints,omitempty"`
	DtlGolfSetupPoints    cv.DTLGolfSetupPoints    `bson:"dtl_golf_setup_points,omitempty"`
	FaceonGolfSetupPoints cv.FaceOnGolfSetupPoints `bson:"faceon_golf_setup_points,omitempty"`
}

func ConvertGolfKeypointsToCVGolfKeypoints(golfKeypoints *GolfKeypoints) *cv.GolfKeypoints {
	return &cv.GolfKeypoints{
		OutputImg:             golfKeypoints.OutputImg,
		DtlGolfSetupPoints:    &golfKeypoints.DtlGolfSetupPoints,
		FaceonGolfSetupPoints: &golfKeypoints.FaceonGolfSetupPoints,
	}
}

func (d *DbManager) CreateGolfKeypoints(ctx context.Context, golfKeypoints *GolfKeypoints) (*GolfKeypoints, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Creating golf keypoint...\n")
	res, err := d.golfKeypointCollection.InsertOne(ctx, golfKeypoints)
	if err != nil {
		return nil, fmt.Errorf("could not create golf keypoint: %w", err)
	}
	objectId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("could create object id")
	}
	golfKeypoints.Id = objectId
	fmt.Printf("Create golf keypoints result: id: %s, userid: %s, inputimgid: %s\n", golfKeypoints.Id, golfKeypoints.UserId, golfKeypoints.InputImageId)
	return golfKeypoints, nil
}

func (d *DbManager) ReadGolfKeypointsForInputImage(ctx context.Context, inputImgId string) (*GolfKeypoints, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Reading golf keypoints for input img id: %s...\n", inputImgId)
	filter := bson.M{"input_image_id": inputImgId}
	var golfKeypoints GolfKeypoints
	if err := d.golfKeypointCollection.FindOne(ctx, filter).Decode(&golfKeypoints); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, fmt.Errorf("no golf keypoints with imgid: %s", inputImgId)
		}
		return nil, fmt.Errorf("could not read golf keypoints: %w", err)
	}
	fmt.Printf("Read golf keypoints result: id: %s, userid: %s, inputimgid: %s\n", golfKeypoints.Id, golfKeypoints.UserId, golfKeypoints.InputImageId)
	return &golfKeypoints, nil
}

func (d *DbManager) UpdateGolfKeypointsForInputImage(ctx context.Context, inputImgId string, newGolfKeypoints *GolfKeypoints) (*GolfKeypoints, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Updating golfkeypoints for inputimgid: %s\n", inputImgId)
	filter := bson.M{"input_img_id": inputImgId}
	update := bson.M{
		"$set": bson.M{
			"user_id":                  newGolfKeypoints.UserId,
			"input_img_id":             newGolfKeypoints.InputImageId,
			"output_img":               newGolfKeypoints.OutputImg,
			"output_keypoints":         newGolfKeypoints.OutputKeypoints,
			"dtl_golf_setup_points":    newGolfKeypoints.DtlGolfSetupPoints,
			"faceon_golf_setup_points": newGolfKeypoints.FaceonGolfSetupPoints,
		},
	}
	var updatedGolfKeypoints GolfKeypoints
	if err := d.golfKeypointCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedGolfKeypoints); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, fmt.Errorf("no golfkeypoints with inputimgid: %s", inputImgId)
		}
		return nil, fmt.Errorf("could not update golfkeypoints: %w", err)
	}
	fmt.Printf("Update golfkeypoints result: %+v...\n", updatedGolfKeypoints)
	return &updatedGolfKeypoints, nil
}

// TODO: If delete image, delete all associated ImageInfos
func (d *DbManager) DeleteGolfKeypointsForInputImage(ctx context.Context, inputImgId string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Deleting golfkeypoints for inputimgid: %s...\n", inputImgId)
	filter := bson.M{"input_img_id": inputImgId}
	res, err := d.golfKeypointCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("could not delete imageinfo %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("did not delete any golfkeypoints, inputimgid %s may not exist", inputImgId)
	}
	fmt.Printf("Delete golfkeypoints result: inputimgid: %s\n", inputImgId)
	return nil
}

// for a user, delete all keypoints associated with it
func (d *DbManager) DeleteGolfKeypointsForUser() {

}
