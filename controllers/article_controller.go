package controllers

import (
	"context"
	"net/http"

	"github.com/AbenezerWork/AASTU-CPC/models"
	"github.com/AbenezerWork/AASTU-CPC/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ArticleController struct {
	Repo *repository.ArticleRepository
}

func NewArticleController(repo *repository.ArticleRepository) *ArticleController {
	return &ArticleController{Repo: repo}
}

// @Summary Create a new article
// @Description Create a new article with the provided JSON body
// @Tags articles
// @Accept json
// @Produce json
// @Param article body models.Article true "Article to create"
// @Success 200 {object} models.Article
// @Router /articlesedit [post]
func (ctrl *ArticleController) CreateArticle(c *gin.Context) {
	var article models.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdArticle, err := ctrl.Repo.Create(context.Background(), &article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, createdArticle)
}

// @Summary Get an article by ID
// @Description Retrieve a single article by its ID
// @Tags articles
// @Produce json
// @Param id path string true "Article ID"
// @Success 200 {object} models.Article
// @Router /articles/{id} [get]
func (ctrl *ArticleController) GetArticleByID(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	article, err := ctrl.Repo.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, article)
}

// @Summary Update an article
// @Description Update an existing article by its ID
// @Tags articles
// @Accept json
// @Produce json
// @Param id path string true "Article ID"
// @Param article body models.Article true "Updated article data"
// @Success 200 {object} models.Article
// @Router /articlesedit/{id} [put]
func (ctrl *ArticleController) UpdateArticle(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var article models.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	article.ID = id
	if err := ctrl.Repo.Update(context.Background(), &article); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, article)
}

// @Summary Delete an article
// @Description Delete an article by its ID
// @Tags articles
// @Produce json
// @Param id path string true "Article ID"
// @Router /articlesedit/{id} [delete]
func (ctrl *ArticleController) DeleteArticle(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := ctrl.Repo.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Article deleted"})
}
