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

type EncounterExecutionRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
	ctx        context.Context
}

func NewEncounterExecutionRepository(client *mongo.Client) *EncounterExecutionRepository {
	database := client.Database("encounters")
	collection := database.Collection("encounterExecutions")
	return &EncounterExecutionRepository{
		client:     client,
		database:   database,
		collection: collection,
		ctx:        context.Background(),
	}
}

func (r *EncounterExecutionRepository) Save(encounterExecution model.EncounterExecution) (model.EncounterExecution, error) {
	nextID, err := r.getNextSequence()
	if err != nil {
		return model.EncounterExecution{}, err
	}

	encounterExecution.ID = uint64(nextID)
	_, err = r.collection.InsertOne(r.ctx, encounterExecution)
	if err != nil {
		return model.EncounterExecution{}, err
	}

	return encounterExecution, nil
}

func (r *EncounterExecutionRepository) FindByID(id uint64) (*model.EncounterExecution, error) {
	filter := bson.M{"id": id}
	var encounterExecution model.EncounterExecution
	err := r.collection.FindOne(r.ctx, filter).Decode(&encounterExecution)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("encounter execution not found: %d", id)
		}
		return nil, err
	}

	return &encounterExecution, nil
}

func (r *EncounterExecutionRepository) FindAll() ([]model.EncounterExecution, error) {
	var encounterExecutions []model.EncounterExecution

	cursor, err := r.collection.Find(r.ctx, bson.M{"_id": bson.M{"$ne": "encounterExecutionID"}})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(r.ctx); err != nil {
			log.Printf("Error closing cursor: %v\n", err)
		}
	}()

	for cursor.Next(r.ctx) {
		var request model.EncounterExecution
		if err := cursor.Decode(&request); err != nil {
			return nil, err
		}
		encounterExecutions = append(encounterExecutions, request)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return encounterExecutions, nil
}

func (r *EncounterExecutionRepository) DeleteByID(id uint64) error {
	filter := bson.M{"id": id}

	_, err := r.collection.DeleteOne(r.ctx, filter)
	return err
}

func (r *EncounterExecutionRepository) Update(encounterExecution model.EncounterExecution) (model.EncounterExecution, error) {
	filter := bson.M{"id": encounterExecution.ID}

	update := bson.M{
		"$set": bson.M{
			"encounterid": encounterExecution.EncounterID,
			"encounter":   encounterExecution.Encounter,
			"touristid":   encounterExecution.TouristID,
			"status":      encounterExecution.Status,
			"starttime":   encounterExecution.StartTime,
			"endtime":     encounterExecution.EndTime,
		},
	}

	_, err := r.collection.UpdateOne(r.ctx, filter, update)
	if err != nil {
		return model.EncounterExecution{}, err
	}

	return encounterExecution, nil
}

func (r *EncounterExecutionRepository) FindAllByTourist(touristID uint64) ([]model.EncounterExecution, error) {
	var encounterExecutions []model.EncounterExecution

	filter := bson.M{"touristid": touristID}

	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Println("Error closing cursor:", err)
		}
	}(cursor, context.Background())

	for cursor.Next(context.Background()) {
		var encounterExecution model.EncounterExecution
		if err := cursor.Decode(&encounterExecution); err != nil {
			return nil, err
		}
		encounterExecutions = append(encounterExecutions, encounterExecution)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return encounterExecutions, nil
}

func (r *EncounterExecutionRepository) FindAllActiveByTourist(touristID uint64) ([]model.EncounterExecution, error) {
	var encounterExecutions []model.EncounterExecution

	// Define a filter to match documents with the given touristID and status Active
	filter := bson.M{
		"touristid": touristID,
		"status":    model.Active, // Assuming model.Active is the constant representing the active status
	}

	// Execute the find operation
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	// Use a deferred function to close the cursor and handle any errors
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			// Handle error when closing the cursor
			// This might indicate a problem with the MongoDB connection or server
			// You may log this error for debugging purposes
			log.Println("Error closing cursor:", err)
		}
	}(cursor, context.Background())

	// Iterate through the cursor and decode documents into EncounterExecution objects
	for cursor.Next(context.Background()) {
		var encounterExecution model.EncounterExecution
		if err := cursor.Decode(&encounterExecution); err != nil {
			return nil, err
		}
		encounterExecutions = append(encounterExecutions, encounterExecution)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return encounterExecutions, nil
}

