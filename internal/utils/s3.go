package utils // Special thanks via https://docs.digitalocean.com/products/spaces/resources/s3-sdk-examples/

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func NewS3Client(key, secret, endpoint, region string) *s3.S3 {
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String(region),
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)
	return s3Client
}

func ListAllS3ObjectsInBucket(s3Client *s3.S3, bucketName string) *s3.ListObjectsOutput {
	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	}

	objects, err := s3Client.ListObjects(input)
	if err != nil {
		log.Println(err.Error())
	}

	return objects
}

func FindMatchingObjectKeyInS3Bucket(objects *s3.ListObjectsOutput, searchKey string) string {
	for _, obj := range objects.Contents {
		objKey := aws.StringValue(obj.Key)

		match := strings.Contains(objKey, searchKey)

		log.Println(objKey, searchKey, " | ", match)
		if match == true {
			return objKey
		}
	}
	return ""
}
