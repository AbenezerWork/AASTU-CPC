package controllers

import (
	"context"
	"net/http"

	"github.com/AbenezerWork/AASTU-CPC/models"
	"github.com/AbenezerWork/AASTU-CPC/repository"
	"github.com/AbenezerWork/AASTU-CPC/utils"
	"github.com/gin-gonic/gin"
)

// SubmissionController handles HTTP requests related to submissions.
type SubmissionController struct {
	Subrepo  *repository.SubmissionRepository
	Probrepo *repository.ProblemRepository
	Userrepo *repository.UserRepository
}

// NewSubmissionController initializes a new SubmissionController.
func NewSubmissionController(sr *repository.SubmissionRepository, pr *repository.ProblemRepository, ur *repository.UserRepository) *SubmissionController {
	return &SubmissionController{
		Subrepo:  sr,
		Probrepo: pr,
		Userrepo: ur,
	}
}

// ValidateSubmission handles POST /submissions/validate
// @Summary Validate a submission
// @Description Validate a user's submission for a problem
// @Tags Submissions
// @Accept json
// @Produce json
// @Param submission body models.Submission true "Submission data"
// @Success 200 {object} models.Submission
// @Router /submissions/validate [post]
func (sc *SubmissionController) ValidateSubmission(c *gin.Context) {
	var submission models.Submission

	if err := c.ShouldBindJSON(&submission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}

	problem, err := sc.Probrepo.GetByID(context.Background(), submission.ProblemID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
		return
	}

	if problem.Source == "codeforces" {
		user, err := sc.Userrepo.GetByID(context.Background(), submission.UserID)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error3": err.Error()})
			return
		}

		err, bl := utils.GetAndCheckAdmission(*problem, submission.Submission, user.CodeforcesUsername)

		if !bl {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		sc.Subrepo.Create(context.Background(), &submission)

	}

	//TODO: finish the submission checker for other platforms
}
