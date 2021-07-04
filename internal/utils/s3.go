package utils // Special thanks via https://docs.digitalocean.com/products/spaces/resources/s3-sdk-examples/

import (
	// "log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Function connects to a specific S3 bucket instance and returns a connected
// instance structure. 
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

// Function returns a list of all the S3 stored objects sin a specific bucket.
func ListAllS3ObjectsInBucket(s3Client *s3.S3, bucketName string) *s3.ListObjectsOutput {
	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	}

	objects, err := s3Client.ListObjects(input)
	if err != nil {
		// log.Println(err.Error())
		panic(err.Error())
	}

	return objects
}

// Function will iterate over all the s3 objects to match the partial key with
// the actual key found in the S3 bucket.
func FindMatchingObjectKeyInS3Bucket(s3Objects *s3.ListObjectsOutput, partialKey string) string {
	for _, obj := range s3Objects.Contents {
		objKey := aws.StringValue(obj.Key)

		match := strings.Contains(objKey, partialKey)

        // If a match happens then it means we have found the ACTUAL KEY in the
		// s3 objects inside the bucket.
		if match == true {
			return objKey
		}
	}
	return ""
}
