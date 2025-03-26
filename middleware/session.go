package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AbenezerWork/AASTU-CPC/repository"
	"github.com/gin-gonic/gin"
)

func AdminAuthRequired(sessionRepo *repository.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		cookies := c.Request.Cookies()

		var session *http.Cookie

		for _, cookie := range cookies {
			fmt.Println("COOKIE", cookie.Value, cookie.Name)
			if cookie.Name == "Set-Cookie" {
				session = cookie
			}
		}
		fmt.Println("HELLO")
		if session == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		sessionModel, err := sessionRepo.GetBySessionID(context.Background(), session.Value)
		if !sessionModel.IsAdmin || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AuthRequired(sessionRepo *repository.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		cookies := c.Request.Cookies()

		var session *http.Cookie

		for _, cookie := range cookies {
			fmt.Println("COOKIE", cookie.Value, cookie.Name)
			if cookie.Name == "Set-Cookie" {
				session = cookie
			}
		}
		fmt.Println("HELLO")
		if session == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		sessionModel, err := sessionRepo.GetBySessionID(context.Background(), session.Value)
		fmt.Println(sessionModel.SessionID, session.Value)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Set("userID", sessionModel.UserID)
		c.Next()
	}
}
