package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (h *Client) Get(ctx context.Context, urlStr string, data url.Values) (*Response, error) {
	req, err := http.NewRequest("GET", urlStr+"?"+data.Encode(), nil)
	if err != nil {
		return nil, err
	}
	return h.Do(ctx, req)
}

func (h *Client) Post(ctx context.Context, url string, contentType string, body io.Reader) (*Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return h.Do(ctx, req)
}

func (h *Client) PostForm(ctx context.Context, url string, data url.Values) (*Response, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return h.Do(ctx, req)
}

func (h *Client) PostJson(ctx context.Context, url string, message json.RawMessage) (*Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return h.Do(ctx, req)
}

func (h *Client) PostFormDownload(ctx context.Context, url string, data url.Values) (*File, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return h.Download(ctx, req)
}
