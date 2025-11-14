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

// TODO: omitempty in .proto?
func (d *DbManager) CreateUser(ctx context.Context, user *cv.User) (*cv.User, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Creating user with name: %s...\n", user.UserName)
	res, err := d.userCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}
	objectId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("could create object id")
	}
	user.Id = objectId.Hex()
	fmt.Printf("Create user result: %+v\n", user)
	return user, nil
}

func (d *DbManager) ReadUser(ctx context.Context, id string) (*cv.User, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Reading user id: %s...\n", id)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("could not convert id to object id %w", err)
	}
	filter := bson.M{"_id": objectId}
	var user cv.User
	if err := d.userCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, fmt.Errorf("no users with id: %s", id)
		}
		return nil, fmt.Errorf("could not read user: %w", err)
	}
	user.Id = objectId.Hex()
	fmt.Printf("Read user result: %+v\n", user)
	return &user, nil
}

func (d *DbManager) UpdateUser(ctx context.Context, id string, user *cv.User) (*cv.User, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Updating user id: %s, username is %s...\n", id, user.UserName)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("could not convert id to object id %w", err)
	}
	filter := bson.M{"_id": objectId}
	update := bson.M{
		"$set": bson.M{
			"username": user.UserName,
			"email":    user.Email,
		},
	}
	var updatedUser cv.User
	if err := d.userCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedUser); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, fmt.Errorf("no users with id: %s", id)
		}
		return nil, fmt.Errorf("could not update user: %w", err)
	}
	updatedUser.Id = objectId.Hex()
	fmt.Printf("Update user result: %+v...\n", updatedUser)
	return &updatedUser, nil
}

// TODO: If delete user, delete associated images and image infos
func (d *DbManager) DeleteUser(ctx context.Context, id string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Deleting user id: %s...\n", id)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("could not convert id to object id %w", err)
	}
	filter := bson.M{"_id": objectId}
	res, err := d.userCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("could not delete user %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("did not delete any users, userid %s may not exist", id)
	}
	fmt.Printf("Delete user result: %+v...\n", res)
	return nil
}
