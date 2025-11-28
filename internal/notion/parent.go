package notion

type PageParent struct {
	Type   string `json:"type"`
	PageID string `json:"page_id"`
}

type DatasourceParent struct {
	Type         string `json:"type"`
	DatasourceID string `json:"database_id"`
}
