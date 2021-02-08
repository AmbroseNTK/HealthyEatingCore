package business

import (
	"context"
	"encoding/json"
	"errors"
	"main/core/models"

	"github.com/cristalhq/jwt"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthBusiness struct {
	DB     *mongo.Database
	signer jwt.Signer
}

func NewAuthBusiness(DB *mongo.Database, key string) (*AuthBusiness, error) {
	signer, err := jwt.NewHS512([]byte(key))
	if err != nil {
		return &AuthBusiness{}, err
	}
	return &AuthBusiness{
		DB:     DB,
		signer: signer,
	}, nil
}

func (b *AuthBusiness) Register(user *models.UserAuth) error {
	result := b.DB.Collection("auth").FindOne(context.TODO(), map[string]interface{}{
		"email": user.Email,
	})
	if err := result.Decode(&map[string]interface{}{}); err == nil {
		return errors.New(user.Email + " has already existed")
	}

	hash, hashError := bcrypt.GenerateFromPassword([]byte(user.Password), 21)
	if hashError != nil {
		return errors.New("Invalid password")
	}
	hashString := string(hash)
	_, insertError := b.DB.Collection("auth").InsertOne(context.TODO(), map[string]interface{}{
		"email": user.Email,
		"hpass": hashString,
	})

	if insertError != nil {
		return errors.New("Cannot create account [" + user.Email + "]")
	}

	return nil

}

func (b *AuthBusiness) Login(user *models.UserAuth) (string, error) {
	result := b.DB.Collection("auth").FindOne(context.TODO(), map[string]interface{}{
		"email": user.Email,
	})
	userInDB := map[string]interface{}{}
	if err := result.Decode(userInDB); err != nil {
		return "", errors.New(user.Email + " did not exist")
	}
	compareError := bcrypt.CompareHashAndPassword([]byte(userInDB["hpass"].(string)), []byte(user.Password))
	if compareError != nil {
		return "", errors.New("Incorrect password")
	}

	profile := b.DB.Collection("users").FindOne(context.TODO(), map[string]interface{}{
		"email": user.Email,
	})

	profileInDB := map[string]interface{}{}
	profile.Decode(profileInDB)

	profileStr, _ := json.Marshal(profileInDB)

	token, err := b.signer.Sign(profileStr)
	if err != nil {
		return "", err
	}
	return string(token), nil
}
