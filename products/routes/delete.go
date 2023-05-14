package routes

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"refyt-backend/libs/events"
	"refyt-backend/libs/uow"
	"refyt-backend/products/repo"
	stripeGateway "refyt-backend/products/stripe"
)

func Delete(productRepo repo.ProductRepository, uowManager uow.UnitOfWorkManager, eventStreamer events.IEventStreamer) gin.HandlerFunc {
	return func(c *gin.Context) {

		uid := c.GetString("uid")

		if uid == "" {
			err := fmt.Errorf("unauthorised user")

			zap.L().Error(err.Error())
			c.JSON(http.StatusUnauthorized, err.Error())
			return
		}

		productID := c.Param("productId")

		err := uowManager.Execute(c, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			err = stripeGateway.DeleteProduct(productID)

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			err = productRepo.DeleteProduct(ctx, uow, productID)

			if err != nil {
				zap.L().Error(err.Error())
				return err
			}

			return nil
		})

		switch {
		case errors.Is(err, repo.ErrProductNotFound):
			zap.L().Error(err.Error())
			c.JSON(http.StatusNotFound, err.Error())
			return
		case err != nil:
			zap.L().Error(err.Error())
			c.JSON(http.StatusInternalServerError, "internal server error")
			return
		}

		zap.L().Info(fmt.Sprintf("Successfully deleted product with id %s.", productID))

		c.JSON(http.StatusNoContent, "")
	}
}
