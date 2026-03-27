package main

import (
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var urlPattern = regexp.MustCompile(`https?://[^\s"'<>` + "`" + `\(\)]+`)

func rewriteResponseHeaders(resp *http.Response, baseURL string) {
	if loc := resp.Header.Get("Location"); loc != "" {
		resp.Header.Set("Location", rewriteSingleURL(loc, baseURL))
	}
	if cl := resp.Header.Get("Content-Location"); cl != "" {
		resp.Header.Set("Content-Location", rewriteSingleURL(cl, baseURL))
	}
}

var rewritableTypes = []string{
	"application/json",
	"text/html",
	"text/xml",
	"text/plain",
	"application/xml",
	"application/xhtml",
	"text/javascript",
	"application/javascript",
}

func shouldRewriteBody(contentType string) bool {
	ct := strings.ToLower(contentType)
	for _, t := range rewritableTypes {
		if strings.Contains(ct, t) {
			return true
		}
	}
	return false
}

func rewriteBody(body []byte, baseURL string) []byte {
	return urlPattern.ReplaceAllFunc(body, func(match []byte) []byte {
		return []byte(rewriteSingleURL(string(match), baseURL))
	})
}

func rewriteSingleURL(rawURL, baseURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	scheme := parsed.Scheme
	if scheme != "http" && scheme != "https" {
		return rawURL
	}
	host := parsed.Hostname()
	portStr := parsed.Port()
	port := 0
	if portStr != "" {
		port, _ = strconv.Atoi(portStr)
	} else if scheme == "https" {
		port = 443
	} else {
		port = 80
	}
	path := parsed.RequestURI()
	return baseURL + "/" + scheme + "/" + host + "/" + strconv.Itoa(port) + path
}
