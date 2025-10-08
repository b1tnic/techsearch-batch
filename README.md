# Qiita記事に自然言語ができるWebアプリ[techsearch](https://techserch.net/)のバッチ処理のソースレポジトリ
(**現在、Qiitaサポートチームにサイト公開可否、データ利用可否を現在問い合わせしている最中なのでどのような文言を検索クエリにしてもダミーデータが返ってくるようになっています。**)

## 概要
日次で深夜二時に実行され、Qiitaの前日投稿分の記事をS3とDynamoDBに保存

## フローチャート
<img width="320" height="962" alt="techsearch-batch-flowchart drawio" src="https://github.com/user-attachments/assets/0386ee23-937b-42e0-8af4-c6ed217c7148" />


## ディレクトリ構造
```
.
├── apiclient apiパッケージ/
│   ├── BedrockClient.go Bedrock接続用クライアント
│   ├── DynamoDBClient.go DynamoDB接続用クライアント
│   ├── QiitaClient.go Qiita接続用クライアント
│   └── S3Client.go S3接続用クライアント
├── article 記事パッケージ/
│   └── Article.go 記事構造体
├── httpclient http接続用ヘルパーパッケージ/
│   └── http.go
├── urlbuilder URL作成用ヘルパーパッケージ/
│   └── url.go
├── go.mod
├── go.sum
└── main.go
```

## 実行方法
```
go run main.go
```

## 開発環境で必要な環境変数一覧
```
AWS_ACCESS_KEY_ID=AWSのアクセスキー
AWS_SECRET_ACCESS_KEY=AWSのシークレットアクセスキー
AWS_REGION=リージョン
PAGE_SIZE=一回当たりのQiita記事取得数
S3_BUCKETNAME=ナレッジベースの保管先S3バケット名
DYNAMODB_TABLENAME=記事保管先DynamoDBテーブル名
KNOWLEDGEBASE_ID=ナレッジベースID
DATASOURCE_ID=ナレッジベースで使用しているデータソースID
QIITA_ACCESS_TOKEN=Qiitaのアクセストークン
```
