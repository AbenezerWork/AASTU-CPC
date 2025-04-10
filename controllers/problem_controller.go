package controllers

import (
	"context"
	"net/http"
	"strconv"

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

// CreateProblem handles POST /problemsedit
// @Summary Create a new problem
// @Description Create a new problem in the database NOTE: Don't enter the id
// @Tags Problems
// @Accept json
// @Produce json
// @Security Auth
// @Param problem body models.Problem true "Problem data"
// @Success 200 {object} models.Problem
// @Failure 401 {object} string "Unauthorized"
// @Router /problemsedit [post]
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

// UpdateProblem handles PUT /problemsedit/:id
// @Summary Update a problem
// @Description Update an existing problem by its ID NOTE: Don't update the id
// @Tags Problems
// @Accept json
// @Produce json
// @Security Auth
// @Param id path string true "Problem ID"
// @Param problem body models.Problem true "Updated problem data"
// @Success 200 {object} models.Problem
// @Failure 401 {object} string "Unauthorized"
// @Router /problemsedit/{id} [put]
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

// DeleteProblem handles DELETE /problemsedit/:id
// @Summary Delete a problem
// @Description Delete a problem by its ID
// @Tags Problems
// @Produce json
// @Security Auth
// @Param id path string true "Problem ID"
// @Success 200 {object} string "Problem deleted successfully"
// @Failure 401 {object} string "Unauthorized"
// @Router /problemsedit/{id} [delete]
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

// @Summary Get all problems
// @Description Retrieve all problems with pagination filters search and sort
// @Tags Problems
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search query"
// @Param sort query string false "Sort order"
// @Param maxRating query int false "Maximum problem rating"
// @Param minRating query int false "Minimum problem rating"
// @Success 200 {array} models.Problem
// @Router /problems [get]
func (ctrl *ProblemController) GetProblems(c *gin.Context) {
	// Parse query parameters
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	search := c.Query("search")
	sort := c.Query("sort")
	maxRating := 4000
	minRating := 0
	if maxRatingStr := c.Query("maxRating"); maxRatingStr != "" {
		var err error
		maxRating, err = strconv.Atoi(maxRatingStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid maxRating"})
			return
		}
	}
	if minRatingStr := c.Query("minRating"); minRatingStr != "" {
		var err error
		minRating, err = strconv.Atoi(minRatingStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid minRating"})
			return
		}
	}
	// Get problems with filters
	problems, err := ctrl.Repo.GetAllProblems(context.Background(), page, limit, search, sort, maxRating, minRating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, problems)
}
