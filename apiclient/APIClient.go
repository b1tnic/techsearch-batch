package apiclient

import (
	"net/http"

	"github.com/b1tnic/techsearch-batch/article"
)

type APIClient interface {
	FetchArticles() (*http.Response, error)
	MapResponseToArticles(*http.Response) ([]article.Article, error)
}
