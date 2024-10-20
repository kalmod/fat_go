package food

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

type FatSecretSearchOptions struct {
	SearchExpression string `url:"search_expression,omitempty"`
	PageNumber       int    `url:"page_number,omitempty"`
	MaxResults       int    `url:"max_results,omitempty"` // 20
	Format           string `url:"format,omitempty"`
}

// This function is to set the parameters of my search
func NewFatSecretSearchOptions(opts ...func(*FatSecretSearchOptions)) FatSecretSearchOptions {
	fs := FatSecretSearchOptions{Format: "json"}
	for _, o := range opts {
		o(&fs)
	}
	return fs
}

func FSS_SearchExpression(s string) func(*FatSecretSearchOptions) {
	return func(fss *FatSecretSearchOptions) {
		fss.SearchExpression = s
	}
}

func FSS_PageNumber(n int) func(*FatSecretSearchOptions) {
	return func(fss *FatSecretSearchOptions) {
		fss.PageNumber = n
	}
}
func FSS_MaxResults(n int) func(*FatSecretSearchOptions) {
	return func(fss *FatSecretSearchOptions) {
		fss.MaxResults = n
	}
}

func NewSearchRequest(searchOptions FatSecretSearchOptions, foodSearch_url string, accessTokenResponse AccessTokenJSON) (*http.Request, error) {

	query_values, _ := query.Values(searchOptions)   // this sets up the values so I can add to my requesturl
	parsed_url, err_url := url.Parse(foodSearch_url) // url we are going to target
	if err_url != nil {
		return nil, err_url
	}
	parsed_url.RawQuery = query_values.Encode() // adds our queries values to our parsed_url

	// create our new request using our complete parsed_url
	search_req, err := http.NewRequest(http.MethodGet, parsed_url.String(), nil)
	if err != nil {
		return nil, err
	}
	// add our headers
	search_req.Header.Add("Content-Type", "application/json")
	search_req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessTokenResponse.AccessToken))

	return search_req, nil
}
