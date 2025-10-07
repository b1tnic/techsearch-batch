package apiclient

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagent"
)

type MyBedrockClient struct {
	KnowledgeBaseID string
	DataSourceID    string
	client          *bedrockagent.Client
}

/** MyBedrockClientを作成、返却するメソッド
 *  @params AWSコンフィグ、ナレッジベースID、データソースID
 *  @return MyBedrockClient
 */
func NewMyBedrockClient(cfg aws.Config, knowledgeBaseID string, dataSourceID string) *MyBedrockClient {
	client := bedrockagent.NewFromConfig(cfg)
	return &MyBedrockClient{
		KnowledgeBaseID: knowledgeBaseID,
		DataSourceID:    dataSourceID,
		client:          client,
	}
}

/** Bedrockのナレッジベースのデータソースを同期するメソッド
 */
func (bc *MyBedrockClient) SyncKnowledgeBaseDataSource() {

	// 同期ジョブ開始
	result, err := bc.client.StartIngestionJob(context.TODO(), &bedrockagent.StartIngestionJobInput{
		KnowledgeBaseId: aws.String(bc.KnowledgeBaseID),
		DataSourceId:    aws.String(bc.DataSourceID),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("同期開始: Job ID = %s", *result.IngestionJob.IngestionJobId)
}
