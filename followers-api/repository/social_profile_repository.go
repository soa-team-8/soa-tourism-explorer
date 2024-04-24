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
	uri := "bolt://followers_database:7687" //TODO: mrk
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

	_, err := session.Run(ctx,
		"CREATE (user:User) SET user.Id = $id, user.Username = $username RETURN user",
		map[string]interface{}{"id": user.ID, "username": user.Username})
	if err != nil {
		return fmt.Errorf("failed to execute user creation query: %w", err)
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

//	func (repo *SocialProfileRepository) Follow(userID, followedUserID uint64) error {
//		ctx := context.Background()
//		session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "followers"})
//		defer func(session neo4j.SessionWithContext, ctx context.Context) {
//			err := session.Close(ctx)
//			if err != nil {
//				log.Printf("Error closing session: %v", err)
//			}
//		}(session, ctx)
//
//		_, err := session.Run(ctx,
//			"MATCH (user:User {Id: $userID}), (followedUser:User {Id: $followedUserID}) "+
//				"CREATE (user)-[:FOLLOWS]->(followedUser)",
//			map[string]interface{}{"userID": userID, "followedUserID": followedUserID})
//		if err != nil {
//			return fmt.Errorf("failed to execute follow query: %w", err)
//		}
//
//		return nil
//	}
func (repo *SocialProfileRepository) Follow(userID, followedUserID uint64) error {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "followers"})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Printf("Error closing session: %v", err)
		}
	}(session, ctx)

	// Check if the follow relationship already exists
	result, err := session.Run(ctx,
		"MATCH (user:User {Id: $userID})-[:FOLLOWS]->(followedUser:User {Id: $followedUserID}) RETURN COUNT(*)",
		map[string]interface{}{"userID": userID, "followedUserID": followedUserID})
	if err != nil {
		return fmt.Errorf("failed to execute follow check query: %w", err)
	}

	if result.Next(ctx) {
		countValue, ok := result.Record().Get("COUNT(*)")
		if !ok {
			return fmt.Errorf("failed to get count of existing follow relationship")
		}
		count := countValue.(int64)
		if count > 0 {
			// Follow relationship already exists, no need to create a new one
			return nil
		}
	}

	// Create the follow relationship
	_, err = session.Run(ctx,
		"MATCH (user:User {Id: $userID}), (followedUser:User {Id: $followedUserID}) "+
			"CREATE (user)-[:FOLLOWS]->(followedUser)",
		map[string]interface{}{"userID": userID, "followedUserID": followedUserID})
	if err != nil {
		return fmt.Errorf("failed to execute follow query: %w", err)
	}

	return nil
}
func (repo *SocialProfileRepository) Unfollow(userID, followedUserID uint64) error {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "followers"})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Printf("Error closing session: %v", err)
		}
	}(session, ctx)

	_, err := session.Run(ctx,
		"MATCH (user:User {Id: $userID})-[r:FOLLOWS]->(followedUser:User {Id: $followedUserID}) "+
			"DELETE r",
		map[string]interface{}{"userID": userID, "followedUserID": followedUserID})
	if err != nil {
		return fmt.Errorf("failed to execute unfollow query: %w", err)
	}

	return nil
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
			return 1, nil
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

func (repo *SocialProfileRepository) GetAllUsers(excludeID uint64) ([]*uint64, error) {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "followers"})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Printf("Error closing session: %v", err)
		}
	}(session, ctx)

	query := fmt.Sprintf("MATCH (user:User) WHERE user.Id <> $excludeID RETURN user.Id")
	result, err := session.Run(ctx, query, map[string]interface{}{"excludeID": int64(excludeID)})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var userIDs []*uint64
	for result.Next(ctx) {
		record := result.Record()
		idValue, ok := record.Get("user.Id")
		if !ok {
			return nil, errors.New("missing user ID in result")
		}

		idFloat, ok := idValue.(int64)
		if !ok {
			return nil, errors.New("user ID is not of type int64")
		}

		id := uint64(idFloat)
		userIDs = append(userIDs, &id)
	}

	return userIDs, nil
}

func (repo *SocialProfileRepository) GetSocialProfile(userID uint64) (*model.SocialProfile, error) {
	socialProfile := &model.SocialProfile{
		UserId:     userID,
		Followers:  []*model.User{},
		Followed:   []*model.User{},
		Followable: []*model.User{},
	}

	// Get followers IDs
	followerIDs, err := repo.GetFollowerIDs(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get followers: %w", err)
	}

	// Get followed IDs
	followedIDs, err := repo.GetFollowedIDs(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get followed users: %w", err)
	}

	// Get followable IDs
	followableIDs, err := repo.GetFollowableIDs(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get followable users: %w", err)
	}

	// Fetch usernames for each user ID and create model.User objects
	for _, id := range followerIDs {
		username, err := repo.GetUsernameByID(*id)
		if err != nil {
			return nil, fmt.Errorf("failed to get username for follower ID %d: %w", *id, err)
		}
		socialProfile.Followers = append(socialProfile.Followers, &model.User{ID: *id, Username: username})
	}

	for _, id := range followedIDs {
		username, err := repo.GetUsernameByID(*id)
		if err != nil {
			return nil, fmt.Errorf("failed to get username for followed ID %d: %w", *id, err)
		}
		socialProfile.Followed = append(socialProfile.Followed, &model.User{ID: *id, Username: username})
	}

	for _, id := range followableIDs {
		username, err := repo.GetUsernameByID(*id)
		if err != nil {
			return nil, fmt.Errorf("failed to get username for followable ID %d: %w", *id, err)
		}
		socialProfile.Followable = append(socialProfile.Followable, &model.User{ID: *id, Username: username})
	}

	return socialProfile, nil
}

