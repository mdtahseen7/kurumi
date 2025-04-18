package instagram

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"govd/util"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var captionPattern = regexp.MustCompile(
	`(?s)<meta property="og:title" content=".*?: &quot;(.*?)&quot;"`,
)

func BuildSignedPayload(contentURL string) (io.Reader, error) {
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	hash := sha256.New()
	_, err := io.WriteString(
		hash,
		contentURL+timestamp+apiKey,
	)
	if err != nil {
		return nil, fmt.Errorf("error writing to SHA256 hash: %w", err)
	}
	secretBytes := hash.Sum(nil)
	secretString := fmt.Sprintf("%x", secretBytes)
	secretString = strings.ToLower(secretString)
	payload := map[string]string{
		"url":  contentURL,
		"ts":   timestamp,
		"_ts":  apiTimestamp,
		"_tsc": "0", // ?
		"_s":   secretString,
	}
	parsedPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling payload: %w", err)
	}
	reader := strings.NewReader(string(parsedPayload))
	return reader, nil
}

func ParseIGramResponse(body []byte) (*IGramResponse, error) {
	var rawResponse interface{}
	if err := json.Unmarshal(body, &rawResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	switch rawResponse.(type) {
	case []interface{}:
		// array of IGramMedia
		var media []*IGramMedia
		if err := json.Unmarshal(body, &media); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return &IGramResponse{
			Items: media,
		}, nil
	case map[string]interface{}:
		// single IGramMedia
		var media IGramMedia
		if err := json.Unmarshal(body, &media); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return &IGramResponse{
			Items: []*IGramMedia{&media},
		}, nil
	default:
		return nil, fmt.Errorf("unexpected response type: %T", rawResponse)
	}
}

func GetCDNURL(contentURL string) (string, error) {
	parsedUrl, err := url.Parse(contentURL)
	if err != nil {
		return "", fmt.Errorf("can't parse igram URL: %v", err)
	}
	queryParams, err := url.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		return "", fmt.Errorf("can't unescape igram URL: %v", err)
	}
	cdnURL := queryParams.Get("uri")
	return cdnURL, nil
}

func GetPostCaption(
	postURL string,
) (string, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		postURL,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", util.ChromeUA)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "it-IT,it;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Referer", "https://www.instagram.com/accounts/onetap/?next=%2F")
	req.Header.Set("Alt-Used", "www.instagram.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", `csrftoken=Ib2Zuvf1y9HkDwXFxkdang; sessionid=8569455296%3AIFQiov2eYfTdSd%3A19%3AAYfVHnaxecWGWhyzxvz60vu5qLn05DyKgN_tTZUXTA; ds_user_id=8569455296; mid=Z_j1vQAEAAGVUE3KuxMR7vBonGBw; ig_did=BC48C8B7-D71B-49EF-8195-F9DE37A57B49; rur="CLN\0548569455296\0541775905137:01f7ebda5b896815e9279bb86a572db6bdc8ebccf3e1f8d5327e2bc5ca187fd5cd932b66"; wd=513x594; datr=x_X4Z_CHqpwtjaRKq7PtCNu3`)
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("TE", "trailers")

	resp, err := httpSession.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get response: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	matches := captionPattern.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		// post has no caption most likely
		return "", nil
	} else {
		return html.UnescapeString(matches[1]), nil
	}
}
