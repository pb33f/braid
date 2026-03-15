package model

import (
	"io"
	"net/http"
	"strings"
	"time"
)

// ConvertResponseIntoHttpResponse converts a HAR Response into a standard Go http.Response.
func ConvertResponseIntoHttpResponse(r Response) *http.Response {
	resp := &http.Response{
		StatusCode: r.StatusCode,
		Status:     r.StatusText,
		Proto:      r.HTTPVersion,
	}

	if r.Cookies != nil {
		h := http.Header{}
		for _, cookie := range r.Cookies {
			exp, _ := time.Parse(time.RFC3339Nano, cookie.Expires)
			c := &http.Cookie{
				Name:     cookie.Name,
				Path:     cookie.Path,
				Value:    cookie.Value,
				Domain:   cookie.Domain,
				Expires:  exp,
				HttpOnly: cookie.HTTPOnly,
				Secure:   cookie.Secure,
			}
			h.Add("Set-Cookie", c.String())
		}
		for _, header := range r.Headers {
			h.Add(header.Name, header.Value)
		}
		if r.RedirectURL != "" {
			h.Add("Location", r.RedirectURL)
		}
		if r.Body.MIMEType != "" {
			h.Add("Content-Type", r.Body.MIMEType)
		}
		if r.Body.Encoding != "" {
			h.Add("Content-Encoding", r.Body.Encoding)
		}
		if r.Body.Compression > 0 {
			h.Add("Content-Length", string(rune(r.Body.Compression)))
		}
		if r.Body.Size > 0 {
			h.Add("Content-Length", string(rune(r.Body.Size)))
		}
		if r.Body.Content != "" {
			resp.Body = io.NopCloser(strings.NewReader(r.Body.Content))
		}
		resp.Header = h
	}

	return resp
}
