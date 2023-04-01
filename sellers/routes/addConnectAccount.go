package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/sellers/repo"
	stripeGateway "refyt-backend/sellers/stripe"
)

func AddSellerAccount(sellerRepo *repo.SellerRepository, stripeKey string) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid := c.GetString("uid")

		fmt.Println("")

		if uid == "" {
			c.JSON(http.StatusUnauthorized, "unauthorized user")
			return
		}

		seller, err := sellerRepo.FindSeller(uid)

		switch {
		case seller.ConnectAccountID != "":
			c.JSON(http.StatusConflict, "seller account already exists")
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		account, err := stripeGateway.CreateSellerAccount(seller, stripeKey)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		seller.AddAccount(account)

		seller, err = sellerRepo.UpdateSeller(seller)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		accountLink, err := stripeGateway.CreateAccountLink(seller, stripeKey)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(200, accountLink)
	}
}
