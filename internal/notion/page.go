package notion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PageRequest struct {
	method string
	pageId string
}

type Cover struct {
	Type string `json:"type"`
}

type PageResponse struct {
	Object         string              `json:"object"`
	ID             string              `json:"id"`
	CreatedTime    time.Time           `json:"created_time"`
	CreatedBy      User                `json:"created_by"`
	LastEditedTime time.Time           `json:"last_edited_time"`
	LastEditedBy   User                `json:"last_edited_by"`
	Archived       bool                `json:"archived"`
	InTrash        bool                `json:"in_trash"`
	Icon           FileObject          `json:"icon"`
	Cover          FileObject          `json:"cover"`
	Properties     map[string]Property `json:"properties"`
	Parent         PageParent          `json:"parent"`
	URL            string              `json:"url"`
	PublicURL      string              `json:"public_url"`
}

func Page(pageId string) *PageRequest {
	return &PageRequest{
		method: "GET",
		pageId: pageId,
	}
}

func (pr *PageRequest) Fetch(c *Client) (*PageResponse, error) {
	url := c.buildURL("pages/" + pr.pageId)
	fmt.Println(url)
	req, err := http.NewRequest(pr.method, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var pageResp PageResponse
	err = json.NewDecoder(resp.Body).Decode(&pageResp)
	if err != nil {
		return nil, err
	}
	return &pageResp, nil
}
