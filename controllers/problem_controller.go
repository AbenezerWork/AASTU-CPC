package controllers

import (
	"context"
	"net/http"

	"github.com/AbenezerWork/AASTU-CPC/models"
	"github.com/AbenezerWork/AASTU-CPC/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProblemController struct {
	Repo *repository.ProblemRepository
}

func NewProblemController(repo *repository.ProblemRepository) *ProblemController {
	return &ProblemController{Repo: repo}
}

func (ctrl *ProblemController) CreateProblem(c *gin.Context) {

	var problem models.Problem
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}
	createdProblem, err := ctrl.Repo.Create(context.Background(), &problem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}
	c.JSON(http.StatusOK, createdProblem)
}

func (ctrl *ProblemController) GetProblemByID(c *gin.Context) {
	session := c.MustGet("session").(models.Session)
	if session.UserID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := c.Param("id")
	problem, err := ctrl.Repo.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, problem)
}

func (ctrl *ProblemController) UpdateProblem(c *gin.Context) {
	session := c.MustGet("session").(models.Session)
	if session.UserID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var problem models.Problem
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	problem.ID = id
	if err := ctrl.Repo.Update(context.Background(), &problem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, problem)
}

func (ctrl *ProblemController) DeleteProblem(c *gin.Context) {
	session := c.MustGet("session").(models.Session)
	if session.UserID == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := ctrl.Repo.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Problem deleted"})
}
