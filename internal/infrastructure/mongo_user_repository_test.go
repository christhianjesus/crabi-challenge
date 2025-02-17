package infrastructure

import (
	"context"
	"testing"
	"time"

	"github.com/christhianjesus/crabi-challenge/internal/domain"
	"github.com/christhianjesus/crabi-challenge/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type mongoUserRepositoryMock struct {
	collection *mocks.MongoCollection
	repo       MongoUserRepository
}

func setupMongoUserRepository(t *testing.T) *mongoUserRepositoryMock {
	mockMongoCollection := mocks.NewMongoCollection(t)

	return &mongoUserRepositoryMock{
		collection: mockMongoCollection,
		repo:       &mongoUserRepository{coll: mockMongoCollection},
	}
}

func TestNewMongoUserRepository_OK(t *testing.T) {
	mongoColl := &mongo.Collection{}

	md := mocks.NewMongoDatabase(t)
	md.On("Collection", mock.AnythingOfType("string")).Return(mongoColl)

	mur := NewMongoUserRepository(md)

	assert.NotNil(t, mur)
	assert.Equal(t, mongoColl, mur.(*mongoUserRepository).coll)
}

func TestCreate_InsertOneError(t *testing.T) {
	murm := setupMongoUserRepository(t)
	murm.collection.On("InsertOne", mock.IsType(nil), mock.AnythingOfType("*infrastructure.mongoUser")).Return(nil, assert.AnError)

	err := murm.repo.CreateUser(context.Context(nil), &domain.User{})

	assert.Error(t, err)
}

func TestCreate_InsertOneOK(t *testing.T) {
	murm := setupMongoUserRepository(t)
	murm.collection.On("InsertOne", mock.IsType(nil), mock.AnythingOfType("*infrastructure.mongoUser")).Return(nil, nil)

	err := murm.repo.CreateUser(context.Context(nil), &domain.User{})

	assert.NoError(t, err)
}

func TestGetIdAndHash_FindOneError(t *testing.T) {
	res := mongo.NewSingleResultFromDocument(nil, nil, nil)

	murm := setupMongoUserRepository(t)
	murm.collection.On("FindOne", mock.IsType(nil), mock.AnythingOfType("bson.M"), mock.AnythingOfType("*options.FindOneOptionsBuilder")).Return(res)

	id, hash, err := murm.repo.GetIdAndHash(context.Context(nil), "")

	assert.Error(t, err)
	assert.EqualError(t, err, mongo.ErrNilDocument.Error())
	assert.Empty(t, id)
	assert.Empty(t, hash)
}

func TestGetIdAndHash_FindOneErrorNoDocuments(t *testing.T) {
	res := mongo.NewSingleResultFromDocument(bson.M{}, mongo.ErrNoDocuments, nil)

	murm := setupMongoUserRepository(t)
	murm.collection.On("FindOne", mock.IsType(nil), mock.AnythingOfType("bson.M"), mock.AnythingOfType("*options.FindOneOptionsBuilder")).Return(res)

	id, hash, err := murm.repo.GetIdAndHash(context.Context(nil), "")

	assert.Error(t, err)
	assert.NotEqual(t, mongo.ErrNilDocument, err)
	assert.Empty(t, id)
	assert.Empty(t, hash)
}

func TestGetIdAndHash_FindOneOK(t *testing.T) {
	mongoID := bson.NewObjectIDFromTimestamp(time.Now())
	password := "XXXXXXXXXXXXX"
	res := mongo.NewSingleResultFromDocument(bson.M{"_id": mongoID, "password": password}, nil, nil)

	murm := setupMongoUserRepository(t)
	murm.collection.On("FindOne", mock.IsType(nil), mock.AnythingOfType("bson.M"), mock.AnythingOfType("*options.FindOneOptionsBuilder")).Return(res)

	id, hash, err := murm.repo.GetIdAndHash(context.Context(nil), "")

	assert.NoError(t, err)
	assert.Equal(t, mongoID.Hex(), id)
	assert.Equal(t, password, hash)
}

func TestGet_FindOneError(t *testing.T) {
	res := mongo.NewSingleResultFromDocument(nil, nil, nil)

	murm := setupMongoUserRepository(t)
	murm.collection.On("FindOne", mock.IsType(nil), mock.AnythingOfType("bson.M"), mock.AnythingOfType("*options.FindOneOptionsBuilder")).Return(res)

	user, err := murm.repo.GetUser(context.Context(nil), "")

	assert.Error(t, err)
	assert.EqualError(t, err, mongo.ErrNilDocument.Error())
	assert.Nil(t, user)
}

func TestGet_FindOneErrorNoDocuments(t *testing.T) {
	res := mongo.NewSingleResultFromDocument(bson.M{}, mongo.ErrNoDocuments, nil)

	murm := setupMongoUserRepository(t)
	murm.collection.On("FindOne", mock.IsType(nil), mock.AnythingOfType("bson.M"), mock.AnythingOfType("*options.FindOneOptionsBuilder")).Return(res)

	user, err := murm.repo.GetUser(context.Context(nil), "")

	assert.Error(t, err)
	assert.NotEqual(t, mongo.ErrNilDocument, err)
	assert.Nil(t, user)
}

func TestGet_FindOneOK(t *testing.T) {
	now := time.Now()
	mongoID := bson.NewObjectIDFromTimestamp(now)
	res := mongo.NewSingleResultFromDocument(bson.M{"_id": mongoID, "password": "XXXXXXXXXXXXX", "created_at": now}, nil, nil)

	murm := setupMongoUserRepository(t)
	murm.collection.On("FindOne", mock.IsType(nil), mock.AnythingOfType("bson.M"), mock.AnythingOfType("*options.FindOneOptionsBuilder")).Return(res)

	user, err := murm.repo.GetUser(context.Context(nil), "")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, mongoID.Hex(), user.ID)
	assert.Equal(t, now.UTC().Truncate(time.Millisecond), user.CreatedAt)
	assert.Empty(t, user.Password)
}
