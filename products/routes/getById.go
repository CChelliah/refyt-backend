package routes

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/products/repo"
)

func Get(productRepo repo.ProductRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		productID := c.Param("productId")

		product, err := productRepo.FindByID(productID)

		switch {
		case errors.Is(err, repo.ErrProductNotFound):
			c.JSON(http.StatusNotFound, err.Error())
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, "internal server error")
			return
		}

		fmt.Println(fmt.Sprintf("Product %s", product.ImageUrls[0]))

		c.JSON(200, product)
	}
}
