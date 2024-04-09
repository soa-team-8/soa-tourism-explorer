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

type EncounterRequestRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
}

func NewEncounterRequestRepository(client *mongo.Client) *EncounterRequestRepository {
	database := client.Database("encounters")
	collection := database.Collection("encounterRequests")
	return &EncounterRequestRepository{
		client:     client,
		database:   database,
		collection: collection,
		ctx:        context.Background(),
	}
}

func (repo *EncounterRequestRepository) Save(encounterReq model.EncounterRequest) (model.EncounterRequest, error) {
	nextID, err := repo.getNextSequence()
	if err != nil {
		return model.EncounterRequest{}, err
	}

	encounterReq.ID = uint64(nextID)
	_, err = repo.collection.InsertOne(repo.ctx, encounterReq)
	if err != nil {
		return model.EncounterRequest{}, err
	}

	return encounterReq, nil
}

func (repo *EncounterRequestRepository) FindAll() ([]model.EncounterRequest, error) {
	// Define an empty slice to hold the results
	var encounterRequests []model.EncounterRequest

	// Find all documents in the collection
	cursor, err := repo.collection.Find(repo.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(repo.ctx); err != nil {
			log.Printf("Error closing cursor: %v\n", err)
		}
	}()

	// Iterate over the cursor and decode documents into the slice
	for cursor.Next(repo.ctx) {
		var request model.EncounterRequest
		if err := cursor.Decode(&request); err != nil {
			return nil, err
		}
		encounterRequests = append(encounterRequests, request)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return encounterRequests, nil
}

func (repo *EncounterRequestRepository) FindByID(id int) (*model.EncounterRequest, error) {
	filter := bson.M{"id": id}
	var encounterRequest model.EncounterRequest
	err := repo.collection.FindOne(repo.ctx, filter).Decode(&encounterRequest)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("encounter request not found: %d", id)
		}
		return nil, err
	}

	return &encounterRequest, nil
}

func (repo *EncounterRequestRepository) Update(encounterReq model.EncounterRequest) (*model.EncounterRequest, error) {
	// Convert the ID to string
	//idStr := strconv.Itoa(int(encounterReq.ID))

	// Define the filter to find the document by ID
	filter := bson.M{"id": encounterReq.ID}

	// Define the update to set the new values
	update := bson.M{
		"$set": bson.M{
			"encounterid": encounterReq.EncounterId,
			"touristid":   encounterReq.TouristId,
			"status":      encounterReq.Status,
		},
	}

	// Perform the update operation
	_, err := repo.collection.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &encounterReq, nil
}

func (repo *EncounterRequestRepository) DeleteByID(id int) error {
	filter := bson.M{"id": id}

	_, err := repo.collection.DeleteOne(repo.ctx, filter)
	return err
}

func (repo *EncounterRequestRepository) getNextSequence() (int32, error) {
	filter := bson.M{"_id": "encounterRequestID"}
	update := bson.M{"$inc": bson.M{"seq": 1}}

	var result bson.M
	err := repo.collection.FindOneAndUpdate(repo.ctx, filter, update).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// If counter doesn't exist, initialize it
			counter := bson.M{"_id": "encounterRequestID", "seq": 2}
			_, err := repo.collection.InsertOne(repo.ctx, counter)
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
