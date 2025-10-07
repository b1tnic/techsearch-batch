package apiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/b1tnic/techsearch-batch/article"
	"github.com/b1tnic/techsearch-batch/httpclient"
	"github.com/b1tnic/techsearch-batch/urlbuilder"
)

type QiitaClient struct {
	httpConnecter httpclient.HttpConnecter
	url           urlbuilder.Url
}

type QiitaArticle struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	LikesCount  int    `json:"likes_count"`
	StocksCount int    `json:"stocks_count"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
	Url         string `json:"url"`
}

/**
 * qiita記事を日付指定で取得するメソッド
 * @param 取得開始日、取得終了日、ページ数
 * @return httpレスポンス
 */
func (qc *QiitaClient) FetchArticles(dateFrom string, dateTo string, page string, perpage string) (*http.Response, error) {

	// 更新日時の指定
	period := fmt.Sprintf("updated:>=%s+updated:<=%s", dateFrom, dateTo)

	// タイムアウト10秒
	qc.httpConnecter.SetHttpClient(10)

	// Url初期化
	qc.url.InitUrl()

	qc.url.AddNewQuery("query", period)
	qc.url.AddNewQuery("page", page)
	qc.url.AddNewQuery("per_page", perpage)
	qc.url.MakeQueryParameters()
	qc.url.MakeCompleteUrl()
	qc.httpConnecter.Url = qc.url.CompleteUrl
	qc.httpConnecter.SetGetRequest()

	log.Printf("Qiita記事を叩きます。Url:%v", qc.httpConnecter.Url)

	// リクエストヘッダの初期化
	qc.httpConnecter.Request.Header = make(http.Header)

	accessToken := os.Getenv("QIITA_ACCESS_TOKEN")
	qc.httpConnecter.AddHeader("Authorization", "Bearer "+accessToken)

	return qc.httpConnecter.SendHttpRequest()
}

/** QiitaAPIのレスポンスヘッダーから、最大ページ数を返却するメソッド
 *  @params Httpレスポンス、一ページ当たりの記事数
 */
func (qc *QiitaClient) GetLastPageInDate(hr *http.Response, perpage string) int {
	log.Printf("%w", hr.Header)
	totalCount, err := strconv.Atoi(hr.Header.Get("Total-Count"))
	if err != nil {
		log.Fatalf("総記事数が変換できませんでした。エラー内容：%v", err)
	}

	perpageInt, _ := strconv.Atoi(perpage)

	if totalCount%perpageInt != 0 {
		return (totalCount/perpageInt + 1)
	}

	return totalCount / perpageInt
}

/**
 * qiitaAPIのレスポンスを記事の構造体にマッピングするメソッド
 * @param レスポンス
 * @return 記事、エラー
 */
func (qc *QiitaClient) MapResponseToArticles(hr *http.Response) ([]article.Article, error) {

	body, err := io.ReadAll(hr.Body)
	if err != nil {
		return nil, fmt.Errorf("Qiita記事のレスポンスボディが出力できませんでした。")
	}
	defer hr.Body.Close()

	var qiitaArticles []QiitaArticle
	json.Unmarshal(body, &qiitaArticles)

	articles := make([]article.Article, len(qiitaArticles))

	for i, qa := range qiitaArticles {
		articles[i] = mapQiitaArticleToArticle(&qa)
	}

	return articles, nil
}

/**
 * Qiita記事のjson構造体から、記事の構造体にマッピングを行う
 * @param qiita記事のjson構造体
 * @return 記事
 */
func mapQiitaArticleToArticle(qa *QiitaArticle) article.Article {
	article := article.Article{
		ID:          qa.ID,
		Title:       qa.Title,
		Body:        qa.Body,
		LikesCount:  qa.LikesCount,
		StocksCount: qa.StocksCount,
		UpdatedAt:   qa.UpdatedAt,
		CreatedAt:   qa.CreatedAt,
		Url:         qa.Url,
		Platform:    "Qiita",
	}

	return article
}

/**
 * qiitaAPIに接続するクライアントを作成
 * @param url
 */
func NewQiitaClient() *QiitaClient {

	// 更新日時、ページ数を指定
	qiitaUrl := urlbuilder.Url{}
	qiitaUrl.BaseUrl = "https://qiita.com/api/v2/items"

	hc := httpclient.NewHttpConnecter()

	qc := QiitaClient{
		httpConnecter: *hc,
		url:           qiitaUrl,
	}

	return &qc
}
