package infra

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

type Minio struct {
	mc     *minio.Client
	bucket string
}

func NewClient(endpoint string, accessKeyID string, secretAccessKey string, bucket string) (*Minio, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})
	return &Minio{mc: minioClient, bucket: bucket}, err
}

func (m *Minio) CreateBucket(ctx context.Context) error {
	exists, err := m.mc.BucketExists(ctx, m.bucket)

	if err != nil {
		return err
	}

	if !exists {
		err = m.mc.MakeBucket(ctx, m.bucket, minio.MakeBucketOptions{Region: "us-east-1"})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Minio) PutObject(ctx context.Context, objectName string, data [][]byte) error {
	log.Printf("put object %s to bucket %s", objectName, m.bucket)
	exists, err := m.mc.BucketExists(ctx, m.bucket)
	if err != nil {
		return err
	}
	if !exists {
		err := m.mc.MakeBucket(ctx, m.bucket, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	var result []byte
	for _, part := range data {
		result = append(result, part...)
	}

	reader := bytes.NewReader(result)

	_, err = m.mc.PutObject(ctx, m.bucket, objectName, reader, int64(len(result)),
		minio.PutObjectOptions{ContentType: "text/plain"})

	if err != nil {
		return err
	}
	return nil
}
