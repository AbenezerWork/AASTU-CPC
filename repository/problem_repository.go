package repository

import (
	"context"

	"github.com/AbenezerWork/AASTU-CPC/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// GetAllProblems retrieves all problems with pagination, search, and sort
func (r *ProblemRepository) GetAllProblems(ctx context.Context, page int, limit int, search string, sort string, maxRating int, minRating int) ([]models.Problem, error) {
	skip := (page - 1) * limit

	// Build filter for search
	filter := bson.M{}
	if search != "" {
		filter = bson.M{
			"$or": []bson.M{
				{"title": bson.M{"$regex": search, "$options": "i"}},
				{"description": bson.M{"$regex": search, "$options": "i"}},
			},
		}
	}

	filter["difficulty"] = bson.M{"$lte": maxRating}

	filter["difficulty"] = bson.M{"$gte": minRating}

	// Build sort options
	sortOptions := bson.D{}
	if sort != "" {
		if sort[0] == '-' {
			sortOptions = append(sortOptions, bson.E{Key: sort[1:], Value: -1})
		} else {
			sortOptions = append(sortOptions, bson.E{Key: sort, Value: 1})
		}
	}

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	if len(sortOptions) > 0 {
		findOptions.SetSort(sortOptions)
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var problems []models.Problem
	if err := cursor.All(ctx, &problems); err != nil {
		return nil, err
	}

	return problems, nil
}
