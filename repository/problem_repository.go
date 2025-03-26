package repository

import (
	"context"

	"github.com/AbenezerWork/AASTU-CPC/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProblemRepository struct {
	collection *mongo.Collection
}

func NewProblemRepository(db *mongo.Database) *ProblemRepository {
	return &ProblemRepository{
		collection: db.Collection("problems"),
	}
}

// Create creates a new problem
func (r *ProblemRepository) Create(ctx context.Context, problem *models.Problem) (*models.Problem, error) {
	result, err := r.collection.InsertOne(ctx, problem)
	if err != nil {
		return nil, err
	}

	problem.ID = result.InsertedID.(primitive.ObjectID)
	return problem, nil
}

// GetByID retrieves a problem by its ID
func (r *ProblemRepository) GetByID(ctx context.Context, id string) (*models.Problem, error) {
	id_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.Problem{}, err
	}

	var problem models.Problem
	err = r.collection.FindOne(ctx, bson.M{"_id": id_}).Decode(&problem)
	if err != nil {
		return nil, err
	}
	return &problem, nil
}

// Update updates an existing problem
func (r *ProblemRepository) Update(ctx context.Context, problem *models.Problem) error {
	_, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": problem.ID},
		problem,
	)
	return err
}

// Delete removes a problem by its ID
func (r *ProblemRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
