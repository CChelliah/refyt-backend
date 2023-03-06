package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"trading-card-app-backend/sellers/repo"
	stripeGateway "trading-card-app-backend/sellers/stripe"
)

func AddSellerAccount(sellerRepo *repo.SellerRepository, stripeKey string) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid := c.Param("uid")

		seller, err := sellerRepo.FindSeller(uid)

		if err != nil {
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
		return
	}
}
