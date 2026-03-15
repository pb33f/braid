package model

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ConvertRequestIntoHttpRequest converts a HAR Request into a standard Go http.Request.
func ConvertRequestIntoHttpRequest(req Request) (*http.Request, error) {
	r := bytes.NewReader([]byte(req.Body.Content))
	query := ""
	if len(req.QueryParams) > 0 {
		query = "?"
		for x, param := range req.QueryParams {
			amp := "&"
			if x < len(req.QueryParams)-1 {
				amp = ""
			}
			query += fmt.Sprintf("%s=%s%s", param.Name, param.Value, amp)
		}
	}

	httpRequest, err := http.NewRequest(req.Method, fmt.Sprintf("%s%s", req.URL, query), nil)

	values := make(url.Values)
	for _, param := range req.Body.Params {
		decoded, _ := url.QueryUnescape(param.Value)
		values.Add(param.Name, decoded)
	}
	if req.Method == http.MethodPost || req.Method == http.MethodPut || req.Method == http.MethodPatch {
		httpRequest.PostForm = values
		httpRequest.Form = values
	}

	if err != nil {
		return nil, err
	}
	for _, header := range req.Headers {
		if !strings.HasPrefix(header.Name, ":") {
			httpRequest.Header.Set(header.Name, header.Value)
		}
	}
	for _, cookie := range req.Cookies {
		httpRequest.AddCookie(&http.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Path:     cookie.Path,
			Domain:   cookie.Domain,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: cookie.HTTPOnly,
			Secure:   cookie.Secure,
		})
	}
	if req.Body.MIMEType != "" {
		httpRequest.Header.Set("Content-Type", req.Body.MIMEType)
	}
	body, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	httpRequest.Body = io.NopCloser(bytes.NewReader(body))
	httpRequest.ContentLength = int64(len(body))
	return httpRequest, nil
}
