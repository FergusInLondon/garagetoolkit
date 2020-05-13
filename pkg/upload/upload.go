package upload

import (
	"io"
	"os"

	"github.com/minio/minio-go/v6"
)

var (
	client     *minio.Client
	bucketName string
)

type Configuration struct {
	MinIOEndpoint    string
	MinIOAccessKeyID string
	MinIOAccessKey   string
	UseSSL           bool
	BucketName       string
	BucketLocation   string
}

func Connect(conf *Configuration) error {
	c, err := minio.New(conf.MinIOEndpoint, conf.MinIOAccessKeyID, conf.MinIOAccessKey, conf.UseSSL)
	if err != nil {
		return err
	}

	client = c
	bucketName = conf.BucketName

	hasBucket, err := client.BucketExists(bucketName)
	if err != nil {
		return err
	}

	if !hasBucket {
		if err := client.MakeBucket(bucketName, conf.BucketLocation); err != nil {
			return err
		}
	}

	return nil
}

func Upload(file io.Reader, fileInfo os.FileInfo) (int, error) {
	return client.PutObject(bucketName, fileInfo.Name(), file, fileInfo.Size(), minio.PutObjectOptions{})
}
