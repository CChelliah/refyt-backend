package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/users/repo"
)

func GetSeller(userRepo *repo.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid := c.GetString("uid")

		if uid == "" {
			c.JSON(http.StatusOK, gin.H{"status": false})
			return
		}

		user, err := userRepo.FindUserByID(uid)

		switch {
		case errors.Is(err, repo.ErrUserNotFound):
			c.JSON(http.StatusNotFound, "user not found")
			return
		case err != nil && !errors.Is(err, repo.ErrUserNotFound):
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		switch {
		case *user.StripeConnectID == "":
			c.JSON(http.StatusOK, gin.H{"status": false})
			return
		case *user.StripeConnectID != "":
			c.JSON(http.StatusOK, gin.H{"status": true})
			return
		}
	}
}
