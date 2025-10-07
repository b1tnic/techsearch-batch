package urlbuilder

type Url struct {
	BaseUrl         string // 例: https://zenn.dev/api/articles
	queries         []Query
	queryParameters string
	CompleteUrl     string
}

type Query struct {
	key   string
	value string
}

/**
 * url全体を作成するメソッド
 */
func (u *Url) MakeCompleteUrl() {
	u.CompleteUrl = u.BaseUrl + u.queryParameters
}

/**
 * urlを初期化するメソッド
 */
func (u *Url) InitUrl() {
	u.CompleteUrl = ""
	u.queries = []Query{}
	u.queryParameters = ""
}

/**
 * urlのクエリ部分を作成するメソッド
 */
func (u *Url) MakeQueryParameters() {
	queryString := "?"
	for i, q := range u.queries {
		queryString += (q.key + "=" + q.value)

		// 最後のクエリ以外には&を付与
		if i != len(u.queries) {
			queryString += "&"
		}
	}

	u.queryParameters = queryString
}

/**
 * クエリを新たに追加するメソッド
 * @param クエリのキー、クエリの値
 */
func (u *Url) AddNewQuery(key string, value string) {
	query := Query{
		key:   key,
		value: value,
	}
	u.queries = append(u.queries, query)
}
