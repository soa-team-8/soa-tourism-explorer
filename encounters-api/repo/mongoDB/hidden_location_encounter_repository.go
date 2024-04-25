package mongoDB

import (
	"context"
	"encounters/model"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type HiddenLocationRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
}

func NewHiddenLocationRepository(client *mongo.Client) *HiddenLocationRepository {
	database := client.Database("encounters")
	collection := database.Collection("encounters")
	return &HiddenLocationRepository{
		client:     client,
		database:   database,
		collection: collection,
		ctx:        context.Background(),
	}
}

func (r *HiddenLocationRepository) Save(hiddenLocationEncounter model.HiddenLocationEncounter) (model.HiddenLocationEncounter, error) {
	nextID, err := r.getNextSequence()
	if err != nil {
		return model.HiddenLocationEncounter{}, err
	}

	hiddenLocationEncounter.ID = uint64(nextID)
	hiddenLocationEncounter.Encounter.ID = uint64(nextID)
	_, err = r.collection.InsertOne(r.ctx, hiddenLocationEncounter)
	if err != nil {
		return model.HiddenLocationEncounter{}, err
	}

	return hiddenLocationEncounter, nil
}

func (r *HiddenLocationRepository) FindById(id uint64) (*model.HiddenLocationEncounter, error) {
	filter := bson.M{"id": id}
	var encounter model.HiddenLocationEncounter
	err := r.collection.FindOne(r.ctx, filter).Decode(&encounter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("hidden location encounter not found: %d", id)
		}
		return nil, err
	}
	return &encounter, nil
}

func (r *HiddenLocationRepository) Update(hiddenLocationEncounter model.HiddenLocationEncounter) (model.HiddenLocationEncounter, error) {
	filter := bson.M{"id": hiddenLocationEncounter.ID}

	// Define the update to set the new values
	update := bson.M{
		"$set": bson.M{
			"encounter":         hiddenLocationEncounter.Encounter,
			"locationlongitude": hiddenLocationEncounter.LocationLongitude,
			"locationlatitude":  hiddenLocationEncounter.LocationLatitude,
			"range":             hiddenLocationEncounter.Range,
		},
	}

	// Perform the update operation
	_, err := r.collection.UpdateOne(r.ctx, filter, update)
	if err != nil {
		return model.HiddenLocationEncounter{}, err
	}

	return hiddenLocationEncounter, nil
}

func (r *HiddenLocationRepository) DeleteById(id uint64) error {
	filter := bson.M{"id": id}

	_, err := r.collection.DeleteOne(r.ctx, filter)
	return err
}

func (r *HiddenLocationRepository) FindAll() ([]model.HiddenLocationEncounter, error) {
	var encounters []model.HiddenLocationEncounter

	// Find all documents in the collection
	cursor, err := r.collection.Find(r.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(r.ctx); err != nil {
			log.Printf("Error closing cursor: %v\n", err)
		}
	}()

	// Iterate over the cursor and decode documents into the slice
	for cursor.Next(r.ctx) {
		var encounter model.HiddenLocationEncounter
		if err := cursor.Decode(&encounter); err != nil {
			return nil, err
		}
		encounters = append(encounters, encounter)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return encounters, nil
}

func (r *HiddenLocationRepository) getNextSequence() (int32, error) {
	filter := bson.M{"_id": "encounterID"}
	update := bson.M{"$inc": bson.M{"seq": 1}}

	var result bson.M
	err := r.collection.FindOneAndUpdate(r.ctx, filter, update).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// If counter doesn't exist, initialize it
			counter := bson.M{"_id": "encounterID", "seq": 2}
			_, err := r.collection.InsertOne(r.ctx, counter)
			if err != nil {
				return 0, err
			}
			return 1, nil
		}
		return 0, err
	}

	seq := result["seq"].(int32)
	return seq, nil
}
