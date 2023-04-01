package routes

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/users/repo"
)

func Get(userRepo *repo.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid := c.GetString("uid")

		if uid == "" {
			c.JSON(http.StatusUnauthorized, "unauthorized user")
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

		c.JSON(200, user)
	}
}
