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

type EncounterRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
}

func NewEncounterRepository(client *mongo.Client) *EncounterRepository {
	database := client.Database("encounters")
	collection := database.Collection("encounters")
	return &EncounterRepository{
		client:     client,
		database:   database,
		collection: collection,
		ctx:        context.Background(),
	}
}

func (r *EncounterRepository) Save(encounter model.Encounter) (model.Encounter, error) {
	nextID, err := r.getNextSequence()
	if err != nil {
		return model.Encounter{}, err
	}

	encounter.ID = uint64(nextID)
	_, err = r.collection.InsertOne(r.ctx, encounter)
	if err != nil {
		return model.Encounter{}, err
	}

	return encounter, nil

	/*
		encounterBSON, err := bson.Marshal(encounter)
		if err != nil {
			return model.Encounter{}, err
		}

		_, err = r.collection.InsertOne(r.ctx, encounterBSON)
		if err != nil {
			return model.Encounter{}, err
		}
		return encounter, nil
	*/
}

func (r *EncounterRepository) FindByID(id uint64) (*model.Encounter, error) {
	filter := bson.M{"id": id}
	var encounter model.Encounter
	err := r.collection.FindOne(r.ctx, filter).Decode(&encounter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("encounter not found: %d", id)
		}
		return nil, err
	}
	return &encounter, nil
}

func (r *EncounterRepository) FindAll() ([]model.Encounter, error) {
	var encounters []model.Encounter

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
		var encounter model.Encounter
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

func (r *EncounterRepository) DeleteByID(id uint64) error {
	filter := bson.M{"id": id}

	_, err := r.collection.DeleteOne(r.ctx, filter)
	return err
}

func (r *EncounterRepository) Update(encounter model.Encounter) (model.Encounter, error) {

	// Define the filter to find the document by ID
	filter := bson.M{"id": encounter.ID}

	// Define the update to set the new values
	update := bson.M{
		"$set": bson.M{
			"authorid":    encounter.AuthorID,
			"name":        encounter.Name,
			"description": encounter.Description,
			"xp":          encounter.XP,
			"status":      encounter.Status,
			"type":        encounter.Type,
			"longitude":   encounter.Longitude,
			"latitude":    encounter.Latitude,
		},
	}

	// Perform the update operation
	_, err := r.collection.UpdateOne(r.ctx, filter, update)
	if err != nil {
		return model.Encounter{}, err
	}

	return encounter, nil
}

func (r *EncounterRepository) FindByIds(ids []uint64) ([]model.Encounter, error) {
	var encounters []model.Encounter

	// Convert IDs to strings
	var stringIDs []string
	for _, id := range ids {
		stringIDs = append(stringIDs, fmt.Sprintf("%d", id))
	}

	// Define the filter to find documents by IDs
	filter := bson.M{"id": bson.M{"$in": stringIDs}}

	// Find documents matching the filter
	cursor, err := r.collection.Find(r.ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := cursor.Close(r.ctx); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Decode documents into encounters slice
	for cursor.Next(r.ctx) {
		var encounter model.Encounter
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

func (r *EncounterRepository) getNextSequence() (int32, error) {
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
