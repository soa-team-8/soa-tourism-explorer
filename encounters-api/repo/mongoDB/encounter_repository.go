package mongoDB

import (
	"context"
	"encounters/model"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
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
}

func (r *EncounterRepository) FindByID(id uint64) (*model.Encounter, error) {
	filter := bson.M{"id": id}
	var encounter struct {
		Encounter bson.M `bson:"encounter"`
	}

	err := r.collection.FindOne(r.ctx, filter).Decode(&encounter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("encounter not found: %d", id)
		}
		return nil, err
	}

	var result model.Encounter

	if len(encounter.Encounter) > 0 {
		err = mapstructure.Decode(encounter.Encounter, &result)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.collection.FindOne(r.ctx, filter).Decode(&result)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, fmt.Errorf("encounter not found: %d", id)
			}
			return nil, err
		}
	}

	return &result, nil
}

func (r *EncounterRepository) FindAll() ([]model.Encounter, error) {
	var encounters []model.Encounter

	cursor, err := r.collection.Find(r.ctx, bson.M{"_id": bson.M{"$ne": "encounterID"}})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(r.ctx); err != nil {
			log.Printf("Error closing cursor: %v\n", err)
		}
	}()

	for cursor.Next(r.ctx) {
		var encounter struct {
			Encounter bson.M `bson:"encounter"`
		}
		if err := cursor.Decode(&encounter); err != nil {
			return nil, err
		}

		var result model.Encounter

		if len(encounter.Encounter) > 0 {
			err := mapstructure.Decode(encounter.Encounter, &result)
			if err != nil {
				return nil, err
			}
		} else {
			err := cursor.Decode(&result)
			if err != nil {
				return nil, err
			}
		}

		encounters = append(encounters, result)
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
	filter := bson.M{"id": encounter.ID}

	update := bson.M{"$set": bson.M{}}

	switch encounter.Type {
	case model.Location:
		update["$set"].(bson.M)["encounter"] = encounter
	case model.Social:
		update["$set"].(bson.M)["encounter"] = encounter
	default:
		update["$set"] = encounter
	}

	_, err := r.collection.UpdateOne(r.ctx, filter, update)
	if err != nil {
		return model.Encounter{}, err
	}

	return encounter, nil
}

// TODO usages???
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
