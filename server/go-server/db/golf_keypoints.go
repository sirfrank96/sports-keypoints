package db

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

type GolfKeypoints struct {
	Id                    primitive.ObjectID        `bson:"_id,omitempty"`
	UserId                string                    `bson:"user_id,omitempty"`
	InputImageId          string                    `bson:"input_image_id,omitempty"`
	OutputImg             []byte                    `bson:"output_img,omitempty"`
	OutputKeypoints       skp.Body25PoseKeypoints   `bson:"output_keypoints,omitempty"`
	DtlGolfSetupPoints    skp.DTLGolfSetupPoints    `bson:"dtl_golf_setup_points,omitempty"`
	FaceonGolfSetupPoints skp.FaceOnGolfSetupPoints `bson:"faceon_golf_setup_points,omitempty"`
}

func ConvertGolfKeypointsToCVGolfKeypoints(golfKeypoints *GolfKeypoints) *skp.GolfKeypoints {
	return &skp.GolfKeypoints{
		DtlGolfSetupPoints:    &golfKeypoints.DtlGolfSetupPoints,
		FaceonGolfSetupPoints: &golfKeypoints.FaceonGolfSetupPoints,
		BodyKeypoints:         &golfKeypoints.OutputKeypoints,
	}
}

func UpdateOutputKeypointsFields(oldKeypoints *skp.Body25PoseKeypoints, newKeypoints *skp.Body25PoseKeypoints) *skp.Body25PoseKeypoints {
	oldReflectVal := reflect.ValueOf(oldKeypoints).Elem()
	newReflectVal := reflect.ValueOf(newKeypoints).Elem()
	numFields := newReflectVal.NumField()
	for i := 0; i < numFields; i++ {
		newField := newReflectVal.Field(i)
		if !newField.IsZero() {
			oldField := oldReflectVal.Field(i)
			oldField.Set(newField)
		}
	}
	return oldKeypoints
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
	filter := bson.M{"input_image_id": inputImgId}
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
	fmt.Printf("Updated golf keypoints result: id: %s, userid: %s, inputimgid: %s\n", updatedGolfKeypoints.Id, updatedGolfKeypoints.UserId, updatedGolfKeypoints.InputImageId)
	return &updatedGolfKeypoints, nil
}

func (d *DbManager) DeleteGolfKeypointsForInputImage(ctx context.Context, inputImgId string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	warning := d.deleteGolfKeypointsForInputImageHelper(ctx, inputImgId)
	if warning != nil {
		return warning
	}
	return nil
}

func (d *DbManager) deleteGolfKeypointsForInputImageHelper(ctx context.Context, inputImgId string) util.Warning {
	fmt.Printf("Deleting golfkeypoints for inputimgid: %s...\n", inputImgId)
	filter := bson.M{"input_image_id": inputImgId}
	res, err := d.golfKeypointCollection.DeleteOne(ctx, filter)
	if err != nil {
		return util.WarningImpl{
			Severity: util.SEVERE,
			Message:  fmt.Sprintf("could not delete imageinfo %w", err),
		}
	}
	if res.DeletedCount == 0 {
		return util.WarningImpl{
			Severity: util.MINOR,
			Message:  fmt.Sprintf("did not delete any golfkeypoints, inputimgid %s may not exist", inputImgId),
		}
	}
	fmt.Printf("Delete golfkeypoints result: inputimgid: %s\n", inputImgId)
	return nil
}
