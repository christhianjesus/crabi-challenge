package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/christhianjesus/crabi-challenge/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoUserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetIdAndHash(ctx context.Context, email string) (string, string, error)
	Get(ctx context.Context, userID string) (*domain.User, error)
}

type mongoUserRepository struct {
	coll mongoCollection
}

type mongoUser struct {
	ID        bson.ObjectID `bson:"_id"`
	Email     string        `bson:"email"`
	Password  string        `bson:"password"`
	FirstName string        `bson:"first_name,omitempty"`
	LastName  string        `bson:"last_name,omitempty"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

// interface added for testing purposes
type mongoDatabase interface {
	Collection(name string, opts ...options.Lister[options.CollectionOptions]) *mongo.Collection
}

// interface added for testing purposes
type mongoCollection interface {
	FindOne(ctx context.Context, filter interface{}, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult
	InsertOne(ctx context.Context, document interface{}, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error)
}

func NewMongoUserRepository(db mongoDatabase) MongoUserRepository {
	return &mongoUserRepository{coll: db.Collection("user")}
}

func (r *mongoUserRepository) Create(ctx context.Context, user *domain.User) error {
	currentTime := time.Now()
	mongoUser := &mongoUser{
		ID:        bson.NewObjectIDFromTimestamp(currentTime),
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	_, err := r.coll.InsertOne(ctx, mongoUser)

	return err
}

func (r *mongoUserRepository) GetIdAndHash(ctx context.Context, email string) (string, string, error) {
	opts := options.FindOne().SetProjection(bson.M{"password": 1})

	var user mongoUser

	err := r.coll.FindOne(ctx, bson.M{"email": email}, opts).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", "", errors.New("Not found")
		}

		return "", "", err
	}

	return user.ID.Hex(), user.Password, nil
}

func (r *mongoUserRepository) Get(ctx context.Context, userID string) (*domain.User, error) {
	opts := options.FindOne().SetProjection(bson.M{"password": 0})
	mongoID, _ := bson.ObjectIDFromHex(userID)

	var user mongoUser

	err := r.coll.FindOne(ctx, bson.M{"_id": mongoID}, opts).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Not found")
		}

		return nil, err
	}

	return &domain.User{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
