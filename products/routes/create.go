package routes

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"refyt-backend/common/uow"
	"refyt-backend/products/domain"
	"refyt-backend/products/repo"
	"refyt-backend/products/s3"
	stripeGateway "refyt-backend/products/stripe"
)

type createProductPayload struct {
	ProductName  string                `form:"productName" binding:"required"`
	Description  string                `form:"description" binding:"required"`
	Quantity     int64                 `form:"quantity" binding:"required"`
	Price        int64                 `form:"price" binding:"required"`
	Size         int64                 `form:"size" binding:"required"`
	RRP          int64                 `form:"rrp" binding:"required"`
	Designer     string                `form:"designer" binding:"required"`
	FitNotes     string                `form:"fitNotes" binding:"required"`
	ProductImage *multipart.FileHeader `form:"productImage"`
}

func Create(productRepo repo.ProductRepository, stripeKey string, uowManager uow.UnitOfWorkManager, accessKey string, secretKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload createProductPayload

		if err := ctx.ShouldBind(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		fmt.Println(payload.ProductImage)

		var product domain.Product

		err := uowManager.Execute(ctx, func(ctx context.Context, uow uow.UnitOfWork) (err error) {

			stripeProduct, err := stripeGateway.NewProduct(payload.ProductName, payload.Price, stripeKey, payload.Description, payload.RRP, payload.Designer, payload.FitNotes)

			if err != nil {
				return err
			}

			product, err = domain.CreateProduct(stripeProduct.ID, payload.ProductName, payload.Description, payload.Quantity, payload.Price, payload.RRP, payload.Designer, payload.FitNotes)

			if err != nil {
				return err
			}

			product, err = productRepo.InsertProduct(ctx, uow, product, "15xf5bidmhbPVSgMWHJSGMb32Vt1")

			if err != nil {
				return err
			}

			err = s3.UploadFile(accessKey, secretKey, payload.ProductImage)

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