func (repo *SocialProfileRepository) GetUsernameByID(userID uint64) (string, error) {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "followers"})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Printf("Error closing session: %v", err)
		}
	}(session, ctx)

	result, err := session.Run(ctx,
		"MATCH (user:User {Id: $userID}) RETURN user.Username",
		map[string]interface{}{"userID": userID})
	if err != nil {
		return "", fmt.Errorf("failed to execute query to get username: %w", err)
	}

	if !result.Next(ctx) {
		return "", errors.New("user not found or username missing")
	}

	record := result.Record()
	usernameValue, ok := record.Get("user.Username")
	if !ok {
		return "", errors.New("missing username in result")
	}

	username, ok := usernameValue.(string)
	if !ok {
		return "", errors.New("username is not a string")
	}

	return username, nil
}

func (repo *SocialProfileRepository) GetFollowerIDs(userID uint64) ([]*uint64, error) {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "followers"})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Printf("Error closing session: %v", err)
		}
	}(session, ctx)

	result, err := session.Run(ctx,
		"MATCH (follower:User)-[:FOLLOWS]->(user:User {Id: $userID}) RETURN follower.Id",
		map[string]interface{}{"userID": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query to get followers: %w", err)
	}

	var followers []*uint64
	for result.Next(ctx) {
		record := result.Record()
		idValue, ok := record.Get("follower.Id")
		if !ok {
			return nil, errors.New("missing follower ID in result")
		}

		idFloat, ok := idValue.(int64)
		if !ok {
			return nil, errors.New("follower ID is not of type int64")
		}

		id := uint64(idFloat)
		followers = append(followers, &id)
	}

	return followers, nil
}

func (repo *SocialProfileRepository) GetFollowedIDs(userID uint64) ([]*uint64, error) {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "followers"})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Printf("Error closing session: %v", err)
		}
	}(session, ctx)

	result, err := session.Run(ctx,
		"MATCH (user:User {Id: $userID})-[:FOLLOWS]->(followed:User) RETURN followed.Id",
		map[string]interface{}{"userID": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query to get followed users: %w", err)
	}

	var followed []*uint64
	for result.Next(ctx) {
		record := result.Record()
		idValue, ok := record.Get("followed.Id")
		if !ok {
			return nil, errors.New("missing followed user ID in result")
		}

		idFloat, ok := idValue.(int64)
		if !ok {
			return nil, errors.New("followed user ID is not of type int64")
		}

		id := uint64(idFloat)
		followed = append(followed, &id)
	}

	return followed, nil
}

func (repo *SocialProfileRepository) GetFollowableIDs(userID uint64) ([]*uint64, error) {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "followers"})
	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			log.Printf("Error closing session: %v", err)
		}
	}(session, ctx)

	result, err := session.Run(ctx,
		"MATCH (user:User {Id: $userID}) "+
			"MATCH (followable:User) WHERE NOT (user)-[:FOLLOWS]->(followable) AND user <> followable "+
			"RETURN followable.Id",
		map[string]interface{}{"userID": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query to get followable users: %w", err)
	}

	var followable []*uint64
	for result.Next(ctx) {
		record := result.Record()
		idValue, ok := record.Get("followable.Id")
		if !ok {
			return nil, errors.New("missing followable user ID in result")
		}

		idFloat, ok := idValue.(int64)
		if !ok {
			return nil, errors.New("followable user ID is not of type int64")
		}

		id := uint64(idFloat)
		followable = append(followable, &id)
	}

	return followable, nil
}

func (repo *SocialProfileRepository) GetRecommendations(userID uint64) ([]*model.User, error) {
	recommendations := []*model.User{}

	followedIDs, err := repo.GetFollowedIDs(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get followed users: %w", err)
	}

	newFollowedIDs := make(map[uint64]struct{})
	for _, id := range followedIDs {
		ids, err := repo.GetFollowedIDs(*id)
		if err != nil {
			return nil, fmt.Errorf("failed to get followed users: %w", err)
		}
		for _, newID := range ids {
			newFollowedIDs[*newID] = struct{}{}
		}
	}

	for id := range newFollowedIDs {
		username, err := repo.GetUsernameByID(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get username for followed ID %d: %w", id, err)
		}
		recommendations = append(recommendations, &model.User{ID: id, Username: username})
	}

	return recommendations, nil
}
