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

func (d *DbManager) CreateInputImage(ctx context.Context, image *cv.Image) (*cv.Image, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Creating image...\n")
	res, err := d.imageCollection.InsertOne(ctx, image)
	if err != nil {
		return nil, fmt.Errorf("could not create image: %w", err)
	}
	objectId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("could create object id")
	}
	image.Id = objectId.Hex()
	fmt.Printf("Create image result: %s, %s\n", image.Id, image.UserId)
	return image, nil
}

func (d *DbManager) ReadInputImage(ctx context.Context, id string) (*cv.Image, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Reading image id: %s...\n", id)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("could not convert id to object id %w", err)
	}
	filter := bson.M{"_id": objectId}
	var image cv.Image
	if err := d.imageCollection.FindOne(ctx, filter).Decode(&image); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, fmt.Errorf("no images with id: %s", id)
		}
		return nil, fmt.Errorf("could not read image: %w", err)
	}
	image.Id = objectId.Hex()
	fmt.Printf("Read image result: imgId: %s, userId: %s\n", image.Id, image.UserId)
	return &image, nil
}

// TODO: If delete image, delete all associated ImageInfos
func (d *DbManager) DeleteInputImage(ctx context.Context, id string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Deleting image id: %s...\n", id)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("could not convert id to object id %w", err)
	}
	filter := bson.M{"_id": objectId}
	res, err := d.imageCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("could not delete image %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("did not delete any images, imageid %s may not exist", id)
	}
	fmt.Printf("Delete image result: imgId: %s\n", id)
	return nil
}

func (d *DbManager) DeleteUsersInputImages() {

}

func (d *DbManager) CreateImageInfo(ctx context.Context, imageInfo *cv.ImageInfo) (*cv.ImageInfo, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Creating image info...\n")
	_, err := d.imageInfoCollection.InsertOne(ctx, imageInfo)
	if err != nil {
		return nil, fmt.Errorf("could not create imageinfo: %w", err)
	}
	fmt.Printf("Create imageinfo result: userid: %s, inputimgid: %s\n", imageInfo.UserId, imageInfo.InputImgId)
	return imageInfo, nil
}

func (d *DbManager) ReadImageInfo(ctx context.Context, inputImgId string) (*cv.ImageInfo, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Reading image info img id: %s...\n", inputImgId)
	filter := bson.M{"inputimgid": inputImgId}
	var imageInfo cv.ImageInfo
	if err := d.imageInfoCollection.FindOne(ctx, filter).Decode(&imageInfo); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, fmt.Errorf("no imageinfos with imgid: %s", inputImgId)
		}
		return nil, fmt.Errorf("could not read imageinfo: %w", err)
	}
	fmt.Printf("Read image info result: userid: %s, inputimgid: %s\n", imageInfo.UserId, imageInfo.InputImgId)
	return &imageInfo, nil
}

func (d *DbManager) UpdateImageInfo(ctx context.Context, inputImgId string, imageInfo *cv.ImageInfo) (*cv.ImageInfo, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Updating imageinfo inputimgid: %s\n", inputImgId)

	filter := bson.M{"input_img_id": inputImgId}
	update := bson.M{
		"$set": bson.M{
			"userid":                        imageInfo.UserId,
			"inputimgid":                    imageInfo.InputImgId,
			"imagetype":                     imageInfo.ImageType,
			"calibrationaxesimg":            imageInfo.CalibrationImgAxes,
			"calibration_vanishingpointimg": imageInfo.CalibrationImgVanishingPoint,
			"outputimg":                     imageInfo.OutputImg,
			"outputkeypoints":               imageInfo.OutputKeypoints,
			"dtlgolfsetuppoints":            imageInfo.DtlGolfSetupPoints,
			"faceongolfsetuppoints":         imageInfo.FaceOnGolfSetupPoints,
		},
	}
	var updatedImageInfo cv.ImageInfo
	if err := d.imageInfoCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedImageInfo); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, fmt.Errorf("no imageinfos with inputimgid: %s", inputImgId)
		}
		return nil, fmt.Errorf("could not update imageinfo: %w", err)
	}
	fmt.Printf("Update image info result: %+v...\n", updatedImageInfo)
	return &updatedImageInfo, nil
}

// TODO: If delete image, delete all associated ImageInfos
func (d *DbManager) DeleteImageInfo(ctx context.Context, inputImgId string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Deleting image info inputimgid: %s...\n", inputImgId)
	filter := bson.M{"inputimgid": inputImgId}
	res, err := d.imageInfoCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("could not delete imageinfo %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("did not delete any imageinfos, imageinfo inputimgid %s may not exist", inputImgId)
	}
	fmt.Printf("Delete image result: inputimgid: %s\n", inputImgId)
	return nil
}

func (d *DbManager) DeleteUsersImageInfos() {

}
