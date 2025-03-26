package repository

import (
	"context"

	"github.com/AbenezerWork/AASTU-CPC/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepository struct {
	Collection *mongo.Collection
}

func NewSessionRepository(db *mongo.Database) *SessionRepository {
	return &SessionRepository{
		Collection: db.Collection("sessions"),
	}
}

func (r *SessionRepository) Create(ctx context.Context, session *models.Session) error {
	_, err := r.Collection.InsertOne(ctx, session)
	return err
}

func (r *SessionRepository) GetBySessionID(ctx context.Context, sessionID string) (*models.Session, error) {
	var session models.Session
	err := r.Collection.FindOne(ctx, bson.M{"session_id": sessionID}).Decode(&session)
	return &session, err
}

func (r *SessionRepository) GetByUserID(ctx context.Context, userID string) (*models.Session, error) {
	var session models.Session
	err := r.Collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&session)
	return &session, err
}

func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"session_id": sessionID})
	return err
}
