package routes

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/libs/uow"
	"refyt-backend/products/domain"
	"refyt-backend/products/repo"
	stripeGateway "refyt-backend/products/stripe"
)

type updateProductPayload struct {
	ProductName *string `json:"productName"`
	Description *string `json:"description"`
	Quantity    *int64  `json:"quantity"`
	Price       *int64  `json:"price"`
	Size        *int64  `json:"size"`
	RRP         *int64  `json:"rrp"`
	Designer    *string `json:"designer"`
	FitNotes    *string `json:"fitNotes"`
}

func Update(productRepo repo.ProductRepository, uowManager uow.UnitOfWorkManager) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid := c.GetString("uid")

		if uid == "" {
			c.JSON(http.StatusUnauthorized, "unauthorized user")
			return
		}

		var payload updateProductPayload

		if err := c.Bind(&payload); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		productID := c.Param("productId")

		var product domain.Product

		err := uowManager.Execute(c, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			stripeProduct, err := stripeGateway.UpdateProduct(payload.ProductName, payload.Price, payload.Description, payload.RRP, payload.Designer, payload.FitNotes, productID)

			if err != nil {
				return err
			}

			product, err = productRepo.UpdateProduct(ctx, uow, stripeProduct.ID, stripeProduct.Name, stripeProduct.Description, 1, stripeProduct.DefaultPrice.UnitAmount)

			if err != nil {
				return err
			}

			return nil
		})

		switch {
		case errors.Is(err, repo.ErrProductNotFound):
			c.JSON(http.StatusNotFound, err.Error())
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, "internal server error")
			return
		}

		c.JSON(200, product)
	}
}
