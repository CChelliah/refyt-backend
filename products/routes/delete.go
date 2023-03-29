package routes

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"trading-card-app-backend/common/uow"
	"trading-card-app-backend/products/repo"
	stripeGateway "trading-card-app-backend/products/stripe"
)

func Delete(productRepo repo.ProductRepository, stripeKey string, uowManager uow.UnitOfWorkManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		productID := ctx.Param("productId")

		err := uowManager.Execute(ctx, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			err = stripeGateway.DeleteProduct(productID, stripeKey)

			if err != nil {
				return err
			}

			err = productRepo.DeleteProduct(ctx, uow, productID)

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

		ctx.JSON(http.StatusNoContent, "")
		return
	}
}
