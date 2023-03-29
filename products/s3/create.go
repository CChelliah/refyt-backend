package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"mime/multipart"
)

func UploadFile(accessKey string, secretKey string, productImage *multipart.FileHeader) (err error) {

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-2"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	},
	)

	if err != nil {
		return err
	}

	client := s3.New(sess)

	if err != nil {
		return err
	}

	file, err := productImage.Open()

	if err != nil {
		return err
	}

	bucket := "refyt"
	key := productImage.Filename

	_, err = client.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}
