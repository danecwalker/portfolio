package notion

import "net/http"

type Client struct {
	baseURL    string
	authToken  string
	httpClient http.Client
}

type HeaderRoundTripper struct {
	rt      http.RoundTripper
	headers map[string]string
}

func (h *HeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for k, v := range h.headers {
		req.Header.Set(k, v)
	}
	return h.rt.RoundTrip(req)
}

func NewClient(authToken string) *Client {
	return &Client{
		baseURL:   "https://api.notion.com/v1/",
		authToken: authToken,
		httpClient: http.Client{
			Timeout: http.DefaultClient.Timeout,
			Transport: &HeaderRoundTripper{
				rt: http.DefaultTransport,
				headers: map[string]string{
					"Authorization":  "Bearer " + authToken,
					"Notion-Version": "2025-09-03",
					"Content-Type":   "application/json",
				},
			},
		},
	}
}

func (c *Client) buildURL(endpoint string) string {
	return c.baseURL + endpoint
}
