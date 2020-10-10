package storage

import (
	"github.com/minio/minio-go"
	"github.com/pingcap/br/pkg/storage"
	"github.com/pingcap/kvproto/pkg/backup"
	"github.com/pingcap/log"
	"go.uber.org/zap"
)

type TempS3Storage struct {
	opts  *backup.S3
	Raw   string
	minio *minio.Client
}

func ConnectToS3(url string) (*TempS3Storage, error) {
	backend, err := storage.ParseBackend(url, &storage.BackendOptions{})
	if err != nil {
		return nil, err
	}
	s3Opts := backend.GetS3()
	log.Info("use temporary S3 storage", zap.Any("config", s3Opts), zap.String("url", url))
	minioClient, err := minio.New(s3Opts.Endpoint, s3Opts.AccessKey, s3Opts.SecretAccessKey, false)
	if err != nil {
		return nil, err
	}
	return &TempS3Storage{opts: s3Opts, Raw: url, minio: minioClient}, nil
}

func (s *TempS3Storage) Cleanup() error {
	for obj := range s.minio.ListObjects(s.opts.Bucket, s.opts.Prefix, true, nil) {
		if obj.Err != nil {
			return obj.Err
		}
		// TODO batching the request
		if err := s.minio.RemoveObject(s.opts.Bucket, obj.Key); err != nil {
			return err
		}
	}
	return nil
}
