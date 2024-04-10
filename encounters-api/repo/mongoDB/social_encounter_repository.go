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

type SocialEncounterRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
}

func NewSocialEncounterRepository(client *mongo.Client) *SocialEncounterRepository {
	database := client.Database("encounters")
	collection := database.Collection("encounters")
	return &SocialEncounterRepository{
		client:     client,
		database:   database,
		collection: collection,
		ctx:        context.Background(),
	}
}

func (r *SocialEncounterRepository) Save(socialEncounter model.SocialEncounter) (model.SocialEncounter, error) {
	nextID, err := r.getNextSequence()
	if err != nil {
		return model.SocialEncounter{}, err
	}

	socialEncounter.ID = uint64(nextID)
	socialEncounter.Encounter.ID = uint64(nextID)
	_, err = r.collection.InsertOne(r.ctx, socialEncounter)
	if err != nil {
		return model.SocialEncounter{}, err
	}

	return socialEncounter, nil
}

func (r *SocialEncounterRepository) FindByID(id uint64) (*model.SocialEncounter, error) {
	filter := bson.M{"id": id}
	var encounter model.SocialEncounter
	err := r.collection.FindOne(r.ctx, filter).Decode(&encounter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("social encounter not found: %d", id)
		}
		return nil, err
	}
	return &encounter, nil
}

func (r *SocialEncounterRepository) Update(socialEncounter model.SocialEncounter) (model.SocialEncounter, error) {
	filter := bson.M{"id": socialEncounter.ID}

	// Define the update to set the new values
	update := bson.M{
		"$set": bson.M{
			"encounter":         socialEncounter.Encounter,
			"requiredpeople":    socialEncounter.RequiredPeople,
			"range":             socialEncounter.Range,
			"activetouristsids": socialEncounter.ActiveTouristsIds,
		},
	}

	// Perform the update operation
	_, err := r.collection.UpdateOne(r.ctx, filter, update)
	if err != nil {
		return model.SocialEncounter{}, err
	}

	return socialEncounter, nil
}

func (r *SocialEncounterRepository) DeleteByID(id uint64) error {
	filter := bson.M{"id": id}

	_, err := r.collection.DeleteOne(r.ctx, filter)
	return err
}

func (r *SocialEncounterRepository) FindAll() ([]model.SocialEncounter, error) {
	var encounters []model.SocialEncounter

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
		var encounter model.SocialEncounter
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

func (r *SocialEncounterRepository) getNextSequence() (int32, error) {
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
