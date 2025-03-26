package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Mentor struct {
	Name  string `bson:"name" json:"name"`
	Tel   string `bson:"tel" json:"tel"`
	Email string `bson:"email" json:"email"`
}

type User struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Score              int64              `bson:"score" json:"score"`
	Role               string             `bson:"role" json:"role"`
	Mentor             Mentor             `bson:"mentor" json:"mentor"`
	UserName           string             `bson:"user_name" json:"user_name" validate:"required"`
	CodeforcesUsername string             `bson:"codeforces_username" json:"codeforces_username" validate:"required"`
	PasswordHash       string             `bson:"password" json:"password" validate:"required"`
}

type Problem struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Author           string             `bson:"author" json:"author"`
	Title            string             `bson:"title" json:"title"`
	ProblemStatement string             `bson:"problem_statement" json:"problem_statement"`
	Source           string             `bson:"source" json:"source"`
	Difficulty       string             `bson:"difficulty" json:"difficulty"`
	ContestID        string             `bson:"contest_id" json:"contest_id"`
	Index            string             `json:"index"`
	Tags             []string           `bson:"tags" json:"tags"`
}

type Article struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Author   string             `bson:"author" json:"author"`
	Title    string             `bson:"title" json:"title"`
	Blog     string             `bson:"blog" json:"blog"`
	Tags     []string           `bson:"tags" json:"tags"`
	Problems []Problem          `bson:"problems" json:"problems"`
	Division string             `bson:"division" json:"division"`
}

type Submission struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     string             `bson:"user_id" json:"user_id"`
	ProblemID  string             `bson:"problem_id" json:"problem_id"`
	Submission string             `bson:"submission" json:"submission"`
}
