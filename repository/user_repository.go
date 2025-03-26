package repository

import (
	"context"

	"github.com/AbenezerWork/AASTU-CPC/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: db.Collection("users"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := r.Collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	id_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.User{}, err
	}
	var user models.User
	err = r.Collection.FindOne(ctx, bson.M{"_id": id_}).Decode(&user)
	return &user, err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.Collection.FindOne(ctx, bson.M{"user_name": username}).Decode(&user)
	return &user, err
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	_, err := r.Collection.ReplaceOne(ctx, bson.M{"_id": user.ID}, user)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
