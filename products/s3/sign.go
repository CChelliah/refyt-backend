package s3

import (
	"refyt-backend/products/domain"
	"strings"
)

func Sign(product domain.Product) domain.Product {

	for j, imageUrl := range product.ImageUrls {

		signedUrl := strings.Replace(imageUrl, "https://refyt.s3.ap-southeast-2.amazonaws.com/", "https://d23djj3zawzk3r.cloudfront.net/", 1)

		product.ImageUrls[j] = signedUrl

	}

	return product

}
