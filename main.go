package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/b1tnic/techsearch-batch/apiclient"
	"github.com/joho/godotenv"
)

// lambda起動
func main() {
	lambda.Start(handler)
}

func handler() {

	// .envファイルを読み込む
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 前日を指定
	FROM_DATE := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	// 当日を指定
	TO_DATE := time.Now().AddDate(0, 0, 0).Format("2006-01-02")
	PAGE_SIZE := os.Getenv("PAGE_SIZE")
	S3_BUCKETNAME := os.Getenv("S3_BUCKETNAME")
	DYNAMODB_TABLENAME := os.Getenv("DYNAMODB_TABLENAME")
	KNOWLEDGEBASE_ID := os.Getenv("KNOWLEDGEBASE_ID")
	DATASOURCE_ID := os.Getenv("DATASOURCE_ID")

	// AWSコンフィグの作成
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// qiita記事を取得
	qc := apiclient.NewQiitaClient()
	qiitaRespArticles, err := qc.FetchArticles(FROM_DATE, TO_DATE, "1", PAGE_SIZE)
	if err != nil {
		log.Fatalf("%v", err)
	}
	// 総ページ数を取得
	totalPageCount := qc.GetLastPageInDate(qiitaRespArticles, PAGE_SIZE)
	log.Printf("%w回、Qiita APIを叩きます。", totalPageCount)

	// S3、DynamoDBクライアントの作成
	s3Client := apiclient.NewMyS3Client(cfg, S3_BUCKETNAME)
	dynamoDBClient := apiclient.NewMyDynamoDBClient(cfg, DYNAMODB_TABLENAME)
	bedrockClient := apiclient.NewMyBedrockClient(cfg, KNOWLEDGEBASE_ID, DATASOURCE_ID)

	// 総ページ数だけ記事を取得し、S3とDynamoDBにアップロード
	for i := 1; i <= totalPageCount && i <= 100; i++ {
		qiitaRespArticles, err := qc.FetchArticles(FROM_DATE, TO_DATE, fmt.Sprint(i), PAGE_SIZE)
		if err != nil {
			log.Fatalf("%v", err)
		}
		log.Printf("%w", qiitaRespArticles)

		// 取得した記事をマッピング
		qiitaArticles, err := qc.MapResponseToArticles(qiitaRespArticles)
		if err != nil {
			log.Fatalf("%v", err)
		}

		// S3、DynamoDBにアップロード
		for _, obj := range qiitaArticles {
			s3Client.UploadSingleObject(obj.ID+".md", strings.NewReader(obj.Body))
			dynamoDBClient.Upload(obj)
		}
	}

	// Bedrockのデータソースを同期する
	bedrockClient.SyncKnowledgeBaseDataSource()
}
