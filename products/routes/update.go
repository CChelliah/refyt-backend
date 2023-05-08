package routes

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"refyt-backend/libs/events"
	"refyt-backend/libs/uow"
	"refyt-backend/products/domain"
	"refyt-backend/products/repo"
	"refyt-backend/products/s3"
	stripeGateway "refyt-backend/products/stripe"
)

type updateProductPayload struct {
	Name          string                  `form:"name" binding:"required"`
	Description   string                  `form:"description" binding:"required"`
	Designer      string                  `form:"designer" binding:"required"`
	Category      string                  `form:"category" binding:"required"`
	FitNotes      string                  `form:"fitNotes" binding:"required"`
	Size          int64                   `form:"size" binding:"required"`
	RRP           int64                   `form:"rrp" binding:"required"`
	Price         int64                   `form:"price" binding:"required"`
	ShippingPrice int64                   `form:"shippingPrice" binding:"required"`
	Images        []*multipart.FileHeader `form:"images"`
}

func Update(productRepo repo.ProductRepository, uowManager uow.UnitOfWorkManager, eventStreamer events.IEventStreamer) gin.HandlerFunc {
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

			product, err = productRepo.FindProductByID(productID)

			if err != nil {
				return err
			}

			stripeProduct, err := stripeGateway.UpdateProduct(payload.Name, payload.Price, payload.Description, payload.RRP, payload.Designer, payload.FitNotes, productID, payload.ShippingPrice)

			if err != nil {
				return err
			}

			imageUrls, err := s3.UploadFile(stripeProduct.ID, payload.Images)

			if err != nil {
				return err
			}

			err = stripeGateway.UpdateProductImages(stripeProduct.ID, imageUrls)

			if err != nil {
				return err
			}

			event := product.Update(&stripeProduct.ID, &payload.Name, &payload.Description, &payload.Price, &payload.RRP, &payload.Designer, &payload.FitNotes, &payload.Category, &payload.Size, &payload.ShippingPrice, &imageUrls)

			if err != nil {
				return err
			}

			product, err = productRepo.UpdateProduct(ctx, uow, stripeProduct.ID, product)

			if err != nil {
				return err
			}

			msg, err := events.ToEventPayload(event, string(events.ProductUpdatedEvent))

			if err != nil {
				return err
			}

			err = eventStreamer.Publish(string(events.ProductTopic), msg)

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
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(200, product)
	}
}
