package business

import (
	"context"
	"main/core/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserBusiness struct {
	DB *mongo.Database
}

func (b *UserBusiness) Create(user models.UserProfile) error {
	_, err := b.DB.Collection("users").InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}
