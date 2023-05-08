package middleware

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(c *gin.Context) {
	firebaseAuth := c.MustGet("firebaseAuth").(*auth.Client)
	authorizationToken := c.GetHeader("Authorization")

	idToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
	if idToken == "" {
		c.Set("uid", idToken)
		c.Next()
		return
	}

	token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}
	c.Set("uid", token.UID)
	c.Next()
}
