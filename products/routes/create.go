package routes

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"refyt-backend/libs/uow"
	"refyt-backend/products/domain"
	"refyt-backend/products/repo"
	"refyt-backend/products/s3"
	stripeGateway "refyt-backend/products/stripe"
)

type createProductPayload struct {
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

func Create(productRepo repo.ProductRepository, uowManager uow.UnitOfWorkManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload createProductPayload

		if err := ctx.ShouldBind(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		var product domain.Product

		err := uowManager.Execute(ctx, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			stripeProduct, err := stripeGateway.NewProduct(payload.Name, payload.Price, payload.Description, payload.RRP, payload.Designer, payload.FitNotes)

			if err != nil {
				return err
			}

			imageUrls, err := s3.UploadFile(stripeProduct.ID, payload.Images)

			if err != nil {
				return err
			}

			fmt.Println("%s", imageUrls)

			product, err = domain.CreateProduct(stripeProduct.ID, payload.Name, payload.Description, payload.Price, payload.RRP, payload.Designer, payload.FitNotes, payload.Category, payload.Size, payload.ShippingPrice)

			if err != nil {
				return err
			}

			product.AddImageUrls(imageUrls)

			if err != nil {
				return err
			}

			err = stripeGateway.UpdateProductImages(product.ProductID, imageUrls)

			if err != nil {
				return err
			}

			product, err = productRepo.InsertProduct(ctx, uow, product, "15xf5bidmhbPVSgMWHJSGMb32Vt1")

			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(200, product)
	}
}
