package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DatasourceRequest struct {
	method       string
	datasourceId string
	body         []byte
}

type DatasourceResponse struct {
	Object string `json:"object"`
	DatasourceGetResponse
	DatasourceQueryResponse
}

type DatasourceGetResponse struct {
	ID               string                 `json:"id"`
	CreatedTime      time.Time              `json:"created_time"`
	LastEditedTime   time.Time              `json:"last_edited_time"`
	Properties       map[string]interface{} `json:"properties"`
	Parent           DatasourceParent       `json:"parent"`
	DatasourceParent PageParent             `json:"database_parent"`
	Archived         bool                   `json:"archived"`
	Inline           bool                   `json:"is_inline"`
	Icon             FileObject             `json:"icon"`
	Cover            FileObject             `json:"cover"`
	URL              string                 `json:"url"`
	InTrash          bool                   `json:"in_trash"`
}

type DatasourceQueryResponse struct {
	Results []PageResponse `json:"results"`
}

func Datasource(datasourceId string) *DatasourceRequest {
	return &DatasourceRequest{
		method:       "GET",
		datasourceId: datasourceId,
		body:         nil,
	}
}

func (pr *DatasourceRequest) Query(filter map[string]interface{}, sort []map[string]interface{}) *DatasourceRequest {
	pr.method = "QUERY"

	payload := map[string]interface{}{}
	if filter != nil {
		payload["filter"] = filter
	}
	if sort != nil {
		payload["sorts"] = sort
	}

	if filter == nil && sort == nil {
		return pr
	}

	jsonBytes, _ := json.Marshal(payload)

	pr.body = jsonBytes
	return pr
}

func (pr *DatasourceRequest) Fetch(c *Client) (*DatasourceResponse, error) {
	var reqBody io.Reader
	var method string
	var suffix string

	switch pr.method {
	case "GET":
		method = http.MethodGet
		reqBody = nil
		suffix = ""
	case "QUERY":
		method = http.MethodPost
		if pr.body == nil {
			reqBody = nil
		} else {
			reqBody = bytes.NewReader(pr.body)
		}
		suffix = "/query"
	default:
		method = http.MethodGet
		reqBody = nil
		suffix = ""
	}

	url := c.buildURL("data_sources/" + pr.datasourceId + suffix)
	fmt.Println(url)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var pageResp DatasourceResponse
	err = json.NewDecoder(resp.Body).Decode(&pageResp)
	if err != nil {
		return nil, err
	}
	return &pageResp, nil
}
