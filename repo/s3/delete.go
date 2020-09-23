package s3

import (
	"os"
	client "wa-service/service/aws"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Delete(filename string) error {
	_, err := s3.New(client.AwsClient).DeleteObject(&s3.DeleteObjectInput{
		Key:    aws.String(filename),
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
	})
	return err
}
