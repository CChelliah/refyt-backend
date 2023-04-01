package routes

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"refyt-backend/libs/uow"
	"refyt-backend/products/repo"
	stripeGateway "refyt-backend/products/stripe"
)

func Delete(productRepo repo.ProductRepository, uowManager uow.UnitOfWorkManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		productID := ctx.Param("productId")

		err := uowManager.Execute(ctx, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			err = stripeGateway.DeleteProduct(productID)

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
	}
}
