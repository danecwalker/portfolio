package notion

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BlockRequest struct {
	method  string
	blockId string
}

type BlockResponse struct {
	Object string `json:"object"`
	BlockGetResponse
	BlockChildrenResponse
}

type BlockGetResponse struct {
	Object           string     `json:"object"`
	ID               string     `json:"id"`
	Parent           PageParent `json:"parent"`
	Type             string     `json:"type"`
	CreatedTime      time.Time  `json:"created_time"`
	CreatedBy        User       `json:"created_by"`
	LastEditedTime   time.Time  `json:"last_edited_time"`
	LastEditedBy     User       `json:"last_edited_by"`
	Archived         bool       `json:"archived"`
	InTrash          bool       `json:"in_trash"`
	HasChildren      bool       `json:"has_children"`
	BulletedListItem *Paragraph `json:"bulleted_list_item,omitempty"`
	Paragraph        *Paragraph `json:"paragraph,omitempty"`
	Heading1         *Paragraph `json:"heading_1,omitempty"`
	Heading2         *Paragraph `json:"heading_2,omitempty"`
	Heading3         *Paragraph `json:"heading_3,omitempty"`
}

type Paragraph struct {
	RichText []RichText             `json:"rich_text"`
	Color    string                 `json:"color"`
	Children *BlockChildrenResponse `json:"children,omitempty"`
}

type BlockChildrenResponse struct {
	Results []BlockResponse `json:"results"`
}

func Blocks(blockId string) *BlockRequest {
	return &BlockRequest{
		method:  "GET",
		blockId: blockId,
	}
}

func (pr *BlockRequest) Query() *BlockRequest {
	pr.method = "CHILDREN"
	return pr
}

func (pr *BlockRequest) Fetch(c *Client) (*BlockResponse, error) {
	var method string
	var suffix string

	switch pr.method {
	case "GET":
		method = http.MethodGet
		suffix = ""
	case "CHILDREN":
		method = http.MethodGet
		suffix = "/children"
	default:
		method = http.MethodGet
		suffix = ""
	}

	url := c.buildURL("blocks/" + pr.blockId + suffix)
	fmt.Println(url)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var pageResp BlockResponse
	err = json.NewDecoder(resp.Body).Decode(&pageResp)
	if err != nil {
		return nil, err
	}
	return &pageResp, nil
}
