package apiclient

import (
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MyS3Client struct {
	downloader *manager.Downloader
	uploader   *manager.Uploader
	client     *s3.Client
	bucketName string
}

/** S3クライアントを返却するメソッド
 *  @params AWSコンフィグ、バケット名
 *  @return MyS3Client
 */
func NewMyS3Client(cfg aws.Config, bucketName string) *MyS3Client {
	client := s3.NewFromConfig(cfg)
	downloader := manager.NewDownloader(client)
	uploader := manager.NewUploader(client)

	return &MyS3Client{
		downloader: downloader,
		uploader:   uploader,
		client:     client,
		bucketName: bucketName,
	}
}

/** オブジェクトをアップロードするメソッド
 *  @params キー、アップロードするオブジェクト
 */
func (c *MyS3Client) UploadSingleObject(key string, reader io.Reader) {
	_, err := c.uploader.Upload(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(key),
		Body:   reader,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("S3へのアップロードが成功しました！")
}
