package routes

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"trading-card-app-backend/common/uow"
	"trading-card-app-backend/products/domain"
	"trading-card-app-backend/products/repo"
	stripeGateway "trading-card-app-backend/products/stripe"
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

func Update(productRepo repo.ProductRepository, stripeKey string, uowManager uow.UnitOfWorkManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload updateProductPayload

		if err := ctx.Bind(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		productID := ctx.Param("productId")

		var product domain.Product

		err := uowManager.Execute(ctx, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			stripeProduct, err := stripeGateway.UpdateProduct(payload.ProductName, payload.Price, stripeKey, payload.Description, payload.RRP, payload.Designer, payload.FitNotes, productID)

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
			ctx.JSON(http.StatusNotFound, err.Error())
			return
		case err != nil:
			ctx.JSON(http.StatusInternalServerError, "internal server error")
			return
		}

		ctx.JSON(200, product)
		return
	}
}
