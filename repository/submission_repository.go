package repository

import (
	"context"

	"github.com/AbenezerWork/AASTU-CPC/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubmissionRepository struct {
	Collection *mongo.Collection
}

func NewSubmissionRepository(db *mongo.Database) *SubmissionRepository {
	return &SubmissionRepository{
		Collection: db.Collection("submission"),
	}
}

func (r *SubmissionRepository) Create(ctx context.Context, submission *models.Submission) error {
	_, err := r.Collection.InsertOne(ctx, submission)
	return err
}

func (r *SubmissionRepository) GetByProblemID(ctx context.Context, userID string) (*models.Submission, error) {
	var submission models.Submission
	err := r.Collection.FindOne(ctx, bson.M{"problem_id": userID}).Decode(&submission)
	return &submission, err
}

func (r *SubmissionRepository) GetByUserID(ctx context.Context, userID string) (*models.Submission, error) {
	var submission models.Submission
	err := r.Collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&submission)
	return &submission, err
}

func (r *SubmissionRepository) Delete(ctx context.Context, submissionID string) error {
	id, err := primitive.ObjectIDFromHex(submissionID)
	if err != nil {
		return err
	}
	_, err = r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
