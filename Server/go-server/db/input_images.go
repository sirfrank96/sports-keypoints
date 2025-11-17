package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	cv "github.com/sirfrank96/go-server/computer-vision-sports-proto"
	"github.com/sirfrank96/go-server/util"
)

type InputImage struct {
	Id                           primitive.ObjectID   `bson:"_id,omitempty"`
	UserId                       string               `bson:"user_id,omitempty"`
	ImageType                    cv.ImageType         `bson:"image_type,omitempty"`
	InputImg                     []byte               `bson:"input_img,omitempty"`
	CalibrationImgAxes           []byte               `bson:"calibration_img_axes,omitempty"`
	CalibrationImgVanishingPoint []byte               `bson:"calibration_img_vanishing_point,omitempty"`
	CalibrationInfo              util.CalibrationInfo `bson:"calibration_info,omitempty"`
}

func (d *DbManager) CreateInputImage(ctx context.Context, inputImg *InputImage) (*InputImage, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Creating input image...\n")
	res, err := d.inputImageCollection.InsertOne(ctx, inputImg)
	if err != nil {
		return nil, fmt.Errorf("could not create input image: %w", err)
	}
	objectId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("could create object id")
	}
	inputImg.Id = objectId
	fmt.Printf("Create input image result: imgId: %s, userId: %s, imageType: %s\n", inputImg.Id, inputImg.UserId, inputImg.ImageType)
	return inputImg, nil
}

func (d *DbManager) ReadInputImage(ctx context.Context, inputImgId string) (*InputImage, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Reading input image id: %s...\n", inputImgId)
	objectId, err := primitive.ObjectIDFromHex(inputImgId)
	if err != nil {
		return nil, fmt.Errorf("could not convert id to object id %w", err)
	}
	filter := bson.M{"_id": objectId}
	var inputImg InputImage
	if err := d.inputImageCollection.FindOne(ctx, filter).Decode(&inputImg); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, fmt.Errorf("no input images with id: %s", inputImgId)
		}
		return nil, fmt.Errorf("could not read input image: %w", err)
	}
	//inputImg.Id = objectId.Hex()
	fmt.Printf("Read input image result: imgId: %s, userId: %s, imageType: %s\n", inputImg.Id, inputImg.UserId, inputImg.ImageType)
	return &inputImg, nil
}

func (d *DbManager) UpdateInputImage(ctx context.Context, inputImgId string, newInputImage *InputImage) (*InputImage, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Updating input image imgid: %s\n", inputImgId)
	objectId, err := primitive.ObjectIDFromHex(inputImgId)
	if err != nil {
		return nil, fmt.Errorf("could not convert id to object id %w", err)
	}
	filter := bson.M{"_id": objectId}
	update := bson.M{
		"$set": bson.M{
			"user_id":                         newInputImage.UserId,
			"image_type":                      newInputImage.ImageType,
			"input_img":                       newInputImage.InputImg,
			"calibration_img_axes":            newInputImage.CalibrationImgAxes,
			"calibration_img_vanishing_point": newInputImage.CalibrationImgVanishingPoint,
			"calibration_info":                newInputImage.CalibrationInfo,
		},
	}
	var updatedInputImage InputImage
	if err := d.inputImageCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedInputImage); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, fmt.Errorf("no input images with imgid: %s", inputImgId)
		}
		return nil, fmt.Errorf("could not update input image: %w", err)
	}
	//updatedInputImage.Id = objectId.Hex()
	fmt.Printf("Update input image result: %+v...\n", updatedInputImage)
	return &updatedInputImage, nil
}

func (d *DbManager) DeleteInputImage(ctx context.Context, inputImgId string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Deleting input image id: %s...\n", inputImgId)
	objectId, err := primitive.ObjectIDFromHex(inputImgId)
	if err != nil {
		return fmt.Errorf("could not convert id to object id %w", err)
	}
	filter := bson.M{"_id": objectId}
	res, err := d.inputImageCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("could not delete input image %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("did not delete any images, imdId %s may not exist", inputImgId)
	}
	fmt.Printf("Delete input image result: imgId: %s\n", inputImgId)
	return nil
}

// TODO: If delete image, delete all associated ImageInfos
// for a user, delete all input images associated with it
func (d *DbManager) DeleteInputImagesForUser() {

}
