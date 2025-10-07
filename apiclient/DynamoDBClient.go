package apiclient

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/b1tnic/techsearch-batch/article"
)

type MyDynamoDBClient struct {
	tableName string
	client    *dynamodb.Client
}

/** DynamoDBにアップロードするメソッド
 *  @params 記事
 */
func (myClient MyDynamoDBClient) Upload(article article.Article) {
	// 構造体をDynamoDBの形式に変換
	item, err := attributevalue.MarshalMap(article)
	if err != nil {
		log.Fatal(err)
	}

	// DynamoDBに保存
	_, err = myClient.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &myClient.tableName, // テーブル名
		Item:      item,
	})
	if err != nil {
		log.Printf("%w", err)
	}

	log.Println("記事がアップロードされました！")
}

/** MyDynamoDBClientを返却するメソッド
 *  @params テーブル名
 */
func NewMyDynamoDBClient(cfg aws.Config, tableName string) *MyDynamoDBClient {
	client := dynamodb.NewFromConfig(cfg)

	return &MyDynamoDBClient{
		tableName: tableName,
		client:    client,
	}
}
