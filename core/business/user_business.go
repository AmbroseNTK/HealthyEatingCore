package business

import (
	"context"
	"main/core/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserBusiness struct {
	DB *mongo.Database
}

func (b *UserBusiness) CreateIndexes() {
	b.DB.Collection("users").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	})
}

func (b *UserBusiness) GetOneById(id string) (models.UserProfile, error) {
	user := new(models.UserProfile)
	result := b.DB.Collection("users").FindOne(context.TODO(), map[string]interface{}{
		"id": id,
	})
	if result.Err() != nil {
		return *user, result.Err()
	}
	err := result.Decode(user)
	if err != nil {
		return *user, err
	}
	return *user, nil
}

func (b *UserBusiness) Create(user models.UserProfile) error {
	_, err := b.DB.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}

func (b *UserBusiness) Update(email string, userUpdated models.UpdatedUserProfile) error {

	updatedResult := b.DB.Collection("users").FindOneAndUpdate(context.TODO(),
		bson.M{"email": email},
		bson.M{"$set": userUpdated})
	if updatedResult.Err() != nil {
		return updatedResult.Err()
	}
	return nil
}

func (b *UserBusiness) Delete(email string) error {
	_, err := b.DB.Collection("users").DeleteOne(context.TODO(), bson.M{"email": email})
	if err != nil {
		return err
	}
	return nil
}
