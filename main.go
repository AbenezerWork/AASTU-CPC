package main

import (
	"context"
	"fmt"
	"log"

	"github.com/AbenezerWork/AASTU-CPC/controllers"
	"github.com/AbenezerWork/AASTU-CPC/models"
	"github.com/AbenezerWork/AASTU-CPC/repository"
	"github.com/AbenezerWork/AASTU-CPC/routers"
	"github.com/AbenezerWork/AASTU-CPC/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://admin:037005@myatlasclusteredu.tknszig.mongodb.net/?retryWrites=true&w=majority&appName=myAtlasClusterEDU")

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("CHECK")
		log.Fatal(err)
	}
	fmt.Println("CHECK")

	db := client.Database("AASTU_CPC")

	articleRepo := repository.NewArticleRepository(db)
	authRepo := repository.NewUserRepository(db)
	problemRepo := repository.NewProblemRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	submissionRepo := repository.NewSubmissionRepository(db)

	articleCtrl := controllers.NewArticleController(articleRepo)
	authCtrl := controllers.NewAuthController(authRepo, sessionRepo)
	problemCtrl := controllers.NewProblemController(problemRepo)
	submissionCtrl := controllers.NewSubmissionController(submissionRepo, problemRepo, authRepo)

	//checking the cf request module
	err, bl := utils.GetAndCheckAdmission(models.Problem{ContestID: "1859", Index: "B"}, "310872613", "FunkyLlama")

	fmt.Println(err, bl)

	r := routers.SetupRouter(articleCtrl, problemCtrl, authCtrl, sessionRepo, submissionCtrl)
	r.Run(":8080")
}
