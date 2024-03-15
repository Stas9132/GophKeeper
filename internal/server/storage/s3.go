package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

const (
	endpoint        = "127.0.0.1:9000"
	accessKeyID     = "aHLytUVhTKOPMYD6nYA2"
	secretAccessKey = "F2Avh18pul7X8IsGhCTeWPnaQNhlOuda3iAYSO30"
	useSSL          = false
)

type S3 struct {
	minioClient *minio.Client
}

func NewS3() (*S3, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	err = minioClient.MakeBucket(context.TODO(), authBucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(context.TODO(), authBucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", authBucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", authBucketName)
	}

	return &S3{minioClient: minioClient}, nil
}

const (
	authBucketName = "auth-bucket"
	location       = "us-east-1"
)

func (s *S3) Register(ctx context.Context, user, password string) (bool, error) {
	// Upload the test file
	// Change the value of filePath if the file is in another location
	objectName := "testdata"
	filePath := "c:\\minio\\file.txt"
	contentType := "application/octet-stream"

	// Upload the test file with FPutObject
	info, err := s.minioClient.FPutObject(ctx, authBucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	return true, nil
}

func (s *S3) Login(ctx context.Context, user, password string) (bool, error) {
	return true, nil

}

func (s *S3) Logout(ctx context.Context) (bool, error) {
	return true, nil

}
