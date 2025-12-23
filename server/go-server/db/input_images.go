package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	skp "github.com/sirfrank96/go-server/sports-keypoints-proto"
	"github.com/sirfrank96/go-server/util"
)

type InputImage struct {
	Id                           primitive.ObjectID   `bson:"_id,omitempty"`
	UserId                       string               `bson:"user_id,omitempty"`
	ImageType                    skp.ImageType        `bson:"image_type,omitempty"`
	InputImg                     []byte               `bson:"input_img,omitempty"`
	Description                  string               `bson:"description,omitempty"`
	Timestamp                    time.Time            `bson:"timestamp,omitempty"`
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

func (d *DbManager) ReadInputImagesForUser(ctx context.Context, userId string) ([]*InputImage, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.readInputImagesForUserHelper(ctx, userId)
}

func (d *DbManager) readInputImagesForUserHelper(ctx context.Context, userId string) ([]*InputImage, error) {
	fmt.Printf("Reading input images for user...\n")
	filter := bson.M{"user_id": userId}
	cursor, err := d.inputImageCollection.Find(ctx, filter)
	if err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, fmt.Errorf("no input images for user")
		}
		return nil, fmt.Errorf("could not read input images for user: %w", err)
	}
	defer cursor.Close(ctx)
	var res []*InputImage
	for cursor.Next(ctx) {
		var inputImage InputImage
		if err := cursor.Decode(&inputImage); err != nil {
			return nil, fmt.Errorf("could not decode input image: %w", err)
		}
		res = append(res, &inputImage)
	}
	return res, nil
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
			"description":                     newInputImage.Description,
			"timestamp":                       newInputImage.Timestamp,
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
	fmt.Printf("Update input image result: imgId: %s, userId: %s, imageType: %s\n", updatedInputImage.Id, updatedInputImage.UserId, updatedInputImage.ImageType)
	return &updatedInputImage, nil
}

func (d *DbManager) DeleteInputImage(ctx context.Context, inputImgId string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.deleteInputImageHelper(ctx, inputImgId)
}

// Deletes golfkeypoints associated with input image, then delete input image
func (d *DbManager) deleteInputImageHelper(ctx context.Context, inputImgId string) error {
	fmt.Printf("Deleting input image id: %s...\n", inputImgId)
	// first delete keypoints associated with input image
	warning := d.deleteGolfKeypointsForInputImageHelper(ctx, inputImgId)
	if warning != nil {
		if warning.GetSeverity() == util.SEVERE {
			return fmt.Errorf("could not delete keypoints associated with input img %s: %w", inputImgId, warning.Error())
		} else {
			fmt.Printf("Minor warning: %s", warning.Error())
		}
	}
	// delete input image
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
