package s3

import (
	"io/ioutil"
	"log"
	"os"
	client "wa-service/service/aws"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func Get(filename string) ([]byte, error) {
	b, err := s3.New(client.AwsClient).GetObject(&s3.GetObjectInput{
		Key:    aws.String(filename),
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	body, _ := ioutil.ReadAll(b.Body)
	return body, nil
}
