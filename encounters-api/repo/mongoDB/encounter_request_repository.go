package mongoDB

import (
	"context"
	"encounters/model"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strconv"
)

type EncounterRequestRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
}

func New(client *mongo.Client) *EncounterRequestRepository {
	database := client.Database("encounters")
	collection := database.Collection("encounterRequests")
	return &EncounterRequestRepository{
		client:     client,
		database:   database,
		collection: collection,
		ctx:        context.Background(),
	}
}

func (repo *EncounterRequestRepository) AcceptRequest(id int) (*model.EncounterRequest, error) {
	// Find the document by ID
	filter := bson.M{"id": id}
	var requestToUpdate model.EncounterRequest
	err := repo.collection.FindOne(repo.ctx, filter).Decode(&requestToUpdate)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("not found %d", id)
		}
		return nil, err
	}

	// Update the status
	requestToUpdate.Accept()

	// Update the document in the collection
	update := bson.M{"$set": requestToUpdate}
	_, err = repo.collection.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &requestToUpdate, nil
}

func (repo *EncounterRequestRepository) RejectRequest(id int) (*model.EncounterRequest, error) {
	// Convert the ID to string
	idStr := strconv.Itoa(id)

	// Find the document by the custom id field
	filter := bson.M{"id": idStr}
	var requestToUpdate model.EncounterRequest
	err := repo.collection.FindOne(repo.ctx, filter).Decode(&requestToUpdate)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("not found %d", id)
		}
		return nil, err
	}

	// Update the status
	requestToUpdate.Reject()

	// Update the document in the collection
	update := bson.M{"$set": requestToUpdate}
	_, err = repo.collection.UpdateOne(repo.ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &requestToUpdate, nil
}

func (repo *EncounterRequestRepository) Save(encounterReq model.EncounterRequest) (model.EncounterRequest, error) {
	encounterReqBSON, err := bson.Marshal(encounterReq)
	if err != nil {
		return model.EncounterRequest{}, err
	}

	_, err = repo.collection.InsertOne(repo.ctx, encounterReqBSON)
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
			return nil, fmt.Errorf("encounter request not found: Not found %d", id)
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
	// Convert the ID to string
	//idStr := strconv.Itoa(id)

	// Define the filter to find the document by ID
	filter := bson.M{"id": id}

	// Perform the delete operation
	_, err := repo.collection.DeleteOne(repo.ctx, filter)
	return err
}
