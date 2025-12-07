package db

import (
	"context"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Username string             `bson:"username,omitempty"`
	Password string             `bson:"password,omitempty"`
	Email    string             `bson:"email,omitempty"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %w", err)
	}
	return string(bytes), nil
}

func VerifyPasswordHash(hash string, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}

func UpdateUserFields(oldUser *User, newUser *User) *User {
	oldReflectVal := reflect.ValueOf(oldUser).Elem()
	newReflectVal := reflect.ValueOf(newUser).Elem()
	numFields := newReflectVal.NumField()
	for i := 0; i < numFields; i++ {
		newField := newReflectVal.Field(i)
		if !newField.IsZero() {
			oldField := oldReflectVal.Field(i)
			oldField.Set(newField)
		}
	}
	return oldUser
}

func (d *DbManager) CreateUser(ctx context.Context, user *User) (*User, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Creating user with name: %s...\n", user.Username)
	res, err := d.userCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}
	objectId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("could create object id")
	}
	user.Id = objectId
	fmt.Printf("Create user result: %+v\n", user)
	return user, nil
}

func (d *DbManager) ReadUser(ctx context.Context, userId string) (*User, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Reading user id: %s...\n", userId)
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("could not convert id to object id %w", err)
	}
	filter := bson.M{"_id": objectId}
	var user User
	if err := d.userCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, fmt.Errorf("no users with id: %s", userId)
		}
		return nil, fmt.Errorf("could not read user: %w", err)
	}
	user.Id = objectId
	fmt.Printf("Read user result: %+v\n", user)
	return &user, nil
}

func (d *DbManager) ReadUserFromUsername(ctx context.Context, userName string) (*User, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Reading user username: %s...\n", userName)
	filter := bson.M{"username": userName}
	var user User
	if err := d.userCollection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, fmt.Errorf("no users with name: %s", userName)
		}
		return nil, fmt.Errorf("could not read user: %w", err)
	}
	/*objectId, err := primitive.ObjectIDFromHex(user.Id.Hex())
	if err != nil {
		return nil, fmt.Errorf("could not convert id to object id %w", err)
	}
	user.Id = objectId*/
	fmt.Printf("Read user from username result: %+v\n", user)
	return &user, nil
}

func (d *DbManager) UpdateUser(ctx context.Context, userId string, user *User) (*User, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Updating user id: %s, username is %s...\n", userId, user.Username)
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("could not convert id to object id %w", err)
	}
	filter := bson.M{"_id": objectId}
	update := bson.M{
		"$set": bson.M{
			"username": user.Username,
			"password": user.Password,
			"email":    user.Email,
		},
	}
	var updatedUser User
	if err := d.userCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedUser); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, fmt.Errorf("no users with id: %s", userId)
		}
		return nil, fmt.Errorf("could not update user: %w", err)
	}
	updatedUser.Id = objectId
	fmt.Printf("Update user result: %+v...\n", updatedUser)
	return &updatedUser, nil
}

// Deletes all input images associated with user (deleteInputImageHelper will also delete golf keypoint associated with each input image)
// Then deletes the user
func (d *DbManager) DeleteUser(ctx context.Context, userId string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	fmt.Printf("Deleting user id: %s...\n", userId)
	// read input img ids associated with user
	inputImages, err := d.readInputImagesForUserHelper(ctx, userId)
	if err != nil {
		return fmt.Errorf("could not read input images associated with user %s: %w", userId, err)
	}
	// delete input imgs associated with user
	for _, inputImg := range inputImages {
		err = d.deleteInputImageHelper(ctx, inputImg.Id.Hex())
		if err != nil {
			return fmt.Errorf("could not delete input image %s associated with user %s: %w", inputImg.Id.Hex(), userId, err)
		}
	}
	// delete user
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return fmt.Errorf("could not convert id to object id %w", err)
	}
	filter := bson.M{"_id": objectId}
	res, err := d.userCollection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("could not delete user %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("did not delete any users, userid %s may not exist", userId)
	}
	fmt.Printf("Delete user result: %+v...\n", res)
	return nil
}
