package s3

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"wa-service/app"
	client "wa-service/service/aws"

	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws"
)

func Upload(filename string, path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		_ = fmt.Errorf("error open file: %v\n", err)
	}
	defer file.Close()
	if err != nil {
		_ = fmt.Errorf("error Marshal file: %v\n", err)
	}
	fileInfo, _ := file.Stat()
	var size = fileInfo.Size()
	buffer := make([]byte, size)
	_, _ = file.Read(buffer)
	//upload to the s3 bucket
	_, err = s3.New(client.AwsClient).PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(os.Getenv("AWS_BUCKET")),
		ACL:           aws.String("public-read"),
		Key:           aws.String(filename),
		Body:          bytes.NewReader(buffer),
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(http.DetectContentType(buffer)),
	})
	if err != nil {
		log.Println("File Failed Upload !!", err)
		return "", err
	}
	log.Println("File Uploaded !!")
	return app.FilePath(filename), nil
}
