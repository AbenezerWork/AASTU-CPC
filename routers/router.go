package routers

import (
	"github.com/AbenezerWork/AASTU-CPC/controllers"
	"github.com/AbenezerWork/AASTU-CPC/middleware"
	"github.com/AbenezerWork/AASTU-CPC/repository"
	"github.com/gin-gonic/gin"
)

func SetupRouter(articleCtrl *controllers.ArticleController, problemCtrl *controllers.ProblemController, authCtrl *controllers.AuthController, sessionRepo *repository.SessionRepository, submissionController *controllers.SubmissionController) *gin.Engine {
	r := gin.Default()

	// Public routes
	r.GET("/articles/:id", articleCtrl.GetArticleByID)
	r.GET("problems/:id", problemCtrl.GetProblemByID)

	// Auth routes
	r.POST("/signup", authCtrl.Signup)
	r.POST("/login", authCtrl.Login)
	r.POST("/logout", authCtrl.Logout)

	//submission
	r.POST("validate-submission", submissionController.ValidateSubmission)

	// User routes
	users := r.Group("/users")
	users.Use(middleware.AdminAuthRequired(sessionRepo))
	{
		users.POST("/", authCtrl.CreateUser)
		users.GET("/:id", authCtrl.GetUserByID)
		users.PUT("/:id", authCtrl.UpdateUser)
		users.DELETE("/:id", authCtrl.DeleteUser)
	}

	// Problem routes
	problems := r.Group("/problemsedit")
	problems.Use(middleware.AuthRequired(sessionRepo))
	{
		problems.POST("/", problemCtrl.CreateProblem)
		problems.PUT("/:id", problemCtrl.UpdateProblem)
		problems.DELETE("/:id", problemCtrl.DeleteProblem)
	}
	articles := r.Group("/articlesedit")
	articles.Use(middleware.AuthRequired(sessionRepo))
	{
		articles.POST("/", articleCtrl.CreateArticle)
		articles.PUT("/:id", articleCtrl.UpdateArticle)
		articles.DELETE("/:id", articleCtrl.DeleteArticle)
	}

	return r
}
