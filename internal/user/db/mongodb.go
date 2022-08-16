package db

import (
	"context"
	"fmt"
	"go_advanced/internal/user"
	"go_advanced/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	d.logger.Debug("Create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %v", err)
	}

	d.logger.Debug("Convert InsertedID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("failed to convert object id to hex: %v", err)
}

func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return u, fmt.Errorf("failed to convert hex to ObjectID: %v", id)
	}

	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)

	if result.Err() != nil {
		return u, fmt.Errorf("failed to find user by id: %s due to error: %v", id, err)
	}

	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user from bd by id: %s due to error: %v", id, err)
	}

	return u, nil
}

func (d *db) FindAll(ctx context.Context) (u []user.User, err error) {

	result, err := d.collection.Find(ctx, bson.M{})

	if result.Err() != nil {
		return u, fmt.Errorf("failed to find users due to error: %v", err)
	}

	if err = result.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed to read all documents: %v", err)
	}

	return u, nil
}

func (d *db) Update(ctx context.Context, user user.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)

	if err != nil {
		return fmt.Errorf("failed to convert user id to ObjectID: %v", err)
	}
	filter := bson.M{"_id": objectID}

	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user. error: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user bytes. %v", err)
	}

	delete(updateUserObj, "_id")

	update := bson.M{
		"$set": updateUserObj,
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update user query: %v", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("not found")
	}

	d.logger.Tracef("Matched %d documents and modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}

func (d *db) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return fmt.Errorf("failed to convert user id to ObjectID: %s", id)
	}
	filter := bson.M{"_id": objectID}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("not found")
	}

	d.logger.Tracef("Deleted %d documents", result.DeletedCount)

	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {
	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
