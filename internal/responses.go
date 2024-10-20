package food

import (
	"encoding/json"
	"io"
	"net/http"
)

type FatResponseHandler struct {
	ResponseData []byte
	FatError     FatSecretError
	FatSearch    FatSecretSearchResult
}

func NewFatResponseHandler(opts ...func(*FatResponseHandler)) FatResponseHandler {
	frh := FatResponseHandler{}
	for _, o := range opts {
		o(&frh)
	}
	return frh
}

func FRH_Data(data []byte) func(*FatResponseHandler) {
	return func(frh *FatResponseHandler) {
		frh.ResponseData = data
	}
}

func FRH_DataRequest(client *http.Client, search_req *http.Request) func(*FatResponseHandler) {
	return func(frh *FatResponseHandler) {
		frh.NewSearch(client, search_req)
	}
}

func FRH_NewError() func(*FatResponseHandler) {
	return func(frh *FatResponseHandler) {
		frh.FatError = FatSecretError{}
	}
}

func FRH_NewSearch() func(*FatResponseHandler) {
	return func(frh *FatResponseHandler) {
		frh.FatSearch = FatSecretSearchResult{}
	}
}

func (frh *FatResponseHandler) ErrorUnmarshal() error {
	err := json.Unmarshal(frh.ResponseData, &frh.FatError)
	if err != nil {
		return err
	}
	return nil
}

func (frh *FatResponseHandler) SearchUnmarshal() error {
	err := json.Unmarshal(frh.ResponseData, &frh.FatSearch)
	if err != nil {
		return err
	}
	return nil
}

func (frh *FatResponseHandler) NewSearch(client *http.Client, search_req *http.Request) error {
	// Make our search request
	res, err := client.Do(search_req)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	frh.ResponseData = data
	return nil
}
