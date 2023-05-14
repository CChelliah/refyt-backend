package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"mime/multipart"
	"os"
)

func UploadFile(productID string, productImages []*multipart.FileHeader) (imageUrls []string, err error) {

	secretKey, exists := os.LookupEnv("AWS_SECRET_ACCESS_KEY")

	if !exists {
		return imageUrls, fmt.Errorf("unable to find secret access key")
	}

	accessKey, exists := os.LookupEnv("AWS_ACCESS_KEY")

	if !exists {
		return imageUrls, fmt.Errorf("unable to find secret access key")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-2"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	},
	)

	if err != nil {
		return imageUrls, err
	}

	client := s3.New(sess)

	if err != nil {
		return imageUrls, err
	}

	imageUrls = []string{}

	for i, productImage := range productImages {

		file, err := productImage.Open()

		if err != nil {
			return imageUrls, err
		}

		bucket := "refyt"
		key := fmt.Sprintf("%s-%d.jpeg", productID, i)

		_, err = client.PutObject(&s3.PutObjectInput{
			Body:   file,
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

		imageUrls = append(imageUrls, fmt.Sprintf("https://d23djj3zawzk3r.cloudfront.net/%s", key))

		if err != nil {
			return imageUrls, err
		}

	}

	return imageUrls, nil
}
