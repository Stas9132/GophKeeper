package storage

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stas9132/GophKeeper/internal/config"
	"github.com/stas9132/GophKeeper/internal/logger"
)

type S3 struct {
	*minio.Client
	logger.Logger
}

func NewS3(l logger.Logger) (*S3, error) {
	minioClient, err := minio.New(config.S3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.S3AccessKeyID, config.S3SecretAccessKey, ""),
		Secure: config.S3UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &S3{minioClient, l}, nil
}
