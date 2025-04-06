package controllers

import (
	"context"
	"net/http"

	"github.com/AbenezerWork/AASTU-CPC/models"
	"github.com/AbenezerWork/AASTU-CPC/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProblemController handles HTTP requests related to problems.
type ProblemController struct {
	Repo *repository.ProblemRepository
}

// NewProblemController initializes a new ProblemController.
func NewProblemController(repo *repository.ProblemRepository) *ProblemController {
	return &ProblemController{Repo: repo}
}

// CreateProblem handles POST /problems
// @Summary Create a new problem
// @Description Create a new problem in the database
// @Tags Problems
// @Accept json
// @Produce json
// @Param problem body models.Problem true "Problem data"
// @Success 200 {object} models.Problem
// @Router /problems [post]
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

// GetProblemByID handles GET /problems/:id
// @Summary Get a problem by ID
// @Description Retrieve a problem by its ID
// @Tags Problems
// @Produce json
// @Param id path string true "Problem ID"
// @Success 200 {object} models.Problem
// @Router /problems/{id} [get]
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

// UpdateProblem handles PUT /problems/:id
// @Summary Update a problem
// @Description Update an existing problem by its ID
// @Tags Problems
// @Accept json
// @Produce json
// @Param id path string true "Problem ID"
// @Param problem body models.Problem true "Updated problem data"
// @Success 200 {object} models.Problem
// @Router /problems/{id} [put]
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

// DeleteProblem handles DELETE /problems/:id
// @Summary Delete a problem
// @Description Delete a problem by its ID
// @Tags Problems
// @Produce json
// @Param id path string true "Problem ID"
// @Router /problems/{id} [delete]
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