func (r *EncounterExecutionRepository) FindAllCompletedByTourist(touristID uint64) ([]model.EncounterExecution, error) {
	var encounterExecutions []model.EncounterExecution

	// Define a filter to match documents with the given touristID and status Completed
	filter := bson.M{
		"touristid": touristID,
		"status":    model.Completed, // Assuming model.Completed is the constant representing the completed status
	}

	// Execute the find operation
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	// Use a deferred function to close the cursor and handle any errors
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			// Handle error when closing the cursor
			// This might indicate a problem with the MongoDB connection or server
			// You may log this error for debugging purposes
			log.Println("Error closing cursor:", err)
		}
	}(cursor, context.Background())

	// Iterate through the cursor and decode documents into EncounterExecution objects
	for cursor.Next(context.Background()) {
		var encounterExecution model.EncounterExecution
		if err := cursor.Decode(&encounterExecution); err != nil {
			return nil, err
		}
		encounterExecutions = append(encounterExecutions, encounterExecution)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return encounterExecutions, nil
}

func (r *EncounterExecutionRepository) FindByEncounter(encounterID uint64) (*model.EncounterExecution, error) {
	var encounterExecution model.EncounterExecution

	// Define a filter to match documents with the given encounterID
	filter := bson.M{"encounterid": encounterID}

	// Execute the find operation
	err := r.collection.FindOne(context.Background(), filter).Decode(&encounterExecution)
	if err != nil {
		return nil, err
	}

	return &encounterExecution, nil
}

func (r *EncounterExecutionRepository) FindAllByEncounter(encounterID uint64) ([]model.EncounterExecution, error) {
	var encounterExecutions []model.EncounterExecution

	// Define a filter to match documents with the given encounterID
	filter := bson.M{"encounterid": encounterID}

	// Execute the find operation
	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	// Use a deferred function to close the cursor and handle any errors
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			// Handle error when closing the cursor
			// This might indicate a problem with the MongoDB connection or server
			// You may log this error for debugging purposes
			log.Println("Error closing cursor:", err)
		}
	}(cursor, context.Background())

	// Iterate through the cursor and decode documents into EncounterExecution objects
	for cursor.Next(context.Background()) {
		var encounterExecution model.EncounterExecution
		if err := cursor.Decode(&encounterExecution); err != nil {
			return nil, err
		}
		encounterExecutions = append(encounterExecutions, encounterExecution)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return encounterExecutions, nil
}

func (r *EncounterExecutionRepository) FindAllByType(encounterID uint64, encounterType model.EncounterType) ([]model.EncounterExecution, error) {
	var encounterExecutions []model.EncounterExecution

	filter := bson.M{
		"encounterid":    encounterID,
		"encounter.type": encounterType,
	}

	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Println("Error closing cursor:", err)
		}
	}(cursor, context.Background())

	for cursor.Next(context.Background()) {
		var encounterExecution model.EncounterExecution
		if err := cursor.Decode(&encounterExecution); err != nil {
			return nil, err
		}
		encounterExecutions = append(encounterExecutions, encounterExecution)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return encounterExecutions, nil
}

func (r *EncounterExecutionRepository) FindByEncounterAndTourist(encounterID, touristID uint64) (*model.EncounterExecution, error) {
	var encounterExecution model.EncounterExecution

	// Define a filter to match documents with the given encounterID and touristID
	filter := bson.M{
		"encounterid": encounterID,
		"touristid":   touristID,
	}

	// Execute the find operation
	err := r.collection.FindOne(context.Background(), filter).Decode(&encounterExecution)
	if err != nil {
		return nil, err
	}

	return &encounterExecution, nil
}

func (r *EncounterExecutionRepository) UpdateRange(encounters []model.EncounterExecution) ([]model.EncounterExecution, error) {
	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.Background())

	// Start the transaction
	if err := session.StartTransaction(); err != nil {
		return nil, err
	}

	defer func() {
		if r := recover(); r != nil {
			// Rollback the transaction if panic occurs
			err := session.AbortTransaction(context.Background())
			if err != nil {
				return
			}
		}
	}()

	// Iterate through encounters and update each one
	for _, encounter := range encounters {
		filter := bson.M{"id": encounter.ID}

		update := bson.M{
			"$set": bson.M{
				"encounterid": encounter.EncounterID,
				"encounter":   encounter.Encounter,
				"touristid":   encounter.TouristID,
				"status":      encounter.Status,
				"starttime":   encounter.StartTime,
				"endtime":     encounter.EndTime,
			},
		}

		_, err := r.collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			err := session.AbortTransaction(context.Background())
			if err != nil {
				return nil, err
			}
			return nil, err
		}
	}

	// Commit the transaction
	if err := session.CommitTransaction(context.Background()); err != nil {
		return nil, err
	}

	return encounters, nil
}

func (r *EncounterExecutionRepository) getNextSequence() (int32, error) {
	filter := bson.M{"_id": "encounterExecutionID"}
	update := bson.M{"$inc": bson.M{"seq": 1}}

	var result bson.M
	err := r.collection.FindOneAndUpdate(r.ctx, filter, update).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			counter := bson.M{"_id": "encounterExecutionID", "seq": 2}
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
