package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/AbenezerWork/AASTU-CPC/models"
	"github.com/AbenezerWork/AASTU-CPC/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("secret-key"))
var validate = validator.New()

type AuthController struct {
	UserRepo    *repository.UserRepository
	SessionRepo *repository.SessionRepository
}

func NewAuthController(ur *repository.UserRepository, sr *repository.SessionRepository) *AuthController {
	return &AuthController{
		UserRepo:    ur,
		SessionRepo: sr,
	}
}

func (ctrl *AuthController) Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.PasswordHash = string(hashedPassword)
	user.ID = primitive.NewObjectID()

	if err := ctrl.UserRepo.Create(context.Background(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.UserRepo.GetByUsername(context.Background(), credentials.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	//check to see if user is already logged in
	if sess, err := ctrl.SessionRepo.GetByUserID(context.Background(), user.ID.Hex()); err == nil || time.Since(sess.Issued) < time.Hour*24*7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already signed in"})
	}

	session, _ := store.Get(c.Request, "session-name")
	session.Values["authenticated"] = true
	session.Values["userID"] = user.ID
	session.ID = primitive.NewObjectID().Hex()
	session.Save(c.Request, c.Writer)

	sessionModel := &models.Session{
		SessionID: session.ID,
		UserID:    user.ID,
		IsAdmin:   user.Role == "admin" || user.Role == "root",
		Issued:    time.Now(),
	}
	if err := ctrl.SessionRepo.Create(context.Background(), sessionModel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "Set-Cookie",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour),
	})

	c.JSON(http.StatusOK, gin.H{"message": "Logged in"})
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	session.Values["authenticated"] = false
	session.Save(c.Request, c.Writer)

	if err := ctrl.SessionRepo.Delete(context.Background(), session.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

func (ctrl *AuthController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate.Struct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.PasswordHash = string(hashedPassword)
	user.ID = primitive.NewObjectID()

	if err := ctrl.UserRepo.Create(context.Background(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func (ctrl *AuthController) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := ctrl.UserRepo.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *AuthController) UpdateUser(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = id
	if err := ctrl.UserRepo.Update(context.Background(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (ctrl *AuthController) DeleteUser(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := ctrl.UserRepo.Delete(context.Background(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
