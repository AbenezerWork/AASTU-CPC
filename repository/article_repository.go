package repository

import (
	"context"

	"github.com/AbenezerWork/AASTU-CPC/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ArticleRepository struct {
	collection *mongo.Collection
}

func NewArticleRepository(db *mongo.Database) *ArticleRepository {
	return &ArticleRepository{
		collection: db.Collection("Articles"),
	}
}

// Create creates a new article
func (r *ArticleRepository) Create(ctx context.Context, article *models.Article) (*models.Article, error) {
	result, err := r.collection.InsertOne(ctx, article)
	if err != nil {
		return nil, err
	}

	article.ID = result.InsertedID.(primitive.ObjectID)
	return article, nil
}

// GetByID retrieves an article by its ID
func (r *ArticleRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Article, error) {
	var article models.Article
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&article)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// GetAll retrieves all articles
func (r *ArticleRepository) GetAll(ctx context.Context) ([]models.Article, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var articles []models.Article
	if err := cursor.All(ctx, &articles); err != nil {
		return nil, err
	}
	return articles, nil
}

// Update updates an existing article
func (r *ArticleRepository) Update(ctx context.Context, article *models.Article) error {
	_, err := r.collection.ReplaceOne(
		ctx,
		bson.M{"_id": article.ID},
		article,
	)
	return err
}

// Delete removes an article by its ID
func (r *ArticleRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// FindByTags finds articles by tags
func (r *ArticleRepository) FindByTags(ctx context.Context, tags []string) ([]models.Article, error) {
	filter := bson.M{"tags": bson.M{"$in": tags}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var articles []models.Article
	if err := cursor.All(ctx, &articles); err != nil {
		return nil, err
	}
	return articles, nil
}
