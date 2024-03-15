package storage

import (
	"bytes"
	"context"
	"crypto/md5"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stas9132/GophKeeper/internal/logger"
	"io"
	"log"
)

const (
	endpoint        = "127.0.0.1:9000"
	accessKeyID     = "aHLytUVhTKOPMYD6nYA2"
	secretAccessKey = "F2Avh18pul7X8IsGhCTeWPnaQNhlOuda3iAYSO30"
	useSSL          = false
)

type S3 struct {
	*minio.Client
	logger.Logger
}

func NewS3(l logger.Logger) (*S3, error) {
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
		if errBucketExists != nil || !exists {
			return nil, errBucketExists
		}
	} else {
		log.Printf("Successfully created %s\n", authBucketName)
	}

	return &S3{minioClient, l}, nil
}

const (
	authBucketName = "auth-bucket"
	location       = "us-east-1"
)

func (s *S3) Register(ctx context.Context, user, password string) (bool, error) {
	if err := s.MakeBucket(ctx, user, minio.MakeBucketOptions{Region: location, ObjectLocking: true}); err != nil {
		return false, err
	}
	hash := md5.Sum([]byte(password))
	info, err := s.PutObject(ctx, authBucketName, user, bytes.NewReader(hash[:]), int64(len(hash)), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		s.Error(err.Error())
		return false, err
	}
	s.Info(info.Key + " stored in auth-storage")
	return true, nil
}

func (s *S3) Login(ctx context.Context, user, password string) (bool, error) {
	data, err := s.GetObject(ctx, authBucketName, user, minio.GetObjectOptions{})
	if err != nil {
		s.Error(err.Error())
		return false, err
	}
	hash := md5.Sum([]byte(password))
	buf, err := io.ReadAll(data)
	if err != nil {
		s.Error(err.Error())
		return false, err
	}
	if c := bytes.Compare(hash[:], buf); c != 0 {
		s.Warn("Invalid credentials for user: " + user)
		return false, err
	}
	s.Info("Login complete")
	return true, nil
}

func (s *S3) Logout(ctx context.Context) (bool, error) {
	return true, nil
}
