package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pingcap/br/pkg/storage"
	"github.com/pingcap/errors"
	"github.com/pingcap/kvproto/pkg/backup"
	"github.com/pingcap/log"
	"go.uber.org/zap"
)

type TempS3Storage struct {
	opts    *backup.S3
	Raw     string
	session *session.Session
	svc     *s3.S3
}

func (s *TempS3Storage) Connect() error {
	qs := s.opts
	awsConfig := aws.NewConfig().
		WithMaxRetries(3).
		WithS3ForcePathStyle(qs.ForcePathStyle).
		WithRegion(qs.Region)
	if qs.Endpoint != "" {
		awsConfig.WithEndpoint(qs.Endpoint)
	}
	var cred *credentials.Credentials
	if qs.AccessKey != "" && qs.SecretAccessKey != "" {
		cred = credentials.NewStaticCredentials(qs.AccessKey, qs.SecretAccessKey, "")
	}
	if cred != nil {
		awsConfig.WithCredentials(cred)
	}
	awsSessionOpts := session.Options{
		Config: *awsConfig,
	}
	ses, err := session.NewSessionWithOptions(awsSessionOpts)
	if err != nil {
		return err
	}

	c := s3.New(ses)

	qs.Prefix += "/"
	s.svc = c
	s.session = ses
	return nil
}

func ConnectToS3(url string) (*TempS3Storage, error) {
	backend, err := storage.ParseBackend(url, &storage.BackendOptions{})
	if err != nil {
		return nil, err
	}
	s3Opts := backend.GetS3()
	s := &TempS3Storage{opts: s3Opts, Raw: url}
	if err := s.Connect(); err != nil {
		return nil, nil
	}
	return s, nil
}

func (s *TempS3Storage) Cleanup() error {
	var marker *string
	req := &s3.ListObjectsInput{
		Bucket:  aws.String(s.opts.Bucket),
		Prefix:  aws.String(s.opts.Prefix),
		MaxKeys: aws.Int64(1024),
	}
	for {
		req.Marker = marker
		res, err := s.svc.ListObjects(req)
		if err != nil {
			return err
		}
		objs := make([]*s3.ObjectIdentifier, 0, len(res.Contents))
		for _, f := range res.Contents {
			objs = append(objs, &s3.ObjectIdentifier{Key: f.Key})
		}
		delInput := &s3.DeleteObjectsInput{
			Bucket: aws.String(s.opts.Bucket),
			Delete: &s3.Delete{
				Objects: objs,
				Quiet:   aws.Bool(false),
			},
		}
		delRes, err := s.svc.DeleteObjects(delInput)
		if err != nil {
			return err
		}
		if len(delRes.Errors) != 0 {
			log.Error("cleanup meet errors", zap.Any("errors", delRes.Errors))
			return errors.New("failed to cleanup due to s3 error")
		}
		log.Info("objects deleted", zap.Int("size", len(delRes.Deleted)))
		if res.IsTruncated != nil && *res.IsTruncated {
			marker = res.NextMarker
		} else {
			return nil
		}
	}
}
