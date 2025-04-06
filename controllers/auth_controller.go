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

// @Summary Signup a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.User true "User details"
// @Router /signup [post]
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

// @Summary Login a user
// @Description Authenticate a user and create a session
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.Credentials true "Login credentials"
// @Router /login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var credentials models.Credentials
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

// @Summary Logout a user
// @Description End the user's session
// @Tags auth
// @Produce json
// @Router /logout [post]
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

// @Summary Create a new user
// @Description Create a new user with the provided JSON body
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User details"
// @Router /users [post]
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

// @Summary Get a user by ID
// @Description Retrieve a user by their ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func (ctrl *AuthController) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := ctrl.UserRepo.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Update a user
// @Description Update an existing user's details
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.User true "Updated user details"
// @Router /users/{id} [put]
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

// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Router /users/{id} [delete]
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
