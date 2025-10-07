package httpclient

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type HttpConnecter struct {
	Url     string
	Request *http.Request
	Client  *http.Client
}

/**
 * HttpConnecterを作成するメソッド
 * @return HttpConnecterへの参照値
 */
func NewHttpConnecter() *HttpConnecter {
	hc := HttpConnecter{}
	return &hc
}

/** Getリクエストをセットするメソッド
 *  作成失敗時、エラーログを返却する。
 *
 *  @return エラー
 */
func (hc *HttpConnecter) SetGetRequest() error {
	req, err := http.NewRequest(http.MethodGet, hc.Url, nil)
	if err != nil {
		return fmt.Errorf("リクエストが作成できませんでした、エラー内容：%v", err)
	}

	hc.Request = req
	return nil
}

// Httpクライアントをセットするメソッド
func (hc *HttpConnecter) SetHttpClient(timeoutSec int) {
	client := &http.Client{
		Timeout: time.Duration(timeoutSec) * time.Second,
	}
	hc.Client = client
}

/** http.Headerにキーと値を追加するメソッド
 *  @param ヘッダーに追加するキー、バリュー
 */
func (hc *HttpConnecter) AddHeader(key string, value string) {
	hc.Request.Header.Add(key, value)
}

/** ClientとRequestを受け取り、通信を試みるメソッド
 *  通信エラー時、エラーログを返却する
 *
 *  @return レスポンス、エラー
 */
func (hc *HttpConnecter) SendHttpRequest() (*http.Response, error) {
	log.Printf("%w", hc.Request)
	res, err := hc.Client.Do(hc.Request)
	if err != nil {
		return res, fmt.Errorf("正常に通信できませんでした、エラー内容:%v", err)
	}

	return res, err
}
