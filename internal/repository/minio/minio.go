package minio

import (
	"github.com/Verce11o/yata-comments/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func NewMinio(cfg *config.Config) *minio.Client {
	client, err := minio.New(cfg.MinioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioConfig.AccessKey, cfg.MinioConfig.SecretKey, ""),
		Secure: cfg.MinioConfig.SSL,
	})

	if err != nil {
		log.Fatalf("")
	}
	return client
}
