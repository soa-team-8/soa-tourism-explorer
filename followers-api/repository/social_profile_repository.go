package repositories

import (
	"context"
	"errors"
	"fmt"
	"followers/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"strconv"
)

type SocialProfileRepository struct {
	driver neo4j.DriverWithContext
}

func NewSocialProfileRepository() (*SocialProfileRepository, error) {
	uri := "bolt://localhost:7687"
	user := "neo4j"
	pass := "followers"
	auth := neo4j.BasicAuth(user, pass, "")

	driver, _ := neo4j.NewDriverWithContext(uri, auth)

	return &SocialProfileRepository{
		driver: driver,
	}, nil
}

func (repo *SocialProfileRepository) CheckConnection() error {
	ctx := context.Background()
	err := repo.driver.VerifyConnectivity(ctx)
	if err != nil {
		return fmt.Errorf("failed to verify connectivity: %w", err)
	}
	return nil
}

func (repo *SocialProfileRepository) CloseDriverConnection(ctx context.Context) error {
	err := repo.driver.Close(ctx)
	if err != nil {
		return fmt.Errorf("failed to close driver connection: %w", err)
	}
	return nil
}

func (repo *SocialProfileRepository) WriteUser(user *model.User) error {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "followers"})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Printf("Error closing session: %v", err)
		}
	}(session, ctx)

	result, err := session.Run(ctx,
		"CREATE (user:User) SET user.Id = $id, user.Username = $username RETURN user",
		map[string]interface{}{"id": user.ID, "username": user.Username})
	if err != nil {
		return fmt.Errorf("failed to execute user creation query: %w", err)
	}
	if !result.Next(ctx) {
		return errors.New("user creation operation did not return any result")
	}

	_, err = session.Run(ctx,
		"MATCH (user:User {Id: $id}) "+
			"CREATE (social:SocialProfile {UserId: $id, FollowersIds: [], FollowedIds: []}) "+
			"CREATE (user)-[:HAS_PROFILE]->(social)",
		map[string]interface{}{"id": user.ID})
	if err != nil {
		return fmt.Errorf("failed to create social profile node: %w", err)
	}

	return nil
}

func (repo *SocialProfileRepository) SaveUser(user *model.User) (bool, error) {
	if user.ID == 0 {
		id, err := repo.GenerateIncrementalID()
		if err != nil {
			return false, fmt.Errorf("failed to generate ID: %w", err)
		}
		user.ID = id
	}

	userInDatabase, err := repo.ReadUser(strconv.FormatUint(user.ID, 10))
	if err != nil {
		return false, fmt.Errorf("failed to read user: %w", err)
	}

	if userInDatabase == (model.User{}) {
		if err := repo.WriteUser(user); err != nil {
			return false, fmt.Errorf("failed to write user: %w", err)
		}
		return true, nil
	}

	return false, nil
}

func (repo *SocialProfileRepository) GenerateIncrementalID() (uint64, error) {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "followers"})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Printf("Error closing session: %v", err)
		}
	}(session, ctx)

	result, err := session.Run(ctx, "MATCH (user:User) RETURN max(user.Id) AS maxId", nil)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %w", err)
	}

	if result.Next(ctx) {
		record := result.Record()
		maxIDValue, ok := record.Get("maxId")
		if !ok {
			return 0, errors.New("missing max ID in result")
		}

		maxID, ok := maxIDValue.(int64)
		if !ok {
			return 0, errors.New("max ID is not of type int64")
		}

		newID := maxID + 1
		return uint64(newID), nil
	}

	return 1, nil
}

func (repo *SocialProfileRepository) ReadUser(userId string) (model.User, error) {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "followers"})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Printf("Error closing session: %v", err)
		}
	}(session, ctx)

	result, err := session.Run(ctx,
		"MATCH (u {Id: $id}) RETURN u.Id, u.Username",
		map[string]interface{}{"id": userId})
	if err != nil {
		return model.User{}, fmt.Errorf("failed to execute query: %w", err)
	}

	if !result.Next(ctx) {
		return model.User{}, nil
	}

	record := result.Record()
	idValue, ok := record.Get("u.Id")
	if !ok {
		return model.User{}, errors.New("missing user ID in result")
	}
	usernameValue, ok := record.Get("u.Username")
	if !ok {
		return model.User{}, errors.New("missing username in result")
	}

	id, err := strconv.ParseUint(fmt.Sprintf("%v", idValue), 10, 64)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to parse user ID: %w", err)
	}

	username, ok := usernameValue.(string)
	if !ok {
		return model.User{}, errors.New("username is not a string")
	}

	return model.User{ID: id, Username: username}, nil
}
