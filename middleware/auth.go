package middleware

import (
	"context"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(c *gin.Context) {
	firebaseAuth := c.MustGet("firebaseAuth").(*auth.Client)
	authorizationToken := c.GetHeader("Authorization")

	idToken := strings.TrimSpace(strings.Replace(authorizationToken, "Bearer", "", 1))
	if idToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id token not available"})
		c.Abort()
		return
	}

	token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		fmt.Println("%s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}
	c.Set("uid", token.UID)
	c.Next()
}
